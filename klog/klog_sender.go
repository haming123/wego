package klog

import "strings"

const tm_send_loop = 30
type LogSender interface {
	SetLogger(logger *Klog)
	BeginSendLoop()
}

type FilePos struct {
	fname 		string
	offset 		int64
	flag_end 	bool
}

func (fpos *FilePos)GetFileCode() string {
	fname := get_short_name(fpos.fname)
	fname = strings.TrimRight(fname, ".log")
	strs := strings.Split(fname, "_")
	if len(strs) == 2 {
		return strs[1]
	}
	return ""
}