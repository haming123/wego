package wego

import (
	"context"
	"github.com/haming123/wego/cache"
)

type SessionEngine struct {
	store      cache.CacheStore
	maxAge     uint
	cookieName string
}

func (this *SessionEngine) SetCookieName(cookieName string) {
	this.cookieName = cookieName
}

func (this *SessionEngine) SetMaxAge(max_age uint) {
	this.maxAge = max_age
}

func (this *SessionEngine) Init(store cache.CacheStore) {
	this.store = store
	this.maxAge = 24 * 3600
	this.cookieName = "sid"
}

func (this *SessionEngine) CreateSid() (string, error) {
	return CreateSid(), nil
}

func (this *SessionEngine) SaveData(ctx context.Context, sid string, data []byte) error {
	if this.store == nil {
		panic("store is nil")
	}
	sid = this.cookieName + "_" + sid
	return this.store.SaveData(ctx, sid, data, this.maxAge)
}

func (this *SessionEngine) ReadData(ctx context.Context, sid string) ([]byte, error) {
	if this.store == nil {
		panic("store is nil")
	}
	sid = this.cookieName + "_" + sid
	return this.store.ReadData(ctx, sid)
}
