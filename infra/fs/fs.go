package fs

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// FS is filesystem wrapper interface
type FS interface {
	GetEntries(path string, recursive bool) ([]Entry, error)
	GetDirs(path string) ([]string, error)
	ReadFile(path string) ([]byte, error)
	Walk(path string, cb func(path string, dir bool, err error) error) error
	CreateDir(path string) error
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
		if info.IsDir() {
			entries, err := f.afs.ReadDir(path)
			if err != nil {
				return err
			}
			if len(entries) > 0 {
				onlyDir := true
				for _, entry := range entries {
					onlyDir = onlyDir && entry.IsDir()
				}
				if onlyDir {
					return nil
				}
			}
		}
		return cb(path, info.IsDir(), err)
	})
}

func (f *fs) ReadFile(name string) ([]byte, error) {
	return f.afs.ReadFile(name)
}

func (f *fs) CreateDir(path string) error {
	return f.afs.MkdirAll(path, 0755)
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
