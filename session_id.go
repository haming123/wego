package wego

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const hexDigit = "0123456789abcdef"

var seq_no uint32 = 0
var mac_addr string = ""

func init() {
	rand.Seed(time.Now().UnixNano())
	seq_no = uint32(rand.Int31())
	mac_addr = getMacAddr()
}

func getMacAddrString(a net.HardwareAddr) string {
	if len(a) == 0 {
		return ""
	}
	buf := make([]byte, 0, len(a)*2)
	for _, b := range a {
		buf = append(buf, hexDigit[b>>4])
		buf = append(buf, hexDigit[b&0xF])
	}
	return string(buf)
}

func getMacAddr() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, netInterface := range netInterfaces {
		uuid := getMacAddrString(netInterface.HardwareAddr)
		if uuid != "" {
			return uuid
		}
	}
	return ""
}

func CreateSid() string {
	seq_no += 1
	return fmt.Sprintf("%s%x%08x", mac_addr, time.Now().Unix(), seq_no)
}

func CreateSidTmp() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	sid := hex.EncodeToString(b)
	return sid, nil
}
