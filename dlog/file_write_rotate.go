package log

import (
	"fmt"
	"os"
	"time"
)

var (
	// permission to:  owner      group      other
	//                 /```\      /```\      /```\
	// octal:            6          6          6
	// binary:         1 1 0      1 1 0      1 1 0
	// what to permit: r w x      r w x      r w x
	// binary         - 1: enabled, 0: disabled
	// what to permit - r: read, w: write, x: execute
	// permission to  - owner: the user that create the file/folder
	//                  group: the users from group that owner is member
	//                  other: all other users
	// DefaultFileMode is the default log file and directory permission.
	DefaultFileMode = os.FileMode(0660)
)

type RotateType int
const (
	ROTATE_DAY RotateType = iota
	ROTATE_HOUR
)

var logChanLength = 32
var dayToSecs int64 = 86400
var hourToSecs int64 = 3600

type FileWriter struct {
	rec 		chan []byte
	file_path 	string
	file     	*BufferWriter
	rotate		RotateType
	tm_file 	int64
}

func NewFileWriter(file_path string, rotate RotateType) *FileWriter {
	fw := &FileWriter{
		rec:		make(chan []byte, logChanLength),
		file_path:	file_path,
		rotate:		rotate,
	}
	return fw
}

func gen_file_name(tm time.Time, file_path string, rotate RotateType) string {
	file_name := file_path + "applog_" + tm.Format("20060102") + "00.log"
	if rotate == ROTATE_HOUR {
		file_name = file_path + "applog_" + tm.Format("2006010215") + ".log"
	}
	return file_name
}

func create_log_file(tm time.Time, file_path string, rotate RotateType) *BufferWriter {
	file := gen_file_name(tm, file_path, rotate)
	logFile, err := NewBufferWriter(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return logFile
}

func (fw *FileWriter) get_cur_log_file(tm time.Time) *BufferWriter {
	dd := tm.Unix() - fw.tm_file
	if fw.file == nil || fw.rotate == ROTATE_DAY && dd > dayToSecs || fw.rotate == ROTATE_HOUR && dd > hourToSecs {
		fw.close_log_file()
		fw.file = create_log_file(tm, fw.file_path, fw.rotate)
		if fw.rotate == ROTATE_HOUR {
			fw.tm_file = tm.Unix() / hourToSecs * hourToSecs
		} else {
			fw.tm_file = tm.Unix() / dayToSecs * dayToSecs
		}
		return fw.file
	}
	return fw.file
}

func (fw *FileWriter)write_log(tm time.Time, data []byte) {
	file := fw.get_cur_log_file(tm)
	_, err := file.Write([]byte(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Print("write log: " + string(data))
}

func (fw *FileWriter) Write(tm time.Time, data []byte)  {
	fw.write_log(tm, data)
}

func (fw *FileWriter) Flush() error {
	if fw.file == nil {
		return nil
	}

	tm := time.Now()
	dd := tm.Unix() - fw.tm_file
	if fw.rotate == ROTATE_DAY && dd > dayToSecs || fw.rotate == ROTATE_HOUR && dd > hourToSecs {
		fw.close_log_file()
		if fw.rotate == ROTATE_HOUR {
			fw.tm_file = tm.Unix() / hourToSecs * hourToSecs
		} else {
			fw.tm_file = tm.Unix() / dayToSecs * dayToSecs
		}
	} else {
		fw.file.Flush()
	}

	//fmt.Println("flush log")
	return nil
}

func (fw *FileWriter) close_log_file()  {
	if fw.file != nil {
		//fmt.Println("close log file: " + fw.file.Name())
		fw.file.Close()
		fw.file = nil
	}
}

func (fw *FileWriter) Close() error {
	fw.close_log_file()
	return nil
}

