package gows

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type DebugLogger interface {
	Output(prefix string, msg string)
}

var debug_log DebugLogger = NewSimpleLogger()
var show_log bool = false
var log_prefix string = "[ws]"

func SetDebugLogger(logger DebugLogger) {
	debug_log = logger
}

func ShowDebugLog(show bool) {
	show_log = show
}

func SetLoggerPrefix(prefix string) {
	log_prefix = prefix
}

func logPrint(v ...interface{}) {
	if show_log == false {
		return
	}
	debug_log.Output(log_prefix, fmt.Sprintln(v...))
}

func logPrintf(format string, v ...interface{}) {
	if show_log == false {
		return
	}
	debug_log.Output(log_prefix, fmt.Sprintf(format, v...))
}

func logPrint4ws(ws *WebSocket, v ...interface{}) {
	if show_log == false {
		return
	}
	logPrintf(ws.RemoteAddr().String() + " " + fmt.Sprintln(v...))
}

func logPrintf4ws(ws *WebSocket, format string, v ...interface{}) {
	if show_log == false {
		return
	}
	logPrintf(ws.RemoteAddr().String() + " " + fmt.Sprintf(format, v...))
}

type SimpleLogger struct {
	mu  sync.Mutex
	out io.Writer
	buf []byte
}

func NewSimpleLogger() *SimpleLogger {
	log_ent := SimpleLogger{}
	log_ent.out = os.Stdout
	return &log_ent
}

func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (lg *SimpleLogger) add_time_info(buf *[]byte, t time.Time) {
	//date
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')

	//time
	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
}

func (lg *SimpleLogger) Output(prefix string, msg string) {
	now := time.Now()
	lg.mu.Lock()
	defer lg.mu.Unlock()

	lg.buf = lg.buf[:0]
	lg.add_time_info(&lg.buf, now)
	lg.buf = append(lg.buf, ' ')
	if len(prefix) > 0 {
		lg.buf = append(lg.buf, prefix...)
		lg.buf = append(lg.buf, ' ')
	}

	lg.buf = append(lg.buf, msg...)
	if len(msg) == 0 || msg[len(msg)-1] != '\n' {
		lg.buf = append(lg.buf, '\n')
	}

	lg.out.Write(lg.buf)
}
