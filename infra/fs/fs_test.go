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

func Test_GetDirs_WhenPathDoesNotExist(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	dirs, err := fs.GetDirs("/app/foo/bar")

	if dirs != nil {
		t.Errorf("Should not return any directories")
	}

	if err == nil {
		t.Errorf("Should return an error")
	}
}

func Test_GetDirs_WhenPathPointsToFile(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	afs.WriteFile("/app/foo", []byte("awesome file"), 0666)
	dirs, err := fs.GetDirs("/app/foo")

	if dirs != nil {
		t.Errorf("Should not return any directories")
	}

	if err == nil {
		t.Errorf("Should return an error")
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

func Test_CreateFile(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	err := fs.CreateFile("/app/foo", "foobarbaz")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if body, err := afs.ReadFile("/app/foo"); err != nil {
		t.Errorf("Unexpected error %v", err)
	} else if actual, expected := string(body), "foobarbaz"; actual != expected {
		t.Errorf("Created file content is %q, want %q", actual, expected)
	}
}

func Test_Remove(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	afs.Mkdir("/foo", 0755)
	afs.WriteFile("/foo/bar", []byte("bar"), 0644)
	afs.WriteFile("/foo/baz", []byte("baz"), 0644)

	if err := fs.Remove("/foo/baz"); err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if ok, _ := afs.Exists("/foo/baz"); ok {
		t.Error("/foo/baz should be deleted")
	}
}

func Test_Exists(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	afs.Mkdir("/foo", 0755)
	afs.WriteFile("/bar", []byte("baz"), 0644)

	tests := []struct {
		in  string
		out bool
	}{
		{in: "/foo", out: true},
		{in: "/bar", out: true},
		{in: "/baz", out: false},
		{in: "foo", out: false},
		{in: "foo/bar/baz", out: false},
	}

	for _, ts := range tests {
		if ok, err := fs.Exists(ts.in); err != nil {
			t.Errorf("Unexpected error %v", err)
		} else if actual, expected := ok, ts.out; actual != expected {
			t.Errorf("Returns %t, want %t", actual, expected)
		}
	}
}

func Test_DirExists(t *testing.T) {
	afs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs := &fs{afs: afs}

	afs.Mkdir("/foo", 0755)
	afs.WriteFile("/bar", []byte("baz"), 0644)

	tests := []struct {
		in  string
		out bool
	}{
		{in: "/foo", out: true},
		{in: "/bar", out: false},
		{in: "/baz", out: false},
		{in: "foo", out: false},
		{in: "foo/bar/baz", out: false},
	}

	for _, ts := range tests {
		if ok, err := fs.DirExists(ts.in); err != nil {
			t.Errorf("Unexpected error %v", err)
		} else if actual, expected := ok, ts.out; actual != expected {
			t.Errorf("Returns %t, want %t", actual, expected)
		}
	}
}
