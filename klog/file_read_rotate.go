package klog

import (
	"errors"
	"os"
	"path/filepath"
)

func get_next_log_file(fpath string, cur_file string) string {
	f_next := ""
	files, _ := filepath.Glob(fpath + "/klog_*.log")
	for _, sname := range files {
		if sname <= cur_file {
			continue
		}

		if f_next == "" {
			f_next = sname
		} else if sname < f_next {
			f_next = sname
		}
	}
	return f_next
}

func open_file_by_pos(pos FilePos) (*os.File, error) {
	if pos.fname == "" {
		return nil, errors.New("file name is null")
	}

	file, err := os.OpenFile(pos.fname, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	if pos.offset > fi.Size() {
		pos.offset = fi.Size()
	}
	file.Seek(pos.offset, os.SEEK_SET)

	loglog.Debugf("Open4Read: %s offset=%d\n", get_short_name(pos.fname), pos.offset)
	return file, nil
}

func fileRead(file *os.File, b []byte) (n int, err error) {
	return file.Read(b)
}