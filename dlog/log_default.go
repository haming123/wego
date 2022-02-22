package log

import "encoding/json"

var log_default *Logger
func init() {
	log_default = NewTermLogger(LOG_DEBUG)
	log_default.ShowCaller(true)
}

func InitTermLogger(log_level Level) *Logger {
	log_default.Close()
	log_new := NewTermLogger(log_level)
	log_new.show_indent = log_default.show_indent
	log_new.caller = log_default.caller
	log_default = log_new
	return log_default
}

func InitFileLogger(log_path string, log_level Level) *Logger {
	log_default.Close()
	log_new := NewFileLogger(log_path, log_level)
	log_new.show_indent = log_default.show_indent
	log_new.caller = log_default.caller
	log_default = log_new
	return log_default
}

func InitFileLoggerHour(log_path string, log_level Level) *Logger {
	log_default.Close()
	log_new := NewFileLoggerHour(log_path, log_level)
	log_new.show_indent = log_default.show_indent
	log_new.caller = log_default.caller
	log_default = log_new
	return log_default
}

func GetLogger() *Logger {
	return log_default
}

func Flush() {
	if log_default != nil {
		log_default.Flush()
	}
}

func Close() {
	if log_default != nil {
		log_default.Close()
	}
}

func SetLevel(log_level Level) {
	log_default.SetLevel(log_level)
}

func ShowCaller(show bool) {
	log_default.ShowCaller(show)
}

func ShowIndent(show bool) {
	log_default.show_indent = show
}

func SetCallDepth(call_depth int) {
	log_default.call_depth = call_depth
}

func Output(level_str string, msg string){
	log_default.outputWrap(level_str, msg)
}

func OutputNoCaller(level_str string, msg string){
	log_default.outputNoCaller(level_str, msg)
}

func Fatal(v ...interface{}) {
	log_default.fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	log_default.fatalf(format, v...)
}

func FatalJSON(v interface{}) {
	log_default.fatalJSON(v)
}

func FatalXML(v interface{}) {
	log_default.fatalXML(v)
}

func Error(v ...interface{}) {
	log_default.error(v...)
}

func Errorf(format string, v ...interface{}) {
	log_default.errorf(format, v...)
}

func ErrorJSON(v interface{}) {
	log_default.errorJSON(v)
}

func ErrorXML(v interface{}) {
	log_default.errorXML(v)
}

func Warn(v ...interface{}) {
	log_default.warn(v...)
}

func Warnf(format string, v ...interface{}) {
	log_default.warnf(format, v...)
}

func WarnJSON(v interface{}) {
	log_default.warnJSON(v)
}

func WarnXML(v interface{}) {
	log_default.warnXML(v)
}

func Info(v ...interface{}) {
	log_default.info(v...)
}

func Infof(format string, v ...interface{}) {
	log_default.infof(format, v...)
}

func InfoJSON(v interface{}) {
	log_default.infoJSON(v)
}

func InfoXML(v interface{}) {
	log_default.infoXML(v)
}

func Debug(v ...interface{}) {
	log_default.debug(v...)
}

func Debugf(format string, v ...interface{}) {
	log_default.debugf(format, v...)
}

func  DebugJSON(v interface{}) {
	log_default.debugJSON(v)
}

func  DebugXML(v interface{}) {
	log_default.debugXML(v)
}

func JsonMarshal(ent interface{}) string{
	data,_ := json.Marshal(ent)
	return string(data)
}

func IndentedJSONMarshal(ent interface{}) string{
	data,_:= json.MarshalIndent(ent, "", "    ")
	return string(data)
}
