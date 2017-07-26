package fs

import (
	"os"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func Test_GetDirs(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	afs.Mkdir("/app/foo", os.ModeDir)
	afs.Mkdir("/app/foo/bar", os.ModeDir)
	afs.Mkdir("/app/baz", os.ModeDir)
	afs.Mkdir("/qux", os.ModeDir)
	afs.Mkdir("/qux/quux", os.ModeDir)
	afs.WriteFile("/app/corge", []byte("awesome file"), 0666)

	dirs, err := fs.GetDirs("/app")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if actual, expected := dirs, []string{"baz", "foo"}; !reflect.DeepEqual(actual, expected) {
		t.Errorf("GetDirs() returns %v, but expected %v", actual, expected)
	}
}

func Test_ReadFile(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	afs.Mkdir("/app", os.ModeDir)
	afs.WriteFile("/app/foobar", []byte("awesome file"), 0666)

	data, err := fs.ReadFile("/app/foobar")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if actual, expected := string(data), "awesome file"; actual != expected {
		t.Errorf("ReadFile() returns %s, but expected %v", actual, expected)
	}
}

func Test_CreateDir(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	p := "/app/foo/bar"
	ok, err := fs.CreateDir(p)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if !ok {
		t.Error("Should return true")
	}

	if existing, err := afs.DirExists(p); err != nil {
		t.Errorf("Unexpected error %v", err)
	} else if !existing {
		t.Error("Should exist")
	}
}

func Test_CreateDir_WhenAlreadyExisted(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	p := "/app/foo/bar"
	afs.Mkdir(p, 0755)
	ok, err := fs.CreateDir(p)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if ok {
		t.Error("Should return false")
	}

	if existing, err := afs.DirExists(p); err != nil {
		t.Errorf("Unexpected error %v", err)
	} else if !existing {
		t.Error("Should exist")
	}
}

func Test_Walk(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	afs.Mkdir("/app", os.ModeDir)
	afs.Mkdir("/app/foo", os.ModeDir)
	afs.Mkdir("/app/bar", os.ModeDir)
	afs.WriteFile("/app/baz", []byte("baz file"), 0666)
	afs.WriteFile("/app/bar/quz", []byte("quz file"), 0666)
	afs.Mkdir("/app/foo/qux", os.ModeDir)
	afs.Mkdir("/app/foo/quux", os.ModeDir)

	expects := map[string]bool{
		"/app":          true,
		"/app/bar":      true,
		"/app/baz":      false,
		"/app/bar/quz":  false,
		"/app/foo/qux":  true,
		"/app/foo/quux": true,
	}
	actuals := map[string]bool{}

	cb := func(path string, dir bool, err error) error {
		if err != nil {
			t.Errorf("Unexpected error %v", err)
			return err
		}
		if _, ok := expects[path]; !ok {
			t.Errorf("Unexpected visted path %q", path)
		}
		actuals[path] = dir
		return nil
	}

	err := fs.Walk("/app", cb)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	for path, expected := range expects {
		if actual, ok := actuals[path]; !ok {
			t.Errorf("%s was not visited", path)
		} else if actual != expected {
			if actual {
				t.Errorf("%s type was incorrect, got: directory, want: file", path)
			} else {
				t.Errorf("%s type was incorrect, got: file, want: directory", path)
			}
		}
	}
}
