package worm

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Level int
const (
	LOG_OFF Level = iota
	LOG_FATAL
	LOG_ERROR
	LOG_WARN
	LOG_INFO
	LOG_DEBUG
)

type SimpleLogger struct {
	mu     	sync.Mutex
	out    	io.Writer
	buf    	[]byte
	level   Level
}

//打印调试日志
type LogPrintCB func(level Level, msg string)

var debug_log *SimpleLogger = NewSimpleLogger()
func SetDebugLogLevel(level Level)  {
	debug_log.SetLevel(level)
}

func print_debug_log(level Level, msg string) {
	if level == LOG_FATAL {
		debug_log.Fatal(msg)
	} else if level == LOG_ERROR {
		debug_log.Error(msg)
	} else if level == LOG_WARN {
		debug_log.Warn(msg)
	} else if level == LOG_INFO {
		debug_log.Info(msg)
	} else if level == LOG_DEBUG {
		debug_log.Debug(msg)
	}
}

func NewSimpleLogger() *SimpleLogger {
	log_ent := SimpleLogger{}
	log_ent.level = LOG_OFF
	log_ent.out = os.Stdout
	return &log_ent
}

func (lg *SimpleLogger) SetLevel(level Level) {
	if level < LOG_OFF {
		level = LOG_OFF
	}
	if level > LOG_DEBUG {
		level = LOG_DEBUG
	}
	lg.level = level
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

func (lg *SimpleLogger) Output(level_str string, msg string) error {
	now := time.Now()
	lg.mu.Lock()
	defer lg.mu.Unlock()

	lg.buf = lg.buf[:0]
	lg.add_time_info(&lg.buf, now)
	lg.buf = append(lg.buf, ' ')
	lg.buf = append(lg.buf, level_str...)
	lg.buf = append(lg.buf, ' ')

	lg.buf = append(lg.buf, msg...)
	if len(msg) == 0 || msg[len(msg)-1] != '\n' {
		lg.buf = append(lg.buf, '\n')
	}

	_, err := lg.out.Write(lg.buf)
	return err
}

func (lg *SimpleLogger) Fatal(v ...interface{}) {
	if lg.level >= LOG_FATAL {
		lg.Output("[F]", fmt.Sprintln(v...))
	}
}

func (lg *SimpleLogger) Fatalf(format string, v ...interface{}) {
	if lg.level >= LOG_FATAL {
		lg.Output("[F]", fmt.Sprintf(format, v...))
	}
}

func (lg *SimpleLogger) Error(v ...interface{}) {
	if lg.level >= LOG_ERROR {
		lg.Output("[E]", fmt.Sprintln(v...))
	}
}

func (lg *SimpleLogger) Errorf(format string, v ...interface{}) {
	if lg.level >= LOG_ERROR {
		lg.Output("[E]", fmt.Sprintf(format, v...))
	}
}

func (lg *SimpleLogger) Warn(v ...interface{}) {
	if lg.level >= LOG_WARN {
		lg.Output("[W]", fmt.Sprintln(v...))
	}
}

func (lg *SimpleLogger) Warnf(format string, v ...interface{}) {
	if lg.level >= LOG_WARN {
		lg.Output("[W]", fmt.Sprintf(format, v...))
	}
}

func (lg *SimpleLogger) Info(v ...interface{}) {
	if lg.level >= LOG_INFO {
		lg.Output("[I]", fmt.Sprintln(v...))
	}
}

func (lg *SimpleLogger) Infof(format string, v ...interface{}) {
	if lg.level >= LOG_INFO {
		lg.Output("[I]", fmt.Sprintf(format, v...))
	}
}

func (lg *SimpleLogger) Debug(v ...interface{}) {
	if lg.level >= LOG_DEBUG {
		lg.Output("[D]", fmt.Sprintln(v...))
	}
}

func (lg *SimpleLogger) Debugf(format string, v ...interface{}) {
	if lg.level >= LOG_DEBUG {
		lg.Output("[D]", fmt.Sprintf(format, v...))
	}
}
