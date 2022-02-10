package klog

import (
	"testing"
)

func TestGetHostIP(t *testing.T) {
	ip := GetHostIP()
	t.Log("GetHostIP:" + ip)
}

func TestGetLocalIP(t *testing.T) {
	ip := GetLocalIP()
	t.Log("GetLocalIP:" + ip)
}
