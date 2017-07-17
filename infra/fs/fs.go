package infra

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// FS is filesystem wrapper interface
type FS interface {
	GetDirs(path string) ([]string, error)
	ReadFile(path string) ([]byte, error)
}

type fs struct {
	afs afero.Afero
}

// NewFS returns FS instance using the os package
func NewFS() FS {
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

func (f *fs) ReadFile(name string) ([]byte, error) {
	return f.afs.ReadFile(name)
}
