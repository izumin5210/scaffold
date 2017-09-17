package template

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/izumin5210/scaffold/domain/model/scaffold/concrete"
)

// Entry represents a scaffold template entry
type Entry interface {
	Path() String
	Dir() String
	Content() String
	IsDir() bool
	Compile(v interface{}) (concrete.Entry, error)
}

type entry struct {
	path    String
	content String
	dir     bool
}

// NewFile returns a new Entry object treated as a file
func NewFile(path, content String) Entry {
	return NewEntry(path, content, false)
}

// NewDir returns a new Entry object treated as a directory
func NewDir(path String) Entry {
	return NewEntry(path, "", true)
}

// NewEntry returns a new Entry object
func NewEntry(path, content String, dir bool) Entry {
	return &entry{
		path:    path,
		content: content,
		dir:     dir,
	}
}

func (e *entry) Path() String {
	return e.path
}

func (e *entry) IsDir() bool {
	return e.dir
}

func (e *entry) Dir() String {
	return String(filepath.Dir(string(e.Path())))
}

func (e *entry) Content() String {
	return e.content
}

func (e *entry) Compile(v interface{}) (concrete.Entry, error) {
	path, err := e.path.Compile(string(e.path), v)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not compile path: %q", e.path)
	}
	content, err := e.content.Compile(string(e.content), v)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not compile content: %q", e.path)
	}
	return concrete.NewEntry(path, content, e.dir), nil
}
