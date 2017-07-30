package scaffold

import "testing"

func Test_NewEntry(t *testing.T) {
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
