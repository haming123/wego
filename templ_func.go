package wego

import (
	"errors"
	"html/template"
	"time"
)

var tmplFuncMap = make(template.FuncMap)
func init() {
	tmplFuncMap["TimeFormat"] = TimeFormat
}

func AddTemplFunc(key string, fn interface{}) error {
	if _, has := tmplFuncMap[key]; has {
		return errors.New("duplicated key")
	}
	tmplFuncMap[key] = fn
	return nil
}

func TimeFormat(t time.Time, layout string) (datestring string) {
	datestring = t.Format(layout)
	return
}
