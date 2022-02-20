package wini

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type ConfigItem struct {
	section	 string
	name	 string
	value	 string
}

func (this *ConfigItem) reset() {
	this.section = ""
	this.name = ""
	this.value = ""
}

/*
在 ini 文件中，每个键值对占用一行，中间使用=隔开。以#开头的内容为注释。
ini 文件是以分区（section）组织的。
分区以[name]开始，在下一个分区前结束。所有分区前的内容属于默认分区
*/
func ParseFile(file_path string, cfg *ConfigData) error {
	err := parseFile(file_path, cfg)
	if err != nil {
		return err
	}
	return nil
}

func parseFile(fileName string, cfg *ConfigData) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	var item ConfigItem
	section := ""
	buf := bufio.NewReader(f)
	for {
		item.reset()
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		line = strings.TrimSpace(line)
		err = parseLine(line, &item)
		if err != nil {
			return errors.New(err.Error() + ":" + line)
		}

		if item.section != "" {
			section = item.section
		} else if item.name != "" {
			cfg.SetData(section, item.name, item.value)
		}
	}

	return nil
}

func parseLine(line string, item *ConfigItem) error {
	cnum := len(line)
	index := strings.Index(line, "#")
	if index >= 0 {
		line = line[0:index]
	}
	index = strings.Index(line, ";")
	if index >= 0 {
		line = line[0:index]
	}

	line = strings.TrimSpace(line)
	cnum = len(line)
	if cnum < 1 {
		return nil
	}

	beg := strings.Index(line, "[")
	if beg >= 0 {
		line = line[1:]
		end := strings.Index(line, "]")
		if end < 1 {
			return errors.New("invalid line")
		}
		item.section = line[0:end]
		return nil
	}

	name, value, ok := SplitAndTrim(line, "=")
	if !ok {
		return errors.New("invalid line")
	}
	value = strings.Trim(value, "\"")

	vLen := len(value)
	if vLen > 3 && value[0] == '$' && value[1] == '{' && value[vLen-1] == '}' {
		value = value[2:vLen-1]
		value = os.Getenv(value)
	}

	item.name = name
	item.value = value
	return nil
}

