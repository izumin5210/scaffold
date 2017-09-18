package template

import (
	"github.com/pkg/errors"

	"github.com/izumin5210/scaffold/domain/model/scaffold/concrete"
)

// Entry represents a scaffold template entry
type Entry interface {
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

func (e *entry) Compile(v interface{}) (concrete.Entry, error) {
	path, err := String(e.path).Compile(v)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not compile path: %q", e.path)
	}
	content, err := String(e.content).Compile(v)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not compile content: %q", e.path)
	}
	return concrete.NewEntry(path, content, e.dir), nil
}
