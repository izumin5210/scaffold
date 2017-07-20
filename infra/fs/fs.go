package fs

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// FS is filesystem wrapper interface
type FS interface {
	GetDirs(path string) ([]string, error)
	ReadFile(path string) ([]byte, error)
	Walk(path string, cb func(path string, dir bool, err error) error) error
}

type fs struct {
	afs afero.Afero
}

// New returns FS instance using the os package
func New() FS {
	return &fs{afs: afero.Afero{Fs: afero.NewOsFs()}}
}

func (f *fs) GetDirs(name string) ([]string, error) {
	entries, err := f.afs.ReadDir(name)
	if err != nil {
		return nil, errors.Cause(err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

func (f *fs) Walk(root string, cb func(path string, dir bool, err error) error) error {
	return f.afs.Walk(root, func(path string, info os.FileInfo, err error) error {
		return cb(path, info.IsDir(), err)
	})
}

func (f *fs) ReadFile(name string) ([]byte, error) {
	return f.afs.ReadFile(name)
}
