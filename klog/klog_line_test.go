package klog

import (
	"fmt"
	"strings"
	"testing"
)

func TestLineEncode (t *testing.T) {
	var buff LogBuffer
	str_val := "account`weixin"
	str_val_encode := "account\\`weixin"
	NewL("login").UserId("user_1").Add("login_type", str_val).Encode(&buff)

	str_data := string(buff.GetBytes())
	if strings.HasSuffix(str_data, " `T \n") == false {
		t.Error("encode error: incorrector line end")
	}
	if strings.Index(str_data, str_val_encode)< 0 {
		t.Error("account`weixin encode error")
	}
}

func TestLineDecode (t *testing.T) {
	var buff LogBuffer
	str_val := "account`weixin"
	NewL("login").UserId("user_1").Add("login_type", str_val).Encode(&buff)

	data := buff.GetBytes()
	nn := len(data)
	if data[nn-1] != '\n' || data[nn-2] != ' ' {
		t.Error("encode error: incorrector line end")
		return
	}
	data = data[0:nn-2]

	flag := GetLineEncodeFlag(data)
	if flag == false {
		t.Error("encode error: incorrector line end")
		t.Error(string(buff.GetBytes()))
		return
	}

	ret := LogLineDecode(data)
	t.Log(ret)
	if ret["login_type"] != "account`weixin" {
		t.Errorf("field decode error: %s", ret["login_type"])
	}
}

func TestKLogLinePool(t *testing.T) {
	ptr := getLineEnt()
	str1 := fmt.Sprintf("prt addr :%p", ptr)
	putLineEnt(ptr)

	ptr2 := getLineEnt()
	str2 := fmt.Sprintf("prt addr :%p", ptr2 )
	putLineEnt(ptr2)

	if str1 != str2 {
		t.Error("ptr1 != ptr2")
	}
}
