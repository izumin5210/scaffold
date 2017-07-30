package scaffold

import (
	"reflect"
	"testing"
)

func Test_NewTemplateFile(t *testing.T) {
	const (
		path    = "/app/foo.go"
		content = "package app"
	)
	f := NewTemplateFile(path, content)

	if actual, expected := f.Path(), path; actual != expected {
		t.Errorf("Path() returns %q, want %q", actual, expected)
	}

	if actual, expected := f.Content(), content; actual != expected {
		t.Errorf("Content() returns %q, want %q", actual, expected)
	}

	if actual, expected := f.IsDir(), false; actual != expected {
		t.Errorf("IsDir() returns %t, want %t", actual, expected)
	}
}

func Test_NewTemplateDir(t *testing.T) {
	const (
		path = "/app/foo.go"
	)
	d := NewTemplateDir(path)

	if actual, expected := d.Path(), path; actual != expected {
		t.Errorf("Path() returns %q, want %q", actual, expected)
	}

	if actual, expected := d.Content(), ""; actual != expected {
		t.Errorf("Content() returns %q, want %q", actual, expected)
	}

	if actual, expected := d.IsDir(), true; actual != expected {
		t.Errorf("IsDir() returns %t, want %t", actual, expected)
	}
}

func Test_NewTemplateEntry(t *testing.T) {
	cases := []struct {
		path    string
		content string
		dir     bool
	}{
		{
			path:    "/app",
			content: "",
			dir:     true,
		},
		{
			path:    "/app/foo.go",
			content: "package app",
			dir:     false,
		},
	}

	for _, c := range cases {
		e := NewEntry(c.path, c.content, c.dir)

		if actual, expected := e.Path(), c.path; actual != expected {
			t.Errorf("Path() returns %q, want %q", actual, expected)
		}

		if actual, expected := e.Content(), c.content; actual != expected {
			t.Errorf("Content() returns %q, want %q", actual, expected)
		}

		if actual, expected := e.IsDir(), c.dir; actual != expected {
			t.Errorf("IsDir() returns %t, want %t", actual, expected)
		}
	}
}

func Test_TemplateEntry_Compile(t *testing.T) {
	cases := []struct {
		in  TemplateEntry
		v   interface{}
		out Entry
	}{
		{
			in:  NewTemplateDir("/app/{{name}}"),
			v:   struct{ Name string }{Name: "foobar"},
			out: NewEntry("/app/foobar", "", true),
		},
		{
			in:  NewTemplateFile("/app/{{name}}.go", "package {{namespace}}\n\ntype {{name}} struct{}"),
			v:   struct{ Name, Namespace string }{Name: "foobar", Namespace: "app"},
			out: NewEntry("/app/foobar.go", "package app\n\ntype foobar struct{}", false),
		},
	}

	for _, c := range cases {
		e, err := c.in.Compile(c.v)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if actual, expected := e, c.out; !reflect.DeepEqual(actual, expected) {
			t.Errorf("Compile() returns %v, want %v", actual, expected)
		}
	}
}

func Test_TemplateEntry_Compile_WhenFailedToCompilePath(t *testing.T) {
	d := NewTemplateDir("/app/{{name}}")
	e, err := d.Compile(nil)

	if err == nil {
		t.Error("Should return an error")
	}

	if e != nil {
		t.Error("Should not return an entry")
	}
}

func Test_TemplateEntry_Compile_WhenFailedToCompileContent(t *testing.T) {
	d := NewTemplateFile("/app/{{name}}.go", "package {{name}")
	e, err := d.Compile(struct{ Name string }{Name: "foobar"})

	if err == nil {
		t.Error("Should return an error")
	}

	if e != nil {
		t.Error("Should not return an entry")
	}
}
