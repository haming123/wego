package worm

import (
	"errors"
	"reflect"
	"sync"
)

type ModelPool struct {
	mutex    sync.Mutex
	ent_type reflect.Type
	pool     []*DbModel
	size     int
}

func NewModelPool(ent interface{}, size ...int) *ModelPool {
	t_ent := reflect.TypeOf(ent)
	t_ent = GetDirectType(t_ent)
	psize := 100
	if len(size) == 1 {
		psize = size[0]
	}

	p := &ModelPool{}
	p.size = psize
	p.pool = make([]*DbModel, psize)
	p.pool = p.pool[:0]
	p.ent_type = t_ent
	return p
}

func (p *ModelPool) createModel(dbs ...*DbSession) *DbModel {
	v_ent_ptr := reflect.New(p.ent_type)
	ent_ptr := v_ent_ptr.Interface()
	if len(dbs) != 1 {
		return Model(ent_ptr)
	} else {
		return dbs[0].Model(ent_ptr)
	}
}

var errPoolModel = errors.New("PoolModel")

//model入池
func (p *ModelPool) Put(md *DbModel) {
	if md == nil {
		return
	}
	if md.ent_type != p.ent_type {
		return
	}
	if md.Err == errPoolModel {
		return
	}

	p.mutex.Lock()
	num := len(p.pool)
	if num >= p.size {
		p.mutex.Unlock()
		return
	}

	md.Reset()
	md.Err = errPoolModel
	p.pool = append(p.pool, md)
	p.mutex.Unlock()
}

//分配一个model
func (p *ModelPool) Get(dbs ...*DbSession) *DbModel {
	p.mutex.Lock()
	num := len(p.pool)
	if num < 1 {
		p.mutex.Unlock()
		return p.createModel(dbs...)
	}

	last := num - 1
	md := p.pool[last]
	p.pool = p.pool[:last]
	p.mutex.Unlock()
	if len(dbs) == 1 {
		md.db_ptr = dbs[0]
	} else {
		md.db_ptr = db_default.def_session
	}
	md.Err = nil
	return md
}
