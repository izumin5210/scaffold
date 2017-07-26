package fs

import (
	"path/filepath"

	"github.com/pkg/errors"
)

// Entry is simple representation of filesystem entries
type Entry interface {
	Path() string
	IsDir() bool
	BaseName() string
	DirName() string
}

type entry struct {
	path string
	dir  bool
}

// NewFile returns a new Entry object treated as a file
func NewFile(path string) (Entry, error) {
	return NewEntry(path, false)
}

// NewDir returns a new Entry object treated as an empty directory
func NewDir(path string) (Entry, error) {
	return NewEntry(path, true)
}

// NewEntry returns a new Entry object
func NewEntry(path string, dir bool) (Entry, error) {
	if !filepath.IsAbs(path) {
		return nil, errors.New("Entry path should be absolute path")
	}
	e := &entry{path: path, dir: dir}
	if l := len(e.path); e.path[l-1:] == "/" {
		e.path = e.path[:l-1]
	}
	return e, nil
}

func (e *entry) Path() string {
	return e.path
}

func (e *entry) IsDir() bool {
	return e.dir
}

func (e *entry) BaseName() string {
	return filepath.Base(e.Path())
}

func (e *entry) DirName() string {
	return filepath.Dir(e.Path())
}
