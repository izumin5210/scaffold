package template

import (
	"reflect"
	"testing"

	"github.com/izumin5210/scaffold/domain/model/scaffold/concrete"
)

func Test_NewFile(t *testing.T) {
	const (
		path    = "/app/foo.go"
		content = "package app"
	)
	f := NewFile(path, content).(*entry)

	if actual, expected := f.path, path; actual != expected {
		t.Errorf("NewFile().path is %q, want %q", actual, expected)
	}

	if actual, expected := f.content, content; actual != expected {
		t.Errorf("NewFile().content is %q, want %q", actual, expected)
	}

	if actual, expected := f.dir, false; actual != expected {
		t.Errorf("NewFile().dir is %t, want %t", actual, expected)
	}
}

func Test_NewDir(t *testing.T) {
	const (
		path = "/app/foo.go"
	)
	d := NewDir(path).(*entry)

	if actual, expected := d.path, path; actual != expected {
		t.Errorf("NewDir().path is %q, want %q", actual, expected)
	}

	if actual, expected := d.content, ""; actual != expected {
		t.Errorf("NewDir().content is %q, want %q", actual, expected)
	}

	if actual, expected := d.dir, true; actual != expected {
		t.Errorf("NewDir().dir is %t, want %t", actual, expected)
	}
}

func Test_NewEntry(t *testing.T) {
	cases := []struct {
		path    string
		content string
		dir     bool
	}{
		{
			path:    "/app/.scaffold/foo",
			content: "",
			dir:     true,
		},
		{
			path:    "/app/.scaffold/foo.go",
			content: "package app",
			dir:     false,
		},
	}

	for _, c := range cases {
		e := NewEntry(c.path, c.content, c.dir).(*entry)

		if actual, expected := e.path, c.path; actual != expected {
			t.Errorf("NewEntry().path is %q, want %q", actual, expected)
		}

		if actual, expected := e.content, c.content; actual != expected {
			t.Errorf("NewEntry().content is %q, want %q", actual, expected)
		}

		if actual, expected := e.dir, c.dir; actual != expected {
			t.Errorf("NewEntry().dir is %t, want %t", actual, expected)
		}
	}
}

func Test_Entry_Compile(t *testing.T) {
	cases := []struct {
		in  Entry
		v   interface{}
		out concrete.Entry
	}{
		{
			in:  NewDir("/app/{{name}}"),
			v:   struct{ Name string }{Name: "foobar"},
			out: concrete.NewEntry("/app/foobar", "", true),
		},
		{
			in:  NewFile("/{{namespace}}/{{name}}.go", "package {{namespace}}\n\ntype {{name}} struct{}"),
			v:   struct{ Name, Namespace string }{Name: "foobar", Namespace: "app"},
			out: concrete.NewEntry("/app/foobar.go", "package app\n\ntype foobar struct{}", false),
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

func Test_Entry_Compile_WhenFailedToCompilePath(t *testing.T) {
	d := NewDir("/app/{{name}}")
	e, err := d.Compile(nil)

	if err == nil {
		t.Error("Should return an error")
	}

	if e != nil {
		t.Error("Should not return an entry")
	}
}

func Test_Entry_Compile_WhenFailedToCompileContent(t *testing.T) {
	d := NewFile("/app/{{name}}.go", "package {{name}")
	e, err := d.Compile(struct{ Name string }{Name: "foobar"})

	if err == nil {
		t.Error("Should return an error")
	}

	if e != nil {
		t.Error("Should not return an entry")
	}
}
