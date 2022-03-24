package klog

import (
	"os"
)

func get_short_name(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' || file[i] == '\\' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return  file
}

func get_file_path(file_path string) string {
	ch := file_path[len(file_path) -1]
	if  ch != '/' && ch != '\\' {
		file_path += string(os.PathSeparator)
	}
	arr := make([]byte, len(file_path))
	for i:= 0; i < len(file_path) ; i++ {
		ch := file_path[i]
		if  ch == '/' || ch == '\\' {
			ch = os.PathSeparator
		}
		arr[i] = ch
	}
	return string(arr)
}
