package gows

import (
	"bufio"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

func genAcceptKey(challengeKey string) string {
	h := sha1.New()
	h.Write(StringToBytes(challengeKey))
	h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func generateChallengeKey() (string, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p), nil
}

func OriginHostCheck(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}
	return strings.EqualFold(u.Host, r.Host)
}

func getSubProtocols(r *http.Request) []string {
	h := strings.TrimSpace(r.Header.Get("Sec-Websocket-Protocol"))
	if h == "" {
		return nil
	}
	protocols := strings.Split(h, ",")
	for i := range protocols {
		protocols[i] = strings.TrimSpace(protocols[i])
	}
	return protocols
}

func Accept(w http.ResponseWriter, r *http.Request, opts *AcceptOptions, headers map[string]string) (*WebSocket, error) {
	if opts == nil {
		opts = &accept_options
	}

	if r.Method != "GET" {
		err := errors.New("websocket: websocket request method is not GET")
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return nil, err
	}
	if !HeaderValueCheck(r.Header, "Connection", "upgrade") {
		err := errors.New("websocket: 'upgrade' token not found in 'Connection' header")
		http.Error(w, err.Error(), http.StatusUpgradeRequired)
		return nil, err
	}
	if !HeaderValueCheck(r.Header, "Upgrade", "websocket") {
		err := errors.New("websocket: 'websocket' token not found in 'Upgrade' header")
		http.Error(w, err.Error(), http.StatusUpgradeRequired)
		return nil, err
	}
	if !HeaderValueCheck(r.Header, "Sec-Websocket-Version", "13") {
		err := errors.New("websocket: unsupported websocket version")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	/*
		if !OriginHostCheck(r) {
			err := errors.New("websocket: request origin not allowed")
			http.Error(w, err.Error(), http.StatusForbidden)
			return nil, err
		}*/
	challengeKey := r.Header.Get("Sec-Websocket-Key")
	if challengeKey == "" {
		err := errors.New("websocket: 'Sec-WebSocket-Key' header is missing")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	subproto := opts.selectSubProtocol(getSubProtocols(r))
	var compress bool = false
	for _, ext := range parseHeaderExtensions(r.Header) {
		if ext.name != "permessage-deflate" {
			continue
		}
		compress = true
		break
	}
	if opts.compress_alloter == nil {
		compress = false
	}

	jacker, ok := w.(http.Hijacker)
	if !ok {
		err := errors.New("websocket: response does not implement http.Hijacker")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	var brw *bufio.ReadWriter
	net_cnn, brw, err := jacker.Hijack()
	if err != nil {
		err := errors.New("websocket: response does not implement http.Hijacker")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	if brw.Reader.Buffered() > 0 {
		net_cnn.Close()
		return nil, errors.New("websocket: client sent data before handshake")
	}

	wr := brw.Writer
	wr.Reset(net_cnn)
	ws := newWebSocket(net_cnn, opts, brw.Reader)
	ws.flateWrite = compress
	//ws.opts = opts

	wr.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	wr.WriteString("Upgrade: websocket\r\n")
	wr.WriteString("Connection: Upgrade\r\n")
	wr.WriteString("Sec-WebSocket-Accept: ")
	wr.WriteString(genAcceptKey(challengeKey))
	wr.WriteString("\r\n")
	if subproto != "" {
		wr.WriteString("Sec-WebSocket-Protocol: ")
		wr.WriteString(subproto)
		wr.WriteString("\r\n")
	}
	if compress == true {
		wr.WriteString("Sec-WebSocket-Extensions: ")
		wr.WriteString("permessage-deflate; server_no_context_takeover; client_no_context_takeover")
		wr.WriteString("\r\n")
	}
	for key, val := range headers {
		if strings.Index(key, "\"Sec-WebSocket") == 0 {
			continue
		}
		wr.WriteString(key)
		wr.WriteString(": ")
		wr.WriteString(val)
		wr.WriteString("\r\n")
	}
	wr.WriteString("\r\n")

	net_cnn.SetDeadline(time.Time{})
	if err = wr.Flush(); err != nil {
		net_cnn.Close()
		return nil, err
	}

	return ws, nil
}
