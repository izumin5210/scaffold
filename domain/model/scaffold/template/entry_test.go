package template

import (
	"testing"
)

func Test_NewFile(t *testing.T) {
	const (
		path    String = "/app/.scaffold/foo.go"
		content String = "package app"
	)
	f := NewFile(path, content)

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

func Test_NewDir(t *testing.T) {
	const (
		path String = "/app/foo.go"
	)
	d := NewDir(path)

	if actual, expected := d.Path(), path; actual != expected {
		t.Errorf("Path() returns %q, want %q", actual, expected)
	}

	if actual, expected := d.Content(), String(""); actual != expected {
		t.Errorf("Content() returns %q, want %q", actual, expected)
	}

	if actual, expected := d.IsDir(), true; actual != expected {
		t.Errorf("IsDir() returns %t, want %t", actual, expected)
	}
}

func Test_NewEntry(t *testing.T) {
	cases := []struct {
		path    String
		content String
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
		e := NewEntry(String(c.path), String(c.content), c.dir)

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

// func Test_Entry_Compile(t *testing.T) {
// 	root := "/app"
// 	tmplRoot := filepath.Join(root, ".scaffold")
// 	getConcPath := func(path string) string {
// 		return filepath.Join(root, path)
// 	}
// 	getTmplPath := func(path string) String {
// 		return String(filepath.Join(tmplRoot, path))
// 	}
// 	cases := []struct {
// 		in  Entry
// 		v   interface{}
// 		out Entry
// 	}{
// 		{
// 			in:  NewTemplateDir(getTmplPath("{{name}}"), tmplRoot),
// 			v:   struct{ Name string }{Name: "foobar"},
// 			out: NewEntry(getConcPath("foobar"), "", true),
// 		},
// 		{
// 			in:  NewTemplateFile(getTmplPath("{{name}}.go"), "package {{namespace}}\n\ntype {{name}} struct{}", tmplRoot),
// 			v:   struct{ Name, Namespace string }{Name: "foobar", Namespace: "app"},
// 			out: NewEntry(getConcPath("foobar.go"), "package app\n\ntype foobar struct{}", false),
// 		},
// 	}

// 	for _, c := range cases {
// 		e, err := c.in.Compile(root, c.v)

// 		if err != nil {
// 			t.Errorf("Unexpected error %v", err)
// 		}

// 		if actual, expected := e, c.out; !reflect.DeepEqual(actual, expected) {
// 			t.Errorf("Compile() returns %v, want %v", actual, expected)
// 		}
// 	}
// }

// func Test_Entry_Compile_WhenFailedToCompilePath(t *testing.T) {
// 	d := NewTemplateDir("/app/.scaffold/{{name}}", "/app/.scaffold")
// 	e, err := d.Compile("/app", nil)

// 	if err == nil {
// 		t.Error("Should return an error")
// 	}

// 	if e != nil {
// 		t.Error("Should not return an entry")
// 	}
// }

// func Test_Entry_Compile_WhenFailedToCompileContent(t *testing.T) {
// 	d := NewTemplateFile("/app/.scaffold/{{name}}.go", "package {{name}", "/app/.scaffold")
// 	e, err := d.Compile("/app", struct{ Name string }{Name: "foobar"})

// 	if err == nil {
// 		t.Error("Should return an error")
// 	}

// 	if e != nil {
// 		t.Error("Should not return an entry")
// 	}
// }
