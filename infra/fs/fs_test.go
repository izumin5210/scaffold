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
	afs.Mkdir("/app/bar/foo", os.ModeDir)
	afs.Mkdir("/app/baz", os.ModeDir)
	afs.Mkdir("/qux", os.ModeDir)
	afs.Mkdir("/qux/quux", os.ModeDir)
	afs.WriteFile("/app/corge", []byte("awesome file"), 0666)

	dirs, err := fs.GetDirs("/app")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if actual, expected := dirs, []string{"/app/foo", "/app/baz"}; reflect.DeepEqual(actual, expected) {
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
