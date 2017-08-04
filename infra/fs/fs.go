package fs

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// FS is filesystem wrapper interface
type FS interface {
	GetEntries(path string, compact bool) ([]Entry, error)
	GetDirs(path string) ([]string, error)
	ReadFile(path string) ([]byte, error)
	CreateDir(path string) (bool, error)
	CreateFile(path string, content string) error
	Exists(path string) (bool, error)
	DirExists(path string) (bool, error)
}

type fs struct {
	afs afero.Afero
}

// New returns FS instance using the os package
func New() FS {
	return &fs{afs: afero.Afero{Fs: afero.NewOsFs()}}
}

func (f *fs) GetDirs(name string) ([]string, error) {
	if ok, err := f.afs.IsDir(name); err != nil {
		return nil, errors.Cause(err)
	} else if !ok {
		return nil, errors.Errorf("GetDirs(string) requires a directory path, %q is a file", name)
	}
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

func (f *fs) CreateDir(path string) (bool, error) {
	if existing, err := f.afs.DirExists(path); err != nil {
		return false, errors.Wrapf(err, "Failed to check existence of %q", path)
	} else if !existing {
		err := f.afs.MkdirAll(path, 0755)
		if err != nil {
			return false, errors.Wrapf(err, "Failed to create directory at %q", path)
		}
		return true, nil
	}
	return false, nil
}

func (f *fs) CreateFile(path string, content string) error {
	return f.afs.WriteFile(path, []byte(content), 0644)
}

func (f *fs) Exists(path string) (bool, error) {
	return f.afs.Exists(path)
}

func (f *fs) DirExists(path string) (bool, error) {
	return f.afs.DirExists(path)
}
