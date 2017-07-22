package scaffold

import (
	"testing"
)

func Test_NewDirectory(t *testing.T) {
	// /app/.scaffold/foo
	// ├── bar
	// └── baz
	entries := []Entry{
		&File{Entry: &entry{path: "/app/.scaffold/foo/bar"}},
		&File{Entry: &entry{path: "/app/.scaffold/foo/baz"}},
	}
	d := NewDirectory("/app/.scaffold/foo", entries)

	if actual, expected := len(d.Children()), 2; actual != expected {
		t.Fatalf("New structure got %d items; want %d items", actual, expected)
	}

	for _, key := range []string{"bar", "baz"} {
		if e := d.Children()[key]; e.IsDir() {
			t.Errorf("%s is a directory, want a file", e.Path())
		}
	}
}

func Test_NewDirectory_EmptyEntries(t *testing.T) {
	d := NewDirectory("/app/.scaffold/foo", []Entry{})

	if actual, expected := len(d.Children()), 0; actual != expected {
		t.Errorf("New structure got %d items; want %d items", actual, expected)
	}
}

func Test_NewStructure_Nested(t *testing.T) {
	// /app/.scaffold/foo
	// └── bar
	//     └── baz
	//     │   ├── qux
	//     │   └── quux
	//     └── corge
	entries := []Entry{
		&File{Entry: &entry{path: "/app/.scaffold/foo/bar/baz/qux"}},
		&File{Entry: &entry{path: "/app/.scaffold/foo/bar/baz/quux"}},
		&File{Entry: &entry{path: "/app/.scaffold/foo/bar/corge"}},
		&Directory{Entry: &entry{path: "/app/.scaffold/foo/bar"}},
		&Directory{Entry: &entry{path: "/app/.scaffold/foo/bar/baz"}},
	}
	d := NewDirectory("/app/.scaffold/foo", entries)

	children := d.Children()

	if actual, expected := len(children), 1; actual != expected {
		t.Fatalf("New structure got %d items; want %d items", actual, expected)
	}

	entry := (d.Children()["bar"].(*Directory))

	if !entry.IsDir() {
		t.Fatalf("%s is a file, want a directory", entry.Path())
	}

	if e := entry.Children()["corge"]; e.IsDir() {
		t.Errorf("%s is a directory, want a file", e.Path())
	}

	entry = (entry.Children()["baz"].(*Directory))

	if !entry.IsDir() {
		t.Fatalf("%s is a file, want a directory", entry.Path())
	}

	for _, e := range entry.Children() {
		if e.IsDir() {
			t.Errorf("%s is a directory, want a file", e.Path())
		}
	}
}

func Test_NewDirectory_WithInvalidEntries(t *testing.T) {
	entries := []Entry{
		&Directory{Entry: &entry{path: "/app/.scaffold/foo/bar"}},
		&Directory{Entry: &entry{path: "/app/.scaffold"}},
		&Directory{Entry: &entry{path: "/app/.scaffold/baz/"}},
		&Directory{Entry: &entry{path: "/app/.scaffold/qux/quux"}},
	}
	d := NewDirectory("/app/.scaffold/foo", entries)

	if actual, expected := len(d.Children()), 1; actual != expected {
		t.Fatalf("New structure got %d items; want %d items", actual, expected)
	}

	bar := d.Children()["bar"]

	if !bar.IsDir() {
		t.Errorf("%s is a file, want a directory", bar.Path())
	}

	if actual, expected := len(bar.(*Directory).Children()), 0; actual != expected {
		t.Errorf("%s has %d items; want %d items", bar.Path(), actual, expected)
	}
}
