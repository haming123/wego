package klog

import "errors"

var log_default *Klog
func InitKlog(lpath string, rtype RotateType) *Klog {
	Close()
	log_default = NewKlog(lpath, rtype)
	return log_default
}

func GetKlog() *Klog {
	return log_default
}

func Close() {
	if log_default != nil {
		log_default.Close()
	}
}

func SetLogSender(sender LogSender) error {
	if log_default == nil {
		return errors.New("klog is nil")
	}
	return log_default.SetLogSender(sender)
}

func NewL(class_name string) *LogLine {
	line := getLineEnt()
	line.out = log_default
	line.ClassName(class_name)
	return line
}
