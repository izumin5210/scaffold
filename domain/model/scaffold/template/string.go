package template

import (
	"bytes"
	"reflect"
	"strings"
	gotemplate "text/template"

	"github.com/pkg/errors"

	"github.com/jinzhu/inflection"
	"github.com/serenize/snaker"
)

// String is a compilable string with text/template package
type String string

// Compile generates textual output applied a parsed template to the specified values
func (s String) Compile(name string, v interface{}) (string, error) {
	tmpl := gotemplate.New(name)
	fmap, err := s.createFuncMap(v)
	if err != nil {
		return string(s), err
	}
	tmpl.Funcs(fmap)
	tmpl, err = tmpl.Parse(string(s))
	if err != nil {
		return string(s), errors.Wrapf(err, "Failed to parse a template: %q", string(s))
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, struct{}{})
	if err != nil {
		return string(s), errors.Wrapf(err, "Failed to execute a template: %q", string(s))
	}
	return string(buf.Bytes()), nil
}

var defaultFuncMap = gotemplate.FuncMap{
	"toUpper":     strings.ToUpper,
	"toLower":     strings.ToLower,
	"camelize":    snaker.SnakeToCamelLower,
	"pascalize":   snaker.SnakeToCamel,
	"underscored": snaker.CamelToSnake,
	"dasherize": func(s string) string {
		return strings.Replace(snaker.CamelToSnake(s), "_", "-", -1)
	},
	"singularize": inflection.Singular,
	"pluralize":   inflection.Plural,
	"firstChild":  func(s string) string { return s[:1] },
}

func (s String) createFuncMap(v interface{}) (gotemplate.FuncMap, error) {
	fmap := gotemplate.FuncMap{}

	if v != nil {
		rv := reflect.ValueOf(v)
		for rv.Kind() != reflect.Struct {
			if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
				rv = rv.Elem()
			} else {
				return nil, errors.Errorf("Unsupported type: %v", rv.Kind())
			}
		}
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			ft := rt.Field(i)
			if ft.PkgPath == "" {
				fv := rv.Field(i)
				name := snaker.CamelToSnake(ft.Name)
				switch fv.Kind() {
				case reflect.String:
					fmap[name] = func() string {
						return fv.String()
					}
				default:
					return nil, errors.Errorf("Unsupported type field: {%s: %v}", ft.Name, fv.Kind())
				}
			}
		}
	}

	for k, f := range defaultFuncMap {
		fmap[k] = f
	}

	return fmap, nil
}
