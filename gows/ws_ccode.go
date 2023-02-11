package gows

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type CloseCode int

const (
	CloseNormalClosure      CloseCode = 1000
	CloseGoingAway          CloseCode = 1001
	CloseProtocolError      CloseCode = 1002
	CloseUnsupportedData    CloseCode = 1003
	CloseReserved           CloseCode = 1004
	CloseNoCloseRcvd        CloseCode = 1005
	CloseAbnormalClosure    CloseCode = 1006
	CloseInvalidPayload     CloseCode = 1007
	ClosePolicyViolation    CloseCode = 1008
	CloseMessageTooBig      CloseCode = 1009
	CloseMandatoryExtension CloseCode = 1010
	CloseInternalError      CloseCode = 1011
	CloseServiceRestart     CloseCode = 1012
	CloseTryAgainLater      CloseCode = 1013
	CloseBadGateway         CloseCode = 1014
	CloseTLSHandshake       CloseCode = 1015
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
	Text string
}

func parseClosePayload(data []byte) (CloseInfo, error) {
	ce := CloseInfo{Code: CloseNoCloseRcvd, Text: ""}
	if len(data) >= 2 {
		ce.Code = CloseCode(binary.BigEndian.Uint16(data))
		if !isValidCloseCode(ce.Code) {
			return ce, fmt.Errorf("invalid close code %v", ce.Code)
		}
		ce.Text = string(data[2:])
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
