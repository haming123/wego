package klog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func get_short_name(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' || file[i] == '\\' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return  file
}

func get_file_path(file_path string) string {
	ch := file_path[len(file_path) -1]
	if  ch != '/' && ch != '\\' {
		file_path += string(os.PathSeparator)
	}
	arr := make([]byte, len(file_path))
	for i:= 0; i < len(file_path) ; i++ {
		ch := file_path[i]
		if  ch == '/' || ch == '\\' {
			ch = os.PathSeparator
		}
		arr[i] = ch
	}
	return string(arr)
}

func GetTimeString(tm time.Time) string {
	const shortForm = "2006-01-02 15:04:05"
	return tm.Format(shortForm)
}

func GetRemoteIp(req *http.Request) string {
	remoteAddr := req.Header.Get("X-Forwarded-For")
	if len(remoteAddr) > 0{
		return 	remoteAddr
	}

	remoteAddr = req.RemoteAddr
	if ip := req.Header.Get("Remote_addr"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func GetHostIP() string {
	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		return ""
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

func GetLocalIP() string{
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	res := conn.LocalAddr().String()
	res = strings.Split(res, ":")[0]
	return res
}

func HttpGet(url string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return string(result), errors.New(fmt.Sprintf("StatusCode=%d :%s",  resp.StatusCode, url))
	}
	return string(result), nil
}

func HttpPost(url string, data interface{}, contentType string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return string(result), errors.New(fmt.Sprintf("StatusCode=%d : %s",  resp.StatusCode, url))
	}
	return string(result), nil
}