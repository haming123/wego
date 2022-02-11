package cache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"github.com/haming123/wego/cache/mod/memcache"
	"github.com/haming123/wego/cache/mod/redis"
)

type UserData struct {
	Ts int64
	Uid int64
	Qdid int64
	Code string
	Role int
	Err int
}

func TestMemcached(t *testing.T) {
	mcc := memcache.New(GetMemcachedAddr())
	mcc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

	it, err := mcc.Get("foo")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(it.Value))
}

var pool *redis.Pool
func init() {
	pool = &redis.Pool{
		MaxIdle:     8, // 最大空闲连接数
		MaxActive:   0, // 和数据库的最大连接数，0 表示没有限制
		IdleTimeout: 100, // 最大空闲时间
		Dial: func() (redis.Conn, error) { // 初始化连接的代码
			c, err := redis.Dial("tcp", GetRedisAddr())
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH",  GetRedisPwd()); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

func TestRedis(t *testing.T) {
	// 从 pool 中取出一个连接
	conn := pool.Get()
	defer conn.Close()

	// 向Redis写入一个数据
	_, err := conn.Do("Set", "name", "Rose")
	if err != nil {
		t.Error(err)
		return
	}
	// 取出
	r, err := redis.String(conn.Do("Get", "name"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("The name is", r)
}

func TestStoreCacheRedis(t *testing.T) {
	InitRedisStore(GetRedisAddr(), GetRedisPwd())

	err := Set("name", "hello")
	if err != nil {
		t.Error(err)
		return
	}

	val_str, err := GetString("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val_str)

	err = Delete("name")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestStoreCacheMemcache(t *testing.T) {
	InitMemcacheStore(GetMemcachedAddr())

	err := Set("name", "hello")
	if err != nil {
		t.Error(err)
		return
	}

	val_str, err := GetString("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val_str)

	err = Delete("name")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMemoryStore(t *testing.T) {
	InitMemoryStore()

	err := Set("name", "hello")
	if err != nil {
		t.Error(err)
		return
	}

	val_str, err := GetString("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val_str)

	err = Delete("name")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMemoryStoreMaxAge(t *testing.T) {
	InitMemoryStore()

	err := Set("name", "hello", 3)
	if err != nil {
		t.Error(err)
		return
	}

	val_str, err := GetString("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val_str)

	time.Sleep(3*time.Second)
	val_str, err = GetString("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val_str)
}

func TestMemoryStoreDataType(t *testing.T) {
	InitMemoryStore()

	Set("str_val", "hello")
	t.Log(GetString("str_val"))

	Set("bool_val", true)
	t.Log(GetBool("bool_val"))

	Set("int_val", 111)
	t.Log(GetInt("int_val"))

	Set("float_val", 123.45)
	t.Log(GetFloat("float_val"))

	Set("time_val", time.Now())
	t.Log(GetTime("time_val"))
}

func TestMemoryStoreStruct(t *testing.T) {
	type User struct {
		Name 	string
		Age 	int
	}
	user := User{Name:"lisi", Age:12}

	InitMemoryStore()

	Set("user", user)

	user2 := User{}
	err := GetStuct("user", &user2)
	if err != nil {
		t.Error(err)
	}
	t.Log(user2)
}

func mem_add()  {
	index := 0
	for {
		key := fmt.Sprintf("key_%d", index)
		val  := fmt.Sprintf("val_aaaaaaaaaaaaaaaaaaaaaaaaaaaaadfddddddddddddddddddddddddd_%d", index)
		Set( key, val, 30)
		index += 1
		if index > 10000 {
			index = 0
		}
		//fmt.Println("add :" + key)
		time.Sleep(time.Second)
	}
}

func mem_add2()  {
	for {
		index := rand.Int()%100
		key := fmt.Sprintf("key_%d", index)
		val  := fmt.Sprintf("val_%d", index)
		Set( key, val, 30)
		//fmt.Println("add :" + key)
		time.Sleep(time.Second)
	}
}

/*
func TestStoreCacheMemoryX(t *testing.T) {
	store := NewMemoryStore(1024)
	SetCacheStore(store)
	go mem_add()
	go mem_add2()
	for {
		index := rand.Int()%100
		key := fmt.Sprintf("key_%d", index)
		val  := fmt.Sprintf("val_%d", index)
		Set( key, val, 30)

		time.Sleep(time.Second*3)
		_, err := GetString( key)
		if err != nil {
			t.Error(err)
			return
		}
	}
}*/

//go test -v -run=none -bench="BenchmarkMemorySet" -benchmem
func BenchmarkMemorySet(b *testing.B) {
	store := NewMemoryStore(1024*1024*100)
	SetCacheStore(store)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%1000000)
		val  := fmt.Sprintf("val_aaaaaaaaaaaaaaaaaaaaaaaaaaaaadfddddddddddddddddddddddddd_%d", i%1000000)
		Set(key, []byte(val))
	}
	b.StopTimer()
}

//go test -v -run=none -bench="BenchmarkMemorySetGet" -benchmem
func BenchmarkMemorySetGet(b *testing.B) {
	store := NewMemoryStore(1024*1024*100)
	SetCacheStore(store)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%1000000)
		val  := fmt.Sprintf("val_aaaaaaaaaaaaaaaaaaaaaaaaaaaaadfddddddddddddddddddddddddd_%d", i%1000000)
		Set(key, []byte(val))

		key = fmt.Sprintf("key_%d", i%123456)
		GetString(key)
	}
	b.StopTimer()
}

