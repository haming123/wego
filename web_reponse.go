package wego

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"sync"
)

type WebResponse struct {
	http.ResponseWriter
	StatusCode 	int
	wroteHeader bool
	gzip_buff 	bytes.Buffer
	gzip_flag 	bool
	gzip_size 	int64
}

func (w *WebResponse) reset() {
	w.ResponseWriter = nil
	w.StatusCode = http.StatusOK
	w.wroteHeader = false
	w.gzip_buff.Reset()
	w.gzip_flag = false
	w.gzip_size = 0
}

func (w *WebResponse) GetStatus() int {
	return w.StatusCode
}

func (w *WebResponse) SetStatus(code int) {
	if w.wroteHeader == false {
		w.StatusCode = code
	}
}

func (w *WebResponse) WroteHeader() bool {
	return w.wroteHeader
}

func (w *WebResponse) WriteHeader(code int) {
	w.StatusCode = code
	if w.wroteHeader == false && w.gzip_flag == false {
		w.ResponseWriter.WriteHeader(code)
		w.wroteHeader = true
	}
}

func (w *WebResponse) Write(data []byte) (int, error) {
	if w.wroteHeader == false && w.gzip_flag == false {
		w.ResponseWriter.WriteHeader(w.StatusCode)
		w.wroteHeader = true
	}

	if data != nil && len(data) > 0 {
		if w.gzip_flag == true {
			return w.gzip_buff.Write(data)
		} else {
			return w.ResponseWriter.Write(data)
		}
	}

	return 0, nil
}

var gzPool = sync.Pool{
	New: func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, gzip.BestSpeed)
		if err != nil {
			panic(err)
		}
		return gz
	},
}

func (w *WebResponse) Flush() error {
	//若没有调用WriteHeader()，则补充WriteHeader的调用
	if w.wroteHeader == false && w.gzip_flag == false {
		w.ResponseWriter.WriteHeader(w.StatusCode)
		w.wroteHeader = true
	}

	if w.gzip_flag == true {
		return w.gzipFlush()
	}

	return nil
}

func (w *WebResponse) gzipFlush() error {
	min_size := w.gzip_size
	if min_size < 256 {
		min_size = 256
	}

	//数据量<min_size则不进行压缩
	data_len := w.gzip_buff.Len()
	if int64(data_len) < min_size {
		var err error
		w.ResponseWriter.WriteHeader(w.StatusCode)
		if data_len > 0 {
			_, err = w.ResponseWriter.Write(w.gzip_buff.Bytes())
		}
		return err
	}

	w.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	w.ResponseWriter.Header().Set("Vary", "Accept-Encoding")
	w.ResponseWriter.WriteHeader(w.StatusCode)

	gz := gzPool.Get().(*gzip.Writer)
	gz.Reset(w.ResponseWriter)

	gz.Write(w.gzip_buff.Bytes())
	gz.Close()

	gz.Reset(ioutil.Discard)
	gzPool.Put(gz)
	w.gzip_buff.Reset()
	return nil
}


