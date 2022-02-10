package cache

import (
	"context"
	"errors"
	"github.com/haming123/wego/cache/mod/redis"
)

//"github.com/gomodule/redigo/redis"

type StoreRedis struct {
	pool	*redis.Pool
}

func (s *StoreRedis) ping() (bool, error) {
	conn := s.pool.Get()
	defer conn.Close()
	data, err := conn.Do("PING")
	if err != nil || data == nil {
		return false, err
	}
	return (data == "PONG"), nil
}

func NewRedisStoreWithPool(pool *redis.Pool) *StoreRedis {
	rs := &StoreRedis{pool}
	_, err := rs.ping()
	if err != nil {
		panic(err)
	}
	return rs
}

/*
type Pool struct {
    // Dial()方法返回一个连接，从在需要创建连接到的时候调用
    Dial func() (Conn, error)
    // TestOnBorrow()方法是一个可选项，该方法用来诊断一个连接的健康状态
    TestOnBorrow func(c Conn, t time.Time) error
   // 最大空闲连接数
    MaxIdle int
    // 一个pool所能分配的最大的连接数目
    // 当设置成0的时候，该pool连接数没有限制
    MaxActive int
    // 空闲连接超时时间，超过超时时间的空闲连接会被关闭。
    // 如果设置成0，空闲连接将不会被关闭
    // 应该设置一个比redis服务端超时时间更短的时间
    IdleTimeout time.Duration
} */
func NewRedisStoreWithDB(address string, password string, db string) *StoreRedis {
	pool := &redis.Pool{
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	return NewRedisStoreWithPool(pool)
}

func NewRedisStore(address string, password string) *StoreRedis {
	pool := &redis.Pool{
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	return NewRedisStoreWithPool(pool)
}

func (s *StoreRedis) Close() error {
	return s.pool.Close()
}

//SETEX 命令：设置指定 key 的值为 value，并将 key 的过期时间设为 seconds (以秒为单位)
//如果 key 已经存在， SETEX 命令将会替换旧的值。
func (s *StoreRedis) SaveData(ctx context.Context, key string, data []byte, max_age uint) error {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}
	defer conn.Close()

	var err error
	if max_age > 0 {
		_, err = conn.Do("SETEX", key, max_age, data)
	} else {
		_, err = conn.Do("SET", key, data)
	}
	return err
}

func (s *StoreRedis) ReadData(ctx context.Context, key string) ([]byte, error) {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	defer conn.Close()

	res, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	data, err := redis.Bytes(res, err)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *StoreRedis) Exist(ctx context.Context, key string) (bool, error) {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return false, err
	}
	defer conn.Close()

	is_exit, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}
	return is_exit, nil
}

func (s *StoreRedis) Delete(ctx context.Context, key string) error {
	conn := s.pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}
	defer conn.Close()

	existed, err := redis.Bool(conn.Do("DEL", key))
	if err == nil && !existed {
		return errors.New("key not found")
	}
	return err
}
