package gows

import (
	"io"
	"unsafe"
)

func trimChar(path string, cc byte) string {
	//trim left
	if len(path) > 0 {
		i := 0
		for ; i < len(path); i++ {
			if path[i] != cc {
				break
			}
		}
		if i > 0 {
			path = path[i:]
		}
	}

	//trim right
	if len(path) > 0 {
		i := len(path) - 1
		for ; i >= 0; i-- {
			if path[i] != cc {
				break
			}
		}
		if i < len(path)-1 {
			path = path[0 : i+1]
		}
	}

	return path
}

func trimSpace(path string) string {
	return trimChar(path, ' ')
}

func WriteAllTo(buff []byte, dst io.Writer) error {
	for len(buff) > 0 {
		nn, err := dst.Write(buff)
		if err != nil {
			return err
		}
		buff = buff[nn:]
	}
	return nil
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
