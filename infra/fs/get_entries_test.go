package fs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func Test_GetEntries(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	rootPath := "/app"
	// app
	// ├── foo
	// │   ├── qux
	// │   │   └── corge.go
	// │   └── qux
	// │       └── grault.go
	// ├── bar.go
	// └── baz
	//     ├── garply
	//     │   └── waldo
	//     ├── garply.go
	//     └── garply_test.go
	entries := []*struct {
		path string
		dir  bool
	}{
		{path: filepath.Join(rootPath, "foo", "qux"), dir: true},
		{path: filepath.Join(rootPath, "foo", "qux", "corge.go"), dir: false},
		{path: filepath.Join(rootPath, "foo", "quux"), dir: true},
		{path: filepath.Join(rootPath, "foo", "quux", "grault.go"), dir: false},
		{path: filepath.Join(rootPath, "bar.go"), dir: false},
		{path: filepath.Join(rootPath, "baz", "garply.go"), dir: false},
		{path: filepath.Join(rootPath, "baz", "garply_test.go"), dir: false},
		{path: filepath.Join(rootPath, "baz", "garply", "waldo"), dir: true},
	}

	for _, e := range entries {
		if e.dir {
			afs.MkdirAll(e.path, 0755)
		} else {
			afs.WriteFile(e.path, []byte("awesome file"), 0644)
		}
	}

	type test struct {
		expectedList []Entry
		compact      bool
	}

	tests := []*test{
		{
			expectedList: []Entry{
				&entry{path: "/app/bar.go", dir: false},
				&entry{path: "/app/baz", dir: true},
				&entry{path: "/app/baz/garply", dir: true},
				&entry{path: "/app/baz/garply.go", dir: false},
				&entry{path: "/app/baz/garply/waldo", dir: true},
				&entry{path: "/app/baz/garply_test.go", dir: false},
				&entry{path: "/app/foo", dir: true},
				&entry{path: "/app/foo/quux", dir: true},
				&entry{path: "/app/foo/quux/grault.go", dir: false},
				&entry{path: "/app/foo/qux", dir: true},
				&entry{path: "/app/foo/qux/corge.go", dir: false},
			},
			compact: false,
		},
		{
			expectedList: []Entry{
				&entry{path: "/app/bar.go", dir: false},
				&entry{path: "/app/baz/garply.go", dir: false},
				&entry{path: "/app/baz/garply/waldo", dir: true},
				&entry{path: "/app/baz/garply_test.go", dir: false},
				&entry{path: "/app/foo/quux/grault.go", dir: false},
				&entry{path: "/app/foo/qux/corge.go", dir: false},
			},
			compact: true,
		},
	}

	for _, ts := range tests {
		entries, err := fs.GetEntries(rootPath, ts.compact)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if actual, expected := len(entries), len(ts.expectedList); actual != expected {
			t.Errorf("GetEntries() returns %d items, want %d items", actual, expected)
		}

		for i, expected := range ts.expectedList {
			if actual := entries[i]; !reflect.DeepEqual(actual, expected) {
				t.Errorf("GetEntries() returns %v, but expected %v", actual, expected)
			}
		}
	}
}

func Test_GetEntries_WhenPathDoesNotExist(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	entries, err := fs.GetEntries("/app", false)

	if err == nil {
		t.Errorf("GetEntries() should return an error")
	}

	if entries != nil {
		t.Errorf("GetEntries() should return %v, want nil", entries)
	}
}

func Test_GetEntries_WhenPathPointsFile(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	afs.WriteFile("/app", []byte("awesome file"), 0644)
	entries, err := fs.GetEntries("/app", false)

	if err == nil {
		t.Errorf("GetEntrie() should return an error")
	}

	if entries != nil {
		t.Errorf("GetEntries() should return %v, want nil", entries)
	}
}

func Test_GetEntries_WithRelPath(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	afs.Mkdir("app/foo", os.ModeDir)
	entries, err := fs.GetEntries("app", false)

	if err == nil {
		t.Errorf("GetEntrie() should return an error")
	}

	if entries != nil {
		t.Errorf("GetEntries() should return %v, want nil", entries)
	}
}
