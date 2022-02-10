package wego

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	log "github.com/haming123/wego/dlog"
	"github.com/haming123/wego/klog"
)

type ShutdownFunc func()
type WebEngine struct {
	RouteGroup
	//路由树对象
	Route 		RouteTree
	//html tepmplates缓存
	Templates 	WebTemplates
	//session引擎
	Session 	SessionEngine
	//配置参数
	Config 		WebConfig
	//查找路由前的hook
	beforeRoute HandlerFunc
	//WebContext对象池
	ctxPool    	sync.Pool
	//服务关闭
	onShutdown 	ShutdownFunc
	//获取Ip的header
	IPHeaders	[]string
	//hanlder for 401
	hanlder_401 HandlerFunc
	//hanlder for 404
	hanlder_404 HandlerFunc
	//hanlder for 500
	hanlder_500 HandlerFunc
	//route info for 404
	route_info_404	*RouteInfo
}

func (web *WebEngine) BeforRouter(handler HandlerFunc) {
	web.beforeRoute = handler
}

func (web *WebEngine) SetShutdown(hook ShutdownFunc) {
	web.onShutdown = hook
}

func (web *WebEngine)SetHandler401(handler HandlerFunc) {
	web.hanlder_401 = handler
}

func (web *WebEngine)SetHandler404(handler HandlerFunc) {
	web.hanlder_404 = handler
}

func (web *WebEngine)SetHandler500(handler HandlerFunc) {
	web.hanlder_500 = handler
}

//空handler
func HandlerNull(c *WebContext) {
}

//缺省404handler
func default_not_fund(c *WebContext) {
	c.WriteText(http.StatusNotFound, "404 NOT FOUND: " + c.Input.Request.URL.Path)
}

func (web *WebEngine)InitRoute404() {
	rinfo := &RouteInfo{}
	rinfo.handler_type = FT_CTX_HANDLER
	rinfo.handler_ctx = default_not_fund
	rinfo.handler_name = "NotFound"
	rinfo.func_name = "NotFound"
	rinfo.handler_ctx = web.hanlder_404
	if rinfo.handler_ctx == nil {
		rinfo.handler_ctx = default_not_fund
	}
	rinfo.filters = []FilterInfo{}
	web.route_info_404 = rinfo
}

func newEngine() *WebEngine {
	web := &WebEngine{}
	web.RouteGroup.InitRoot(web)
	web.Templates.Init()
	web.Config.JsonPrefix = "while(1);"
	web.IPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}
	web.ctxPool.New = func() interface{} {
		return newContext()
	}
	return web
}

func NewWeb() (*WebEngine, error) {
	web := newEngine()

	//设置Config的缺省值
	err := web.Config.GetStruct(&web.Config)
	if err != nil{
		return nil, err
	}

	//设置模块的调试日志的显示
	SetDebugLogLevel(Level(web.Config.ShowDebugLog))
	klog.SetDebugLogLevel(klog.Level(web.Config.ShowDebugLog))

	//初始化dlog
	err = web.Config.InitDlog()
	if err != nil{
		return nil, err
	}

	return web, nil
}

func InitWeb(file_name ...string) (*WebEngine, error) {
	web := newEngine()
	err := web.Config.LoadConfig(file_name...)
	if err != nil{
		return nil, err
	}

	//设置模块的调试日志的显示
	SetDebugLogLevel(Level(web.Config.ShowDebugLog))
	klog.SetDebugLogLevel(klog.Level(web.Config.ShowDebugLog))

	//参数处理
	if web.Config.JsonPrefix == "" {
		web.Config.JsonPrefix = "while(1);"
	}
	if web.Config.IPHeader != "" {
		web.IPHeaders = strings.Split(web.Config.IPHeader, ",")
	}

	//初始化dlog
	err = web.Config.InitDlog()
	if err != nil{
		return nil, err
	}

	return web, nil
}

func (web *WebEngine) initModule() error {
	//klog初始化
	err := web.Config.InitKlog()
	if err != nil{
		return err
	}

	//初始化缓存store
	err = web.Config.InitCache()
	if err != nil{
		return err
	}

	//session初始化
	err = web.Config.InitSession(&web.Session)
	if err != nil{
		return err
	}

	//将filter用于路由
	web.Route.initRouteFilter()
	//405 route init
	web.InitRoute404()
	return nil
}

func (web *WebEngine) cleanModule() {
	if web.onShutdown != nil {
		debug_log.Info("call closeHook")
		web.onShutdown()
	}

	//关闭dlog
	if web.Config.DlogParam.Output == 1 {
		debug_log.Info("close dlog")
		log.Close()
	}

	//关闭klog
	if web.Config.KlogParam.KlognOn {
		debug_log.Info("close klog")
		klog.Close()
	}
}

//首先关闭所有开启的监听器，然后关闭所有闲置连接，最后等待活跃的连接均闲置了才终止服务。
//若传入的context在服务完成终止前已超时，则Shutdown方法返回context的错误，否则返回任何由关闭服务监听器所引起的错误。
//当Shutdown方法被调用时，Serve、ListenAndServe及ListenAndServeTLS方法会立刻返回ErrServerClosed错误。
//用户的退出指令一般是SIGTERM或SIGINT（常常对应bash的Ctrl + C）
func gracefullShutdown(server *http.Server, quit chan<- bool) {
	waiter := make(chan os.Signal, 1)
	signal.Notify(waiter, syscall.SIGTERM, syscall.SIGINT)
	<-waiter

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		debug_log.Error(err)
	}
	close(quit)
}

func (web *WebEngine)nweHttpSever(addr string) *http.Server {
	server := &http.Server{ Addr: addr, Handler:web}
	if web.Config.ServerParam.ReadTimeout > 0 {
		server.ReadTimeout = time.Duration(web.Config.ServerParam.ReadTimeout) * time.Second
	}
	if web.Config.ServerParam.WriteTimeout > 0 {
		server.WriteTimeout = time.Duration(web.Config.ServerParam.WriteTimeout) * time.Second
	}
	return server
}

func (web *WebEngine) RunHTTP(addr string) error {
	err := web.initModule()
	if err != nil {
		return err
	}

	server := web.nweHttpSever(addr)
	quit := make(chan bool, 1)
	go gracefullShutdown(server, quit)

	debug_log.Infof("run http server host= %s\n", addr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-quit
	debug_log.Info("http server exit")
	web.cleanModule()
	return nil
}

func (web *WebEngine) RunTLS(addr string, certFile string, keyFile string) error {
	err := web.initModule()
	if err != nil {
		return err
	}

	server := web.nweHttpSever(addr)
	quit := make(chan bool, 1)
	go gracefullShutdown(server, quit)

	debug_log.Infof("run https server host= %s\n", addr)
	err = server.ListenAndServeTLS(certFile, keyFile)
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-quit
	debug_log.Info("http server exit")
	web.cleanModule()
	return nil
}

func (web *WebEngine) Run(addr ...string) (err error) {
	cfg := web.Config.ServerParam
	if cfg.UseHttps {
		address := fmt.Sprintf("%s:%d", cfg.HttpsAddr, cfg.HttpsPort)
		if len(addr) > 0 {
			address = addr[0]
		}
		certFile := cfg.HttpsCertFile
		keyFile := cfg.HttpsKeyFile
		return web.RunTLS(address, certFile, keyFile)
	} else {
		address := fmt.Sprintf("%s:%d", cfg.HttpAddr, cfg.HttpPort)
		if len(addr) > 0 {
			address = addr[0]
		}
		return web.RunHTTP(address)
	}
}

//修正path，并重新获取路由，若获取到则跳转的修正的路由
//返回：是否跳转
func (web *WebEngine) cleanAndRedirect(c *WebContext) bool {
	req := c.Input.Request
	method := c.Input.Request.Method
	path := cleanPath(c.Path)
	if c.Path != path {
		rinfo := web.Route.getRoute(method, path, &c.RouteParam)
		if rinfo != nil {
			code := http.StatusMovedPermanently
			if req.Method != http.MethodGet {
				code = http.StatusTemporaryRedirect
			}
			c.Redirect(code, path)
			debug_log.Debug("redirect to: " + path)
			return true
		}
	}
	return false
}

func (web *WebEngine)shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}
	if strings.Contains(req.Header.Get("Connection"), "Upgrade") {
		return false
	}
	if strings.Contains(req.Header.Get("Accept"), "text/event-stream") {
		return false
	}
	return true
}

func (web *WebEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := web.ctxPool.Get().(*WebContext)
	defer web.ctxPool.Put(c)

	c.reset()
	c.Config = &web.Config
	c.engine = web
	c.Input.Request = req
	c.Output.ResponseWriter = w
	c.Path = req.URL.Path
	if web.Config.ShowUrlLog {
		c.Start = time.Now()
	}

	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			debug_log.Errorf("%s\n", trace(message))
			c.AbortWithError(http.StatusInternalServerError, errors.New("Internal Server Error"))
		}
	}()

	if web.Config.ShowUrlLog == true {
		printReqInfo(c)
	}

	//执行BeforRouter钩子
	if web.beforeRoute != nil {
		web.beforeRoute(c)
		debug_log.Debug("call beforeRoute hook")
	}

	//获取路由
	method := c.Input.Request.Method
	rinfo := web.Route.getRoute(method, c.Path, &c.RouteParam)
	if rinfo == nil {
		c.Route = web.route_info_404
	} else {
		c.Route = rinfo
	}

	//是否gzip
	if web.Config.ServerParam.EnableGzip {
		c.Output.gzip_flag = web.shouldCompress(req)
		c.Output.gzip_size = web.Config.ServerParam.GzipSize
	}

	//开启session
	if web.Session.store != nil {
		err := c.Session.Read()
		if err != nil {
			debug_log.Error(err)
		}
	}

	//执行BeforeExec钩子
	//优先执行函数钩子，再者是struct钩子，再者是group钩子
	if c.Route.before_func != nil {
		c.Route.before_func(c)
		debug_log.Debug("call BeforeExec function")
	} else if c.Route.before_mthd != nil {
		c.Route.before_mthd.BeforeExec(c)
		debug_log.Debug("call BeforeExec method")
	}

	//执行过滤器
	c.filters = c.Route.filters
	if c.state.Status == 0 {
		c.Next()
	}

	//执行handler
	if c.state.Status == 0 {
		callHandler(c)
	}

	//执行AfterExe钩子
	//优先执行函数钩子，再者是struct钩子，再者是group钩子
	if c.Route.after_func != nil {
		c.Route.after_func(c)
		debug_log.Debug("call AfterExec function")
	} else if c.Route.after_mthd != nil {
		c.Route.after_mthd.AfterExec(c)
		debug_log.Debug("call AfterExec method")
	}

	err := c.Output.Flush()
	if err != nil {
		debug_log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, errors.New("Internal Server Error"))
	}

	if web.Config.ShowUrlLog == true {
		printExeInfo(c)
	}
}
