package fs

import (
	"fmt"
	"testing"
)

func Test_NewFile(t *testing.T) {
	p := "/app/foo.go"
	f, err := NewFile(p)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if f.IsDir() {
		t.Error("IsDir() returns true, want false")
	}

	if actual, expected := f.Path(), p; actual != expected {
		t.Errorf(".Path() returns %q, want %q", actual, expected)
	}
}

func Test_NewFile_WithRelPath(t *testing.T) {
	f, err := NewFile("app/foo.go")

	if err == nil {
		t.Error("NewFile() with relative path should return an error")
	}

	if f != nil {
		t.Error("NewFile() with relative path should not return an entry")
	}
}

func Test_NewDir(t *testing.T) {
	p := "/app/foo"
	f, err := NewDir(p)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if !f.IsDir() {
		t.Errorf("EmptyDir(%s).IsDir() returns false, want true", p)
	}

	if actual, expected := f.Path(), p; actual != expected {
		t.Errorf("EmptyDir().Path() returns %q, want %q", actual, expected)
	}
}

func Test_NewDir_WithRelPath(t *testing.T) {
	f, err := NewDir("app/foo")

	if err == nil {
		t.Error("EmptyDir() with relative path should return an error")
	}

	if f != nil {
		t.Error("NewFile() with relative path should not return an entry")
	}
}

func Test_NewDir_WhenEndsWithSlash(t *testing.T) {
	p := "/app/foo"
	f, err := NewDir(fmt.Sprintf("%s/", p))

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if actual, expected := f.Path(), p; actual != expected {
		t.Errorf("EmptyDir().Path() returns %q, want %q", actual, expected)
	}
}

func Test_Entry_BaseName(t *testing.T) {
	testcases := []struct{ path, name string }{
		{path: "/app", name: "app"},
		{path: "/app/foo", name: "foo"},
		{path: "/app/bar/", name: "bar"},
		{path: "/app/baz.go", name: "baz.go"},
	}

	for _, tc := range testcases {
		d, err := NewDir(tc.path)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if actual, expected := d.BaseName(), tc.name; actual != expected {
			t.Errorf("Name() returns %q, want %q", actual, expected)
		}
	}
}

func Test_Entry_DirName(t *testing.T) {
	testcases := []struct{ path, name string }{
		{path: "/app", name: "/"},
		{path: "/app/foo", name: "/app"},
		{path: "/app/foo/bar", name: "/app/foo"},
		{path: "/app/foo/bar/", name: "/app/foo"},
		{path: "/app/baz/qux/quux.go", name: "/app/baz/qux"},
	}

	for _, tc := range testcases {
		d, err := NewDir(tc.path)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if actual, expected := d.DirName(), tc.name; actual != expected {
			t.Errorf("DirName() returns %q, want %q", actual, expected)
		}
	}
}
