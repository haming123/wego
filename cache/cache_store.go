package cache

import "context"

type CacheStore interface {
	SaveData(ctx context.Context, key string, data []byte, max_age uint) error
	ReadData(ctx context.Context, key string) ([]byte, error)
	Exist(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
}
