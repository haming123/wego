package wego

import (
	"fmt"
	"html/template"
	"path/filepath"
	"sync"
)

type WebTemplates struct {
	mux   		sync.Mutex
	tpl_map		map[string]*template.Template
	delim_left 	string
	delim_right string
}

func (this *WebTemplates) Init()  {
	this.tpl_map = make(map[string]*template.Template)
	this.delim_left  = "{{"
	this.delim_right = "}}"
}

func (this *WebTemplates) SetDelim(left, right string){
	this.delim_left = left
	this.delim_right = right
}

func (this *WebTemplates) getTemplate(filenames ...string) (*template.Template, error) {
	this.mux.Lock()
	defer this.mux.Unlock()

	if len(filenames) == 0 {
		return nil, fmt.Errorf("must have one file name")
	}

	key := ""
	for _, item := range filenames {
		key += item
	}
	tpl, ok := this.tpl_map[key]
	if ok {
		//fmt.Println("use template cache")
		return tpl, nil
	}

	//fmt.Println("load template")
	name := filepath.Base(filenames[0])
	t, err := template.New(name).Delims(this.delim_left, this.delim_right).Funcs(tmplFuncMap).ParseFiles(filenames...)
	if err == nil {
		this.tpl_map[key] = t
	}
	return t, err
}
