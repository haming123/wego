package log

import "testing"

func Test_get_short_name(t *testing.T)  {
	file_path := "D:\\GoApp\\src\\log4go\\main\\logs\\applog_20211129.log"
	file_name := get_short_name(file_path)
	if file_name != "applog_20211129.log" {
		t.Errorf("get_short_name return: %s", file_name)
	}
}
