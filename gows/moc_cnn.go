package gows

import (
	"bytes"
	"net"
	"time"
)

type mocConn struct {
	buff bytes.Buffer
}

func (cnn *mocConn) Read(b []byte) (n int, err error) {
	return cnn.buff.Read(b)
}

func (cnn *mocConn) Write(b []byte) (n int, err error) {
	return cnn.buff.Write(b)
}

func (cnn *mocConn) Close() error {
	cnn.buff.Reset()
	return nil
}

func (cnn *mocConn) LocalAddr() net.Addr {
	return nil
}

func (cnn *mocConn) RemoteAddr() net.Addr {
	return nil
}

func (cnn *mocConn) SetDeadline(t time.Time) error {
	return nil
}

func (cnn *mocConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (cnn *mocConn) SetWriteDeadline(t time.Time) error {
	return nil
}
