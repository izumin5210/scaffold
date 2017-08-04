package scaffold

import (
	"testing"
)

func Test_NewConcreteEntry(t *testing.T) {
	cases := []struct {
		path     string
		content  string
		dir      bool
		tmplPath string
	}{
		{
			path:     "/app",
			content:  "",
			dir:      true,
			tmplPath: "/{{name}}",
		},
		{
			path:     "/app/foo.go",
			content:  "package app",
			dir:      false,
			tmplPath: "/{{name}}/foo.go",
		},
	}

	for _, c := range cases {
		e := NewConcreteEntry(c.path, c.content, c.dir, c.tmplPath)

		if actual, expected := e.Path(), c.path; actual != expected {
			t.Errorf("Path() returns %q, want %q", actual, expected)
		}

		if actual, expected := e.Content(), c.content; actual != expected {
			t.Errorf("Content() returns %q, want %q", actual, expected)
		}

		if actual, expected := e.IsDir(), c.dir; actual != expected {
			t.Errorf("IsDir() returns %t, want %t", actual, expected)
		}

		if actual, expected := e.TemplatePath(), c.tmplPath; actual != expected {
			t.Errorf("TemplatePath() returns %q, want %q", actual, expected)
		}
	}
}

func Test_NewConcreteFile(t *testing.T) {
	path := "/app/foo.go"
	content := "package app"
	tmplPath := "/{{name}}/foo.go"
	e := NewConcreteFile(path, content, tmplPath)

	if actual, expected := e.Path(), path; actual != expected {
		t.Errorf("Path() returns %q, want %q", actual, expected)
	}

	if actual, expected := e.Content(), content; actual != expected {
		t.Errorf("Content() returns %q, want %q", actual, expected)
	}

	if actual, expected := e.IsDir(), false; actual != expected {
		t.Errorf("IsDir() returns %t, want %t", actual, expected)
	}

	if actual, expected := e.TemplatePath(), tmplPath; actual != expected {
		t.Errorf("TemplatePath() returns %q, want %q", actual, expected)
	}
}

func Test_NewConcreteDir(t *testing.T) {
	path := "/app/foo.go"
	tmplPath := "/{{name}}/foo.go"
	e := NewConcreteDir(path, tmplPath)

	if actual, expected := e.Path(), path; actual != expected {
		t.Errorf("Path() returns %q, want %q", actual, expected)
	}

	if actual, expected := e.Content(), ""; actual != expected {
		t.Errorf("Content() returns %q, want %q", actual, expected)
	}

	if actual, expected := e.IsDir(), true; actual != expected {
		t.Errorf("IsDir() returns %t, want %t", actual, expected)
	}

	if actual, expected := e.TemplatePath(), tmplPath; actual != expected {
		t.Errorf("TemplatePath() returns %q, want %q", actual, expected)
	}
}
