package wego

import (
	"errors"
	"github.com/haming123/wego/cache"
	log "github.com/haming123/wego/dlog"
	"github.com/haming123/wego/klog"
	"github.com/haming123/wego/wini"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ServerConfig struct {
	//是否启用 HTTPS，默认是false
	UseHttps bool `ini:"use_https"`
	//Http监听地址，默认为空
	HttpAddr string `ini:"http_addr"`
	//Http监听端口，默认为 8080
	HttpPort uint `ini:"http_port;default=8080"`
	//Https监听地址，默认为空
	HttpsAddr string `ini:"https_addr"`
	//Https监听端口，默认为 10443
	HttpsPort uint `ini:"https_port;default=10443"`
	//HTTPS证书路径
	HttpsCertFile string `ini:"cert_file"`
	//HTTPS证书 keyfile 的路径
	HttpsKeyFile string `ini:"key_file"`
	//设置 HTTP 的超时时间
	ReadTimeout time.Duration `ini:"read_timeout"`
	//设置 HTTP 的超时时间
	WriteTimeout time.Duration `ini:"write_timeout"`
	//POST请求的默认内存缓存大小，单位：M
	MaxBody int64 `ini:"max_body"`
	//是否开启 gzip，输出的内容会进行 gzip，根据Accept-Encoding来判断
	EnableGzip bool `ini:"gzip_on"`
	//压缩长度阈值，只有超过gzip_size的会被压缩返回
	GzipSize int64 `ini:"gzip_size"`
}

type SessionConfig struct {
	//缓是否开启session, 默认为 false
	SessionOn bool `ini:"session_on"`
	//session类型：cookie、cache
	SessionStore string `ini:"session_store;default=cookie"`
	//客户端的cookie的名称前缀
	CookieName string `ini:"cookie_name;default=wego"`
	//保存session数据的cookie域名, 默认空
	Domain string `ini:"domain"`
	//session 过期时间，单位：秒，默认值是3600
	LifeTime uint `ini:"life_time;default=3600"`
	//设置cookie的SameSite属性
	SameSite http.SameSite `ini:"samesite"`
	//session数据的hash字符串,若session的存储类型为cookie,则必须提供
	HashKey string `ini:"hash_key"`
}

type MemoryDbConfig struct {
	//用于缓存的内存大小，单位：M, 缺省：0（不限制）
	MaxSize uint64 `ini:"max_size"`
}

type RedisConfig struct {
	//Redis地址
	Address string `ini:"address"`
	//Redis登录密码
	DbPwd string `ini:"db_pwd"`
	//Redis数据库号码
	DbNum int `ini:"db_num"`
}

type MemcacheConfig struct {
	//逗号分隔的 memcached 主机列表
	Address string `ini:"address"`
}

type DlogConfig struct {
	//日志输出类型配置，0 终端 1 文件
	Output int `ini:"output"`
	//日志输出级别: 0 OFF 1 FATAL 2 ERROR 3 WARN 4 INFO 5 DEBUG
	Level int `ini:"level;default=5"`
	//日志文件存储路径，缺省：logs
	Path string `ini:"path;default=logs"`
	//是否在日志里面显示源码文件名和行号，默认 true
	ShowCaller bool `ini:"show_caller;default=true"`
	//是否将json、xml展开显示
	ShowIndent bool `ini:"show_indent;default=true"`
}

type KlogConfig struct {
	//是否开启Klog, 默认为 false
	KlognOn bool `ini:"klog_on"`
	//日志文件存储路径，缺省：logs
	Path string `ini:"path;default=logs"`
	//文件轮换类型：0 天 1 小时
	Rotate int `ini:"rotate"`
}

type WebConfig struct {
	wini.ConfigData
	//应用名称
	AppName string `ini:"app_name"`
	//服务器名称
	ServerName string `ini:"server_name"`
	//缓是否开启缓存
	CacheOn bool `ini:"cache_on"`
	//cache类型：redis、memcache、memory
	CacheStore string `ini:"cache_store;default=memory"`
	//是否显示请求日志，默认为 true
	ShowUrlLog bool `ini:"show_url_log;default=true"`
	//是否显示请求日志，默认为 true
	ShowSqlLog bool `ini:"show_sql_log;default=true"`
	//设置调试日志级别：0 OFF 1 FATAL 2 ERROR 3 WARN 4 INFO 5 DEBUG
	ShowDebugLog int `ini:"show_debug_log"`
	//防JSON劫持的前缀字符串
	JsonPrefix string `ini:"json_prefix"`
	//获取client ip的header
	IPHeader string `ini:"ip_header"`
	//Web服务配置
	ServerParam ServerConfig `ini:"server"`
	//Session配置
	SessionParam SessionConfig `ini:"session"`
	//内存缓存配置
	MemoryDbParam MemoryDbConfig `ini:"memory"`
	//Redis缓存配置
	RedisParam RedisConfig `ini:"redis"`
	//Memcache缓存配置
	MemcacheParam MemcacheConfig `ini:"memcache"`
	//日志模块配置
	DlogParam DlogConfig `ini:"dlog"`
	//统计日志配置
	KlogParam KlogConfig `ini:"klog"`
}

//加载配置文件
func getDefaultConfigFile() string {
	cur_path, _ := os.Getwd()
	fileName := "app.conf"
	file_path := filepath.Join(cur_path, fileName)
	if _, err := os.Stat(file_path); err == nil {
		return fileName
	} else {
		return ""
	}
}

//加载配置文件
func (this *WebConfig) LoadConfig(file_name ...string) error {
	fileName := "app.conf"
	if len(file_name) > 0 {
		fileName = file_name[0]
	}

	cur_path, err := os.Getwd()
	if err != nil {
		return err
	}

	file_path := filepath.Join(cur_path, fileName)
	err = wini.ParseFile(file_path, &this.ConfigData)
	if err != nil {
		return err
	}

	err = this.ConfigData.GetStruct(this)
	if err != nil {
		return err
	}

	debug_log.Infof("load config file = %s\n", file_path)
	return nil
}

//通过配置参数初始化Dlog
func (this *WebConfig) InitDlog() error {
	cfg := this.DlogParam
	if cfg.Output == 1 && len(cfg.Path) < 1 {
		return errors.New("log path is empty")
	}
	if cfg.Output == 1 {
		log.InitFileLogger(cfg.Path, log.Level(cfg.Level))
	} else {
		log.InitTermLogger(log.Level(cfg.Level))
	}
	log.ShowCaller(cfg.ShowCaller)
	log.ShowIndent(cfg.ShowIndent)

	str_log := "init dlog output=Term"
	if cfg.Output == 1 {
		str_log = "init dlog output=File"
	}
	str_log += " level=" + log.Level(cfg.Level).String()
	if len(cfg.Path) > 0 {
		str_log += " path=" + cfg.Path
	}
	debug_log.Info(str_log)
	return nil
}

//通过配置参数初始化Dlog
func (this *WebConfig) InitKlog() error {
	cfg := this.KlogParam
	if cfg.KlognOn == false {
		return nil
	}

	if len(cfg.Path) < 1 {
		return errors.New("log path is empty")
	}

	klog.InitEngine(cfg.Path, klog.RotateType(cfg.Rotate))
	klog.SetDebugLogLevel(klog.Level(this.ShowDebugLog))

	str_log := "init klog rotate=day"
	if cfg.Rotate == 1 {
		str_log = "init klog output=hour"
	}
	if len(cfg.Path) > 0 {
		str_log += " path=" + cfg.Path
	}
	debug_log.Info(str_log)
	return nil
}

//通过配置参数初创建CookieStore
func (this *WebConfig) NewCookieStore() (cache.CacheStore, error) {
	cfg_store := this.SessionParam
	if len(cfg_store.HashKey) < 1 {
		return nil, errors.New("hash_key is empty")
	}

	cookie_name := ""
	if cfg_store.CookieName != "" {
		cookie_name = cfg_store.CookieName + "_data"
	}

	store := cache.NewCookieStore(cookie_name, cfg_store.HashKey)
	return store, nil
}

//通过配置参数初创建RedisStore
func (this *WebConfig) NewRedisStore() (cache.CacheStore, error) {
	cfg_store := this.RedisParam
	if len(cfg_store.Address) < 1 {
		return nil, errors.New("redis address is empty")
	}
	store := cache.NewRedisStore(cfg_store.Address, cfg_store.DbPwd)
	return store, nil
}

//通过配置参数初创建MemcacheStore
func (this *WebConfig) NewMemcacheStore() (cache.CacheStore, error) {
	cfg_store := this.MemcacheParam
	if len(cfg_store.Address) < 1 {
		return nil, errors.New("memcached address is empty")
	}
	address := strings.Split(cfg_store.Address, ",")
	store := cache.NewMemcacheStore(address...)
	return store, nil
}

//通过配置参数初创建MemoryStore
func (this *WebConfig) NewMemoryStore() (cache.CacheStore, error) {
	cfg_store := this.MemoryDbParam
	max_size := cfg_store.MaxSize * 1024 * 1024
	store := cache.NewMemoryStore(max_size)
	return store, nil
}

//通过配置参数初始化Cache
func (this *WebConfig) InitCache() error {
	if this.CacheOn == false {
		return nil
	}
	if cache.GetCacheStore() != nil {
		return errors.New("cache store is create")
	}

	var err error
	var store cache.CacheStore
	if this.CacheStore == "redis" {
		store, err = this.NewRedisStore()
	} else if this.CacheStore == "memcache" {
		store, err = this.NewMemcacheStore()
	} else if this.CacheStore == "memory" {
		store, err = this.NewMemoryStore()
	} else {
		err = errors.New("invalid store type")
	}
	if err != nil {
		return err
	}

	cache.SetCacheStore(store)
	debug_log.Infof("init cache store=%s\n", this.CacheStore)
	return nil
}

//通过配置参数初始化Session
func (this *WebConfig) InitSession(sess *SessionEngine) error {
	cfg_sess := this.SessionParam
	if cfg_sess.SessionOn == false {
		return nil
	}

	if cfg_sess.SessionStore == "cookie" {
		store, err := this.NewCookieStore()
		if err != nil {
			return err
		}
		sess.Init(store)
	} else {
		store := cache.GetCacheStore()
		if store == nil {
			return errors.New("cache store is nil")
		}
		sess.Init(store)
	}

	max_age := cfg_sess.LifeTime
	if max_age > 0 {
		sess.SetMaxAge(max_age)
	}

	if cfg_sess.CookieName != "" {
		sess.SetCookieName(cfg_sess.CookieName)
	}

	debug_log.Infof("init session store=%s CookieName=%s max_age=%d\n",
		cfg_sess.SessionStore, cfg_sess.CookieName, cfg_sess.LifeTime)
	return nil
}
