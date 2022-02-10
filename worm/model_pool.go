package worm

import (
	"reflect"
	"sync"
)

type ModelPool struct {
	mutex 	sync.Mutex
	pool  	[]*DbModel
	size 	int
}

//md.md_pool != nil，说明是ModelPool分配的
//只有md.md_pool != nil才可以加入ModelPool
//加入ModelPool后令md.md_pool = nil，可以防止重复加入
func (p *ModelPool) Put(md *DbModel) {
	if md == nil {
		return
	}

	p.mutex.Lock()
	if md.md_pool == nil {
		p.mutex.Unlock()
		return
	}
	md.Reset()

	if debug_log.level >= LOG_DEBUG {
		t_ent := reflect.TypeOf(md.ent_ptr).Elem()
		debug_log.Debugf("put %s to pool", t_ent.Name())
	}

	num := len(p.pool)
	if p.size == 0 || num < p.size {
		p.pool = append(p.pool, md)
	}

	p.mutex.Unlock()
}

//分配一个model
//分配的model必须指向md_pool
func (p *ModelPool) Get() *DbModel {
	p.mutex.Lock()

	num := len(p.pool)
	if num < 1 {
		p.mutex.Unlock()
		return nil
	}

	last := num - 1
	md := p.pool[last]
	p.pool = p.pool[:last]
	md.md_pool = p

	if debug_log.level >= LOG_DEBUG {
		t_ent := reflect.TypeOf(md.ent_ptr).Elem()
		debug_log.Debugf("get %s from pool", t_ent.Name())
	}

	p.mutex.Unlock()
	return md
}