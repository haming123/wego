package cache

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"
)

type DataItem struct {
	key 		string
	data 		[]byte
	tm_rw		time.Time
	tm_long		time.Duration
}

type DataCache struct {
	data_map    map[string]*list.Element
	data_list   *list.List
	temp_list   *list.List
	data_lock   sync.Mutex
	data_size	uint64
	max_size	uint64
	tm_clean 	time.Time
}

func NewMemoryStore(max_size ...uint64) *DataCache {
	mem := &DataCache{}
	mem.data_map = make(map[string]*list.Element)
	mem.data_list = list.New()
	mem.temp_list = list.New()
	mem.tm_clean = time.Now()
	if len(max_size) > 0 {
		mem.max_size = max_size[0]
	}
	if mem.max_size < 1 {
		mem.max_size = 900000000000	//90G
	}
	return mem
}

//若存在指定的item，修改item的属性，并将item前移到列头
//若存在指定的item，并且正在清理中，则从两个队列中移出，并添加到数据队列的列头
//若不存在，创建一个item，并添加到data_list列头
func (mem *DataCache) setData(key string, data []byte, max_age uint) {
	tm := time.Now()
	if ele, ok := mem.data_map[key]; ok {
		item := ele.Value.(*DataItem)
		item.tm_rw = tm
		mem.data_size -= uint64(len(item.data))
		item.data = data
		mem.data_size += uint64(len(item.data))
		item.tm_long = time.Second * time.Duration(max_age)
		if mem.temp_list.Len() > 0 {
			mem.temp_list.Remove(ele);mem.data_list.Remove(ele)
			mem.data_map[key] = mem.data_list.PushFront(item)
		} else {
			mem.data_list.MoveToFront(ele)
		}
		fmt.Printf("set: %s size=%d data=%d temp=%d clean=%d\n",
			key, mem.data_size, mem.data_list.Len(), mem.temp_list.Len(),
			tm.Unix()-mem.tm_clean.Unix())
	} else {
		item := &DataItem{key:key, data:data, tm_rw:tm, tm_long:time.Second * time.Duration(max_age)}
		mem.data_map[key] = mem.data_list.PushFront(item)
		mem.data_size += uint64(len(data))
		fmt.Printf("add: %s size=%d data=%d temp=%d clean=%d\n",
			key, mem.data_size, mem.data_list.Len(), mem.temp_list.Len(),
			tm.Unix()-mem.tm_clean.Unix())
	}

	//若到达了清理时间，首先将数据队列与临时队列交换，然后启动清理协程
	if mem.data_size > mem.max_size && tm.Sub(mem.tm_clean) > 60*time.Second {
		fmt.Printf("begin clean ...data=%d temp=%d\n", mem.data_list.Len(), mem.temp_list.Len())
		mem.tm_clean = tm.AddDate(1, 0, 0)
		data_list := mem.data_list
		mem.data_list = mem.temp_list
		mem.temp_list = data_list
		go mem.dataClean(tm)
	}
}

//从data_map中获取一个DataItem
//若存在，并且没有过期，将DataItem移动到data_list列头
//若存在，并且已经过期，删除该DataItem
func (mem *DataCache) getData(key string) ([]byte, bool) {
	tm := time.Now()
	if ele, ok := mem.data_map[key]; ok {
		item := ele.Value.(*DataItem)
		if item.tm_long > 0 && tm.Sub(item.tm_rw) > item.tm_long {
			mem.temp_list.Remove(ele);mem.data_list.Remove(ele)
			mem.data_size -= uint64(len(item.data))
			delete(mem.data_map, item.key)
			fmt.Println("get(delete): " + item.key)
			return nil, false
		} else if mem.temp_list.Len() > 0 {
			item.tm_rw = tm
			mem.temp_list.Remove(ele);mem.data_list.Remove(ele)
			mem.data_map[key] = mem.data_list.PushFront(item)
			fmt.Println("get(Remove): " + item.key)
			return item.data, true
		} else {
			item.tm_rw = tm
			mem.data_list.MoveToFront(ele)
			fmt.Println("get(MoveToFront): " + item.key)
			return item.data, true
		}
	}
	return nil, false
}

//删除一个DataItem
func (mem *DataCache) deleteData(key string) {
	if ele, ok := mem.data_map[key]; ok {
		item := ele.Value.(*DataItem)
		mem.temp_list.Remove(ele);mem.data_list.Remove(ele)
		mem.data_size -= uint64(len(item.data))
		delete(mem.data_map, item.key)
	}
}

//清理过期的数据
func (mem *DataCache) dataClean(tm time.Time) {
	//删除过期数据
	for {
		fmt.Printf("cleanning invalid...size=%d data=%d temp=%d\n", mem.data_size, mem.data_list.Len(), mem.temp_list.Len())
		//获取临时队列头列部item
		mem.data_lock.Lock()
		ele := mem.temp_list.Front()
		if ele == nil {
			mem.data_lock.Unlock()
			break
		}
		//若没有过期，添加到的data_list尾部
		//若过期，则删除
		item := mem.temp_list.Remove(ele).(*DataItem)
		if item.tm_long > 0 &&  tm.Sub(item.tm_rw) > item.tm_long {
			mem.data_size -= uint64(len(item.data))
			delete(mem.data_map, item.key)
			fmt.Println("delete invalid One: " + item.key)
		} else {
			mem.data_map[item.key] = mem.data_list.PushBack(item)
		}
		mem.data_lock.Unlock()
		//time.Sleep(time.Millisecond*10)
	}

	//强制删除最久没有访问的数据，直到小于指定的数据量
	for ; mem.data_size > mem.max_size ; {
		fmt.Printf("cleanning last...size=%d data=%d temp=%d\n", mem.data_size, mem.data_list.Len(), mem.temp_list.Len())
		mem.data_lock.Lock()
		ele := mem.data_list.Back()
		if ele == nil {
			mem.data_lock.Unlock()
			break
		}
		item := mem.data_list.Remove(ele).(*DataItem)
		mem.data_size -= uint64(len(item.data))
		delete(mem.data_map, item.key)
		mem.data_lock.Unlock()
		fmt.Println("delete last One: " + item.key)
		//time.Sleep(time.Millisecond*10)
	}

	//设置清理时间
	mem.tm_clean = time.Now()
}

func (mem *DataCache) SaveData(ctx context.Context, sid string, data []byte, max_age uint) error {
	mem.data_lock.Lock()
	mem.setData(sid, data, max_age)
	mem.data_lock.Unlock()
	return nil
}

func (mem *DataCache) ReadData(ctx context.Context, sid string) ([]byte, error) {
	mem.data_lock.Lock()
	data, found := mem.getData(sid)
	mem.data_lock.Unlock()
	if !found {
		return nil, nil
	}
	return data, nil
}

func (mem *DataCache) Exist(ctx context.Context, sid string) (bool, error) {
	mem.data_lock.Lock()
	_, found := mem.getData(sid)
	mem.data_lock.Unlock()
	if !found {
		return false, nil
	}
	return true, nil
}

func (mem *DataCache) Delete(ctx context.Context, key string) error {
	mem.data_lock.Lock()
	mem.deleteData(key)
	mem.data_lock.Unlock()
	return nil
}
