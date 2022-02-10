package cache

import (
	"context"
	"errors"
	"github.com/haming123/wego/cache/mod/memcache"
)

//"github.com/bradfitz/gomemcache/memcache"

type StoreMemcached struct {
	mccon	*memcache.Client
	addrs 	[]string
}

func NewMemcacheStore(addr ...string) *StoreMemcached {
	mc := memcache.New(addr...)
	rs := &StoreMemcached{mccon:mc, addrs :addr}
	return rs
}

func (s *StoreMemcached) SaveData(ctx context.Context, key string, data []byte, max_age uint) error {
	mc := s.mccon
	if mc == nil {
		return errors.New("memcached not init")
	}

	var err error
	if max_age > 0 {
		err = mc.Set(&memcache.Item{Key: key, Value: data, Expiration:int32(max_age)})
	} else {
		err = mc.Set(&memcache.Item{Key: key, Value: data})
	}

	return err
}

func (s *StoreMemcached) ReadData(ctx context.Context, key string) ([]byte, error) {
	mc := s.mccon
	if mc == nil {
		return nil, errors.New("memcached not init")
	}

	item, err := mc.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}
		return nil, err
	}

	return item.Value, nil
}

func (s StoreMemcached) Exist(ctx context.Context, key string) (bool, error) {
	mc := s.mccon
	if mc == nil {
		return false, errors.New("memcached not init")
	}
	_, err := mc.Get(key)
	return err == nil, err
}

func (s StoreMemcached) Delete(ctx context.Context, key string) error {
	mc := s.mccon
	if mc == nil {
		return errors.New("memcached not init")
	}
	return mc.Delete(key)
}

