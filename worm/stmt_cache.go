package worm

import (
	"container/list"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type StmtItem struct {
	query 	string
	stmt   	*sql.Stmt
	tm_add 	int64
}

type StmtCache struct {
	cache       map[string]*list.Element
	cacheList   *list.List
	wasteList   *list.List
	mux_cahce   sync.RWMutex
	mux_clean   sync.Mutex
	max_size	int
	tm_clean 	int64
}

func NewStmtCache(max_num int) *StmtCache {
	ent := &StmtCache{}
	ent.cache = make(map[string]*list.Element)
	ent.cacheList = list.New()
	ent.wasteList = list.New()
	ent.max_size = max_num
	return ent
}

//从cacheList中获取一个stmt
//若存在，将stmt移动到列头
func (sc *StmtCache) GetStmt(query string) (*sql.Stmt, bool) {
	if ele, ok := sc.cache[query]; ok {
		sc.cacheList.MoveToFront(ele)
		return ele.Value.(*StmtItem).stmt, true
	}
	return nil, false
}

//添加一个stmt到从cacheList
//若存在，则直接移动到列头，否则添加一个item到cacheList中
//若stmt的数量超出max_size，则将最后的stmt移动到wasteList中
func (sc *StmtCache) SetStmt(query string, stmt *sql.Stmt) {
	if ele, ok := sc.cache[query]; ok {
		sc.cacheList.MoveToFront(ele)
		ele.Value.(*StmtItem).stmt = stmt
		return
	}

	ent := &StmtItem {}
	ent.stmt = stmt
	ent.query = query
	ele := sc.cacheList.PushFront(ent)
	sc.cache[query] = ele

	if sc.max_size != 0 && sc.cacheList.Len() > sc.max_size {
		sc.RemoveOldest()
	}
}

//cache_list的容量超过max_size时，将陈旧的stmt移到wasteList中
//若距离上次清理时间超过指定的值，则启动清理协程
func (sc *StmtCache) RemoveOldest() {
	ele := sc.cacheList.Back()
	if ele == nil {
		return
	}

	//从列表中移出
	sc.cacheList.Remove(ele)
	ent := ele.Value.(*StmtItem)
	//从map中移出
	delete(sc.cache, ent.query)

	//加入到废弃列表
	sc.mux_clean.Lock()
	ent.tm_add = time.Now().Unix()
	sc.wasteList.PushFront(ent)
	sc.mux_clean.Unlock()
}

//若距离上次清理时间>30秒，则进行实际清理
func (sc *StmtCache) TryCleanWasteList() {
	tm_now := time.Now().Unix()
	//fmt.Printf("time=%d on=%d off=%d\n", tm_now - sc.tm_clean, sc.cacheList.Len(), sc.wasteList.Len())
	if tm_now - sc.tm_clean < 30 {
		return
	}
	go sc.CleanWasteList(tm_now)
}

//清理wasteList中的stmt
//为了防止stmt清理时被其他协程使用，因此wasteList中的stmt一定要超过60秒才能close
//若不过过60秒，继续加入到wasteList中，等待下次清理
func (sc *StmtCache) CleanWasteList(tm_now int64) {
	sc.mux_clean.Lock()
	defer sc.mux_clean.Unlock()

	sc.tm_clean = tm_now
	if sc.wasteList.Len() < 1 {
		return
	}
	//fmt.Println("...CleanWasteList beg")

	tmp_list := sc.wasteList
	sc.wasteList = list.New()
	for e := tmp_list.Front(); e != nil; e = e.Next() {
		ent := e.Value.(*StmtItem)
		if tm_now - ent.tm_add > 60 {
			fmt.Println("close: " + ent.query + "!!!!!!")
			ent.stmt.Close()
		} else {
			//fmt.Printf("life=%d:%s......\n", tm_now - ent.tm_add, ent.query)
			sc.wasteList.PushFront(ent)
		}
	}
	//fmt.Println("...CleanWasteList end")
}

func (sc *StmtCache) CleanStmtList()  {
	sc.mux_cahce.Lock()
	defer sc.mux_cahce.Unlock()

	tmp_list := sc.cacheList
	sc.wasteList = list.New()
	for e := tmp_list.Front(); e != nil; e = e.Next() {
		ent := e.Value.(*StmtItem)
		ent.stmt.Close()
	}
}

func (sc *StmtCache) Close() {
	sc.CleanWasteList(0)
	sc.CleanStmtList()
}

//返回参数中的状态：1 没有命中cache， 2 命中cache
const (
	STMT_NOT_USE 		int = 0
	STMT_EXE_PREPARE 	int = 1
	STMT_USE_CACHE		int = 2
)
func (sc *StmtCache)Prepare(dbs *DbSession, db *sql.DB, query string) (*sql.Stmt, int, error) {
	sc.mux_cahce.RLock()
	stmt, ok := sc.GetStmt(query)
	if ok {
		sc.mux_cahce.RUnlock()
		sc.TryCleanWasteList()
		return stmt, STMT_USE_CACHE, nil
	}
	sc.mux_cahce.RUnlock()

	// double check
	sc.mux_cahce.Lock()
	stmt, ok = sc.GetStmt(query)
	if ok {
		sc.mux_cahce.Unlock()
		sc.TryCleanWasteList()
		return stmt, STMT_USE_CACHE, nil
	}

	var err error
	stmt, err = db.Prepare(query)
	if err == nil {
		sc.SetStmt(query, stmt)
	}
	sc.mux_cahce.Unlock()
	//dbs.logPrint(nil, "[SQL] Prepare")

	sc.TryCleanWasteList()
	return stmt, STMT_EXE_PREPARE, err
}
