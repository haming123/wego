package log

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"runtime"
	"sync"
	"time"
)

const tm_flush_loop = 10
type LogWriter interface {
	Write(tm time.Time, data []byte)
	Flush() error
	Close() error
}

type Logger struct {
	mu     sync.Mutex
	level   Level
	caller  bool
	out    	LogWriter
	buf 	*LogBuffer
	close_chan	chan bool
	call_depth 	int
	show_indent bool
}

func NewTermLogger(log_level Level) *Logger {
	var log Logger
	log.out = NewTermWriter()
	log.level = log_level
	log.caller = true
	log.call_depth = 3
	log.buf = NewLogBuffer()

	return &log
}

func NewFileLogger(log_path string, log_level Level) *Logger {
	file_path := get_file_path(log_path)

	var log Logger
	log.out = NewFileWriter(file_path, ROTATE_DAY)
	log.buf = NewLogBuffer()
	log.level = log_level
	log.caller = true
	log.call_depth = 3
	log.close_chan = make(chan bool)
	go log.flush_loop()

	return &log
}

func NewFileLoggerHour(log_path string, log_level Level) *Logger {
	file_path := get_file_path(log_path)

	var log Logger
	log.out = NewFileWriter(file_path, ROTATE_HOUR)
	log.buf = NewLogBuffer()
	log.level = log_level
	log.caller = true
	log.call_depth = 3
	log.close_chan = make(chan bool)
	go log.flush_loop()

	return &log
}

func (l *Logger) SetLevel(lvl Level) {
	l.level = lvl
}

func (l *Logger) ShowCaller(show bool) {
	l.caller = show
}

func (l *Logger) SetCallDepth(call_depth int) {
	l.call_depth = call_depth
}

func (l *Logger) ShowIndent(show bool) {
	l.show_indent = show
}

func (l *Logger) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.out.Flush()
	return nil
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.out.Close()
	if l.close_chan != nil {
		close(l.close_chan)
		l.close_chan = nil
	}

	//fmt.Println("logger Close")
	return nil
}

func (lg *Logger) flush_loop() {
	for {
		select {
		case <-time.After(tm_flush_loop * time.Second):
			//fmt.Println("Flush")
			lg.out.Flush()
		case <-lg.close_chan:
			//fmt.Println("stop flush loop")
			return
		}
	}
}

func (lg *Logger)output(level_str string, msg string, showCaller bool) {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	now := time.Now()
	src_info := ""
	if showCaller {
		_, source, line, _ := runtime.Caller(lg.call_depth)
		src_info = fmt.Sprintf("%s:%d", get_short_name(source), line)
	}

	lg.buf.Reset()
	lg.buf.WriteTimeString(now)
	lg.buf.WriteByte(' ')
	lg.buf.WriteString(level_str)
	if src_info != "" {
		lg.buf.WriteByte(' ')
		lg.buf.WriteString(src_info)
	}
	lg.buf.WriteByte(' ')
	lg.buf.WriteString(msg)
	if len(msg) == 0 || msg[len(msg)-1] != '\n' {
		lg.buf.WriteByte('\n')
	}

	lg.out.Write(now, lg.buf.GetBytes())
}

func (lg *Logger) fatal(v ...interface{}) {
	lvl := LOG_FATAL
	if lg.level >= lvl {
		if len(v) == 1 {
			str, ok := v[0].(string)
			if ok {
				lg.output(lvl.String(), str, lg.caller)
			} else {
				lg.output(lvl.String(), fmt.Sprintln(v...), lg.caller)
			}
		} else {
			lg.output(lvl.String(), fmt.Sprintln(v...), lg.caller)
		}
	}
}

func (lg *Logger) fatalf(format string, v ...interface{}) {
	lvl := LOG_FATAL
	if lg.level >= lvl {
		lg.output(lvl.String(), fmt.Sprintf(format, v...), lg.caller)
	}
}

func (lg *Logger) fatalJSON(v interface{}) {
	var data []byte
	if lg.show_indent {
		data,_ = json.MarshalIndent(v, "", "    ")
	} else {
		data,_ = json.Marshal(v)
	}

	lvl := LOG_FATAL
	if lg.level >= lvl {
		lg.output(lvl.String(), string(data), lg.caller)
	}
}

func (lg *Logger) fatalXML(v interface{}) {
	var data []byte
	if lg.show_indent {
		data,_ = xml.MarshalIndent(v, "", "  ")
	} else {
		data,_ = xml.Marshal(v)
	}

	lvl := LOG_FATAL
	if lg.level >= lvl {
		lg.output(lvl.String(), string(data), lg.caller)
	}
}

func (lg *Logger) error(v ...interface{}) {
	lvl := LOG_ERROR
	if lg.level >= lvl {
		if len(v) == 1 {
			str, ok := v[0].(string)
			if ok {
				lg.output(lvl.String(), str, lg.caller)
			} else {
				lg.output(lvl.String(), fmt.Sprintln(v...), lg.caller)
			}
		} else {
			lg.output(lvl.String(), fmt.Sprintln(v...), lg.caller)
		}
	}
}

func (lg *Logger) errorf(format string, v ...interface{}) {
	lvl := LOG_ERROR
	if lg.level >= lvl {
		lg.output(lvl.String(), fmt.Sprintf(format, v...), lg.caller)
	}
}

func (lg *Logger) errorJSON(v interface{}) {
	var data []byte
	if lg.show_indent {
		data,_ = json.MarshalIndent(v, "", "    ")
	} else {
		data,_ = json.Marshal(v)
	}

	lvl := LOG_ERROR
	if lg.level >= lvl {
		lg.output(lvl.String(), string(data), lg.caller)
	}
}

func (lg *Logger) errorXML(v interface{}) {
	var data []byte
	if lg.show_indent {
		data,_ = xml.MarshalIndent(v, "", "  ")
	} else {
		data,_ = xml.Marshal(v)
	}

	lvl := LOG_ERROR
	if lg.level >= lvl {
		lg.output(lvl.String(), string(data), lg.caller)
	}
}

func (lg *Logger) warn(v ...interface{}) {
	lvl := LOG_WARN
	if lg.level >= lvl {
		if len(v) == 1 {
			str, ok := v[0].(string)
			if ok {
				lg.output(lvl.String(), str, lg.caller)
			} else {
				lg.output(lvl.String(), fmt.Sprintln(v...), lg.caller)
			}
		} else {
			lg.output(lvl.String(), fmt.Sprintln(v...), lg.caller)
		}
	}
}

func (lg *Logger) warnf(format string, v ...interface{}) {
	lvl := LOG_WARN
	if lg.level >= lvl {
		lg.output(lvl.String(), fmt.Sprintf(format, v...), lg.caller)
	}
}

func (lg *Logger) warnJSON(v interface{}) {
	var data []byte
	if lg.show_indent {
		data,_ = json.MarshalIndent(v, "", "    ")
	} else {
		data,_ = json.Marshal(v)
	}

	lvl := LOG_WARN
	if lg.level >= lvl {
		lg.output(lvl.String(), string(data), lg.caller)
	}
}

func (lg *Logger) warnXML(v interface{}) {
	var data []byte
	if lg.show_indent {
		data,_ = xml.MarshalIndent(v, "", "  ")
	} else {
		data,_ = xml.Marshal(v)
	}

	lvl := LOG_WARN
	if lg.level >= lvl {
		lg.output(lvl.String(), string(data), lg.caller)
	}
}

func (lg *Logger) info(v ...interface{}) {
	lvl := LOG_INFO
	if lg.level >= lvl {
		if len(v) == 1 {
			str, ok := v[0].(string)
			if ok {
				lg.output(lvl.String(), str, lg.caller)
			} else {
				lg.output(lvl.String(), fmt.Sprintln(v...), lg.caller)
			}
		} else {
			lg.output(lvl.String(), fmt.Sprintln(v...), lg.caller)
		}
	}
}

func (lg *Logger) infof(format string, v ...interface{}) {
	lvl := LOG_INFO
	if lg.level >= lvl {
		lg.output(lvl.String(), fmt.Sprintf(format, v...), lg.caller)
	}
}

func (lg *Logger) infoJSON(v interface{}) {
	var data []byte
	if lg.show_indent {
		data,_ = json.MarshalIndent(v, "", "    ")
	} else {
		data,_ = json.Marshal(v)
	}

	lvl := LOG_INFO
	if lg.level >= lvl {
		lg.output(lvl.String(), string(data), lg.caller)
	}
}

func (lg *Logger) infoXML(v interface{}) {
	var data []byte
	if lg.show_indent {
		data,_ = xml.MarshalIndent(v, "", "  ")
	} else {
		data,_ = xml.Marshal(v)
	}

	lvl := LOG_INFO
	if lg.level >= lvl {
		lg.output(lvl.String(), string(data), lg.caller)
	}
}

func (lg *Logger) debug(v ...interface{}) {
	lvl := LOG_DEBUG
	if lg.level >= lvl {
		if len(v) == 1 {
			str, ok := v[0].(string)
			if ok {
				lg.output(lvl.String(), str, lg.caller)
			} else {
				lg.output(lvl.String(), fmt.Sprintln(v...), lg.caller)
			}
		} else {
			lg.output(lvl.String(), fmt.Sprintln(v...), lg.caller)
		}
	}
}

func (lg *Logger) debugf(format string, v ...interface{}) {
	lvl := LOG_DEBUG
	if lg.level >= lvl {
		lg.output(lvl.String(), fmt.Sprintf(format, v...), lg.caller)
	}
}

func (lg *Logger) debugJSON(v interface{}) {
	var data []byte
	if lg.show_indent {
		data,_ = json.MarshalIndent(v, "", "    ")
	} else {
		data,_ = json.Marshal(v)
	}

	lvl := LOG_DEBUG
	if lg.level >= lvl {
		lg.output(lvl.String(), string(data), lg.caller)
	}
}

func (lg *Logger) debugXML(v interface{}) {
	var data []byte
	if lg.show_indent {
		data,_ = xml.MarshalIndent(v, "", "  ")
	} else {
		data,_ = xml.Marshal(v)
	}

	lvl := LOG_DEBUG
	if lg.level >= lvl {
		lg.output(lvl.String(), string(data), lg.caller)
	}
}

//用于与Fatal/Error...一致的call_depth
func (lg *Logger)outputWrap(level_str string, msg string){
	lg.output(level_str, msg, lg.caller)
}

//用于与Fatal/Error...一致的call_depth
func (lg *Logger)outputNoCaller(level_str string, msg string){
	lg.output(level_str, msg, false)
}

func (lg *Logger)Output(level_str string, msg string){
	lg.outputWrap(level_str, msg)
}

func (lg *Logger)OutputNoCaller(level_str string, msg string){
	lg.outputNoCaller(level_str, msg)
}

func (lg *Logger)Fatal(v ...interface{}) {
	lg.fatal(v...)
}

func (lg *Logger)Fatalf(format string, v ...interface{}) {
	lg.fatalf(format, v...)
}

func (lg *Logger)FatalJSON(v interface{}) {
	lg.fatalJSON(v)
}

func (lg *Logger)FatalXML(v interface{}) {
	lg.fatalXML(v)
}

func (lg *Logger)Error(v ...interface{}) {
	lg.error(v...)
}

func (lg *Logger)Errorf(format string, v ...interface{}) {
	lg.errorf(format, v...)
}

func (lg *Logger)ErrorJSON(v interface{}) {
	lg.errorJSON(v)
}

func (lg *Logger)ErrorXML(v interface{}) {
	lg.errorXML(v)
}

func (lg *Logger)Warn(v ...interface{}) {
	lg.warn(v...)
}

func (lg *Logger)Warnf(format string, v ...interface{}) {
	lg.warnf(format, v...)
}

func (lg *Logger)WarnJSON(v interface{}) {
	lg.warnJSON(v)
}

func (lg *Logger)WarnXML(v interface{}) {
	lg.warnXML(v)
}

func (lg *Logger)Info(v ...interface{}) {
	lg.info(v...)
}

func (lg *Logger)Infof(format string, v ...interface{}) {
	lg.infof(format, v...)
}

func (lg *Logger)InfoJSON(v interface{}) {
	lg.infoJSON(v)
}

func (lg *Logger)InfoXML(v interface{}) {
	lg.infoXML(v)
}

func (lg *Logger)Debug(v ...interface{}) {
	lg.debug(v...)
}

func (lg *Logger)Debugf(format string, v ...interface{}) {
	lg.debugf(format, v...)
}

func (lg *Logger)DebugJSON(v interface{}) {
	lg.debugJSON(v)
}

func (lg *Logger)DebugXML(v interface{}) {
	lg.debugXML(v)
}
