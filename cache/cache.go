package cache

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"
	"github.com/haming123/wego/cache/mod/redis"
)

var cache_store CacheStore

func InitMemoryStore(max_size ...uint64) {
	cache_store = NewMemoryStore(max_size...)
}

func InitMemcacheStore(addr ...string) {
	cache_store = NewMemcacheStore(addr...)
}

func InitRedisStore(address string, password string) {
	cache_store = NewRedisStore(address, password)
}

func InitRedisStoreWithDB(address string, password string, db string) {
	cache_store = NewRedisStoreWithDB(address, password, db)
}

func InitRedisStoreWithPool(pool *redis.Pool) {
	cache_store = NewRedisStoreWithPool(pool)
}

func SetCacheStore(store CacheStore) {
	cache_store = store
}

func GetCacheStore() CacheStore {
	return cache_store
}

func Delete(key string) error {
	if cache_store == nil {
		panic("cache store is nil")
	}
	return cache_store.Delete(nil, key)
}

func Set(key string, value interface{}, max_age ...uint) error {
	if cache_store == nil {
		panic("cache store is nil")
	}

	var max_age_val uint = 0
	if len(max_age) > 0 {
		max_age_val = max_age[0]
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cache_store.SaveData(nil, key, data, max_age_val)
}

func GetBytes(key string) ([]byte, error) {
	if cache_store == nil {
		panic("cache store is nil")
	}
	data, err := cache_store.ReadData(nil, key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetString(key string) (string, error) {
	data, err := GetBytes(key)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetBool(key string) (bool, error) {
	data, err := GetBytes(key)
	if err != nil {
		return false, err
	}

	var val bool
	err = json.Unmarshal(data, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func GetInt(key string) (int64, error) {
	data, err := GetBytes(key)
	if err != nil {
		return 0, err
	}

	var val int64
	err = json.Unmarshal(data, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func GetFloat(key string) (float64, error) {
	data, err := GetBytes(key)
	if err != nil {
		return 0, err
	}

	var val float64
	err = json.Unmarshal(data, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func GetTime(key string) (time.Time, error) {
	data, err := GetBytes(key)
	if err != nil {
		return time.Time{}, err
	}

	var val time.Time
	err = json.Unmarshal(data, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func GetStuct(key string, ptr interface{}) error {
	if ptr == nil {
		return errors.New("ptr must be *Struct")
	}

	v_ent := reflect.ValueOf(ptr)
	if v_ent.Kind() != reflect.Ptr {
		return errors.New("ptr must be *Struct")
	}

	data, err := GetBytes(key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, ptr)
	if err != nil {
		return err
	}
	return nil
}
