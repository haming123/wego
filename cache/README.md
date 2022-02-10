# wego/cache

### 介绍
wego/cache是用于存取应用缓存数据一个模块。缓存数据通常存储在缓存存储引擎中，wego/cache通过存储引擎提供的API来访问缓存数据。
缓存存储引擎通常采用内存来存储缓存数据，并且采用hash表来组织和管理缓存数据，因此缓存存储引擎通常具有比数据库更高的访问性能。
使用缓存可以提升应用的访问性能，并能够大大降低应用对数据库的访问压力。另外Web服务的session数据也多存储于缓存引擎中。
wego/cache支持目前支持三种缓存引擎：
```
1）memcache
2）redis
3）memory
```
其中memory引擎使用本地内存来存储缓存数据，由于不需要通过网络来访问数据，memory引擎具有更高的访问性能。

### 安装
go get github.com/haming123/wego/cache

### 快速上手
实现看看如何使用memory引擎来存取缓存数据：
```go
package main
import (
	"fmt"
	"wego/cache"
)
func main()  {
	cache.InitMemoryStore()

	err := cache.Set("name", "hello")
	if err != nil {
		fmt.Println(err)
		return
	}

	val_str, err := cache.GetString("name")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(val_str)

	err = cache.Delete("name")
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

### 各种类型的数据的存取
为了方便各种类型的数据的存取, wego/cache提供了GetString、GetInt...等快捷函数：
```go
func TestMemoryStoreDataType(t *testing.T) {
	cache.InitMemoryStore()

	cache.Set("str_val", "hello")
	t.Log(cache.GetString("str_val"))

	cache.Set("bool_val", true)
	t.Log(cache.GetBool("bool_val"))

	cache.Set("int_val", 111)
	t.Log(cache.GetInt("int_val"))

	cache.Set("float_val", 123.45)
	t.Log(cache.GetFloat("float_val"))

	cache.Set("time_val", time.Now())
	t.Log(cache.GetTime("time_val"))
}
```

### 存取结构体数据
 wego/cache存储struct时，通过JSON序列化函数将结构体序列化为[]byte，然后将byte数组存储到存储引擎中。JSON序列化时只会序列化公开字段（大写开头的字段），
 因此在使用wego/cache时需要确保不要使用非公开字段。用户也可以通过struct的tag来设定是否需要序列化一个字段，或者是为字段指定一个序列化名称。
 ```go
 func TestMemoryStoreStruct(t *testing.T) {
 	type User struct {
 		Name 	string
 		Age 	int
 	}
 	user := User{Name:"lisi", Age:12}
 
 	cache.InitMemoryStore()
 
 	cache.Set("user", user)
 
 	user2 := User{}
 	err := cache.GetStuct("user", &user2)
 	if err != nil {
 		t.Error(err)
 	}
 	t.Log(user2)
 }
```

### 设定数据的有效期
在使用Set方法存储缓存数据时，您可以为数据指定一个有效期（单位：秒）
 ```go
func TestMemoryStoreMaxAge(t *testing.T) {
	cache.InitMemoryStore()

	err := cache.Set("name", "hello", 3)
	if err != nil {
		t.Error(err)
		return
	}

	val_str, err := cache.GetString("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val_str)

	time.Sleep(3*time.Second)
	val_str, err = cache.GetString("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val_str)
}
```

### 使用redis缓存数据
 ```go
func TestStoreCacheRedis(t *testing.T) {
	cache.InitRedisStore("127.0.0.1:6379", "redis_pwd")

	err := cache.Set("name", "hello")
	if err != nil {
		t.Error(err)
		return
	}

	val_str, err := cache.GetString("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val_str)

	err = cache.Delete("name")
	if err != nil {
		t.Error(err)
		return
	}
}
```

### 使用memcache缓存数据
 ```go
func TestStoreCacheMemcache(t *testing.T) {
	cache.InitMemcacheStore("127.0.0.1:11211")

	err := cache.Set("name", "hello")
	if err != nil {
		t.Error(err)
		return
	}

	val_str, err := cache.GetString("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val_str)

	err = cache.Delete("name")
	if err != nil {
		t.Error(err)
		return
	}
}
```

### 自定义存储引擎
wego/cache 模块采用了接口的方式来实现缓存功能，用户可以实现该缓存接口，从而使用自己的缓存引擎来存取缓存数据：
 ```go
type CacheStore interface {
	SaveData(ctx context.Context, key string, data []byte, max_age uint) error
	ReadData(ctx context.Context, key string) ([]byte, error)
	Exist(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
}
```
要使用自定义的引擎，需要首先创建引擎对象，例如
 ```go
store := &MyStroe{...}
```
然后通过cache.SetCacheStore设置存储引擎， 例如：
 ```go
cache.SetCacheStore(store)
```
然后结可以使用相应的存取函数了。