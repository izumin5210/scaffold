package scaffold

import (
	"bytes"
	gotemplate "text/template"
)

// Template is a template renderer for file paths and contents
type Template interface {
	Compile(text string) (string, error)
}

type template struct {
	name    string
	funcMap gotemplate.FuncMap
}

// NewTemplate creates a new template instance
func NewTemplate(name string) Template {
	tmpl := &template{name: name}
	tmpl.funcMap = tmpl.createFuncMap()
	return tmpl
}

func (t *template) Compile(text string) (string, error) {
	tmpl, err := t.getInstance().Parse(text)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, struct{}{})
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}

func (t *template) getInstance() *gotemplate.Template {
	tmpl := gotemplate.New(t.name)
	tmpl.Funcs(t.funcMap)
	return tmpl
}

func (t *template) createFuncMap() gotemplate.FuncMap {
	return gotemplate.FuncMap{
		"name": func() string { return t.name },
	}
}
