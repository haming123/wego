package gows

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type CloseCode int

const (
	//正常关闭; 无论为何目的而创建, 该链接都已成功完成任务.
	CloseNormalClosure CloseCode = 1000
	//终端离开：可能因为服务端错误, 也可能因为浏览器正从打开连接的页面跳转离开.
	CloseGoingAway CloseCode = 1001
	//协议错误：由于协议错误而中断连接.
	CloseProtocolError CloseCode = 1002
	//数据格式错误：由于接收到不允许的数据类型而断开连接
	CloseUnsupportedData CloseCode = 1003
	//保留
	CloseReserved CloseCode = 1004
	//没有收到预期的状态码.
	CloseNoCloseRcvd CloseCode = 1005
	//异常关闭：用于期望收到状态码时连接非正常关闭 (也就是说, 没有发送关闭帧).
	CloseAbnormalClosure CloseCode = 1006
	//由于收到了格式不符的数据而断开连接 (如文本消息中包含了非 UTF-8 数据).
	CloseInvalidPayload CloseCode = 1007
	//由于收到不符合约定的数据而断开连接. 这是一个通用状态码, 用于不适合使用 1003 和 1009 状态码的场景.
	ClosePolicyViolation CloseCode = 1008
	//由于收到过大的数据帧而断开连接.
	CloseMessageTooBig CloseCode = 1009
	//缺少扩展：客户端终止连接，因为期望一个或多个拓展, 但服务器没有.
	CloseMandatoryExtension CloseCode = 1010
	//内部错误：服务器终止连接，因为遇到异常
	CloseInternalError CloseCode = 1011
	//服务重启：服务器由于重启而断开连接.
	CloseServiceRestart CloseCode = 1012
	//稍后再试：服务器由于临时原因断开连接。
	CloseTryAgainLater CloseCode = 1013
	//错误的网关.
	CloseBadGateway CloseCode = 1014
	//握手错误：表示连接由于无法完成 TLS 握手而关闭 (例如无法验证服务器证书).
	CloseTLSHandshake CloseCode = 1015
)

func isValidCloseCode(code CloseCode) bool {
	switch code {
	case CloseReserved, CloseNoCloseRcvd, CloseAbnormalClosure, CloseTLSHandshake:
		return false
	}

	if code >= CloseNormalClosure && code <= CloseBadGateway {
		return true
	}
	if code >= 3000 && code <= 4999 {
		return true
	}

	return false
}

type CloseInfo struct {
	Code CloseCode
	Info string
}

func (cc *CloseInfo) Error() string {
	return cc.Info
}

func parseClosePayload(data []byte) (CloseInfo, error) {
	ce := CloseInfo{Code: CloseNoCloseRcvd, Info: ""}
	if len(data) >= 2 {
		ce.Code = CloseCode(binary.BigEndian.Uint16(data))
		if !isValidCloseCode(ce.Code) {
			return ce, fmt.Errorf("invalid close code %v", ce.Code)
		}
		ce.Info = string(data[2:])
	}
	return ce, nil
}

func MarshalCloseInfo(code CloseCode, text string) ([]byte, error) {
	if code == CloseNoCloseRcvd {
		return []byte{}, nil
	}

	if len(text) > maxControlFrameSize-2 {
		return []byte{}, errors.New("control frame length > 123")
	}

	var buff [maxControlFrameSize]byte
	data := buff[0 : 2+len(text)]
	binary.BigEndian.PutUint16(data, uint16(code))
	copy(data[2:], text)
	return data, nil
}
