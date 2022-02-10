#### 向日志服务器发送日志
klog支持向wego日志服务器发送日志：
```go
package main
import (
	"klog"
	"time"
)
func main() {
	klog.InitKlog("./logs", klog.ROTATE_HOUR);
	klog.SetLogSender(klog4go.NewSender4Wego("http://39.108.252.54:10080", "acc1", "app1", klog.GetLocalIP(), 9090))
	defer klog.CloseDefaultKlog()

	klog.NewL("login").UserId("user_1").Add("login_type", "account`weixin").Output()
	klog.NewL("index").Client("127.0.0.1:666").Add("exe", 12).Output()
}
```
NewSender4Wego函数的参数如下：
NewSender4Wego(log_srv string, accid string, appid string, local_ip string, port int)
其中：
    log_srv：是日志服务器的地址
    accid：是在日志平台申请的账号ID
    appid：是在日志平台中创建的应用的ID
    local_ip：是本机的IP地址
    port:是写日志的应用对应的端口号，若没有则输入0
