package concrete

import (
	"path/filepath"
)

// Entry is an existing file or directory that created from the scaffold template
type Entry interface {
	Path() string
	Dir() string
	Content() string
	IsDir() bool
}

type entry struct {
	path    string
	content string
	dir     bool
}

// NewFile returns a new TemplateEntry object treated as a file
func NewFile(path, content string) Entry {
	return NewEntry(path, content, false)
}

// NewDir returns a new TemplateEntry object treated as a directory
func NewDir(path string) Entry {
	return NewEntry(path, "", true)
}

// NewEntry returns a new Entry object
func NewEntry(path, content string, dir bool) Entry {
	return &entry{
		path:    path,
		content: content,
		dir:     dir,
	}
}

func (e *entry) Path() string {
	return e.path
}

func (e *entry) Content() string {
	return e.content
}

func (e *entry) IsDir() bool {
	return e.dir
}

func (e *entry) Dir() string {
	return filepath.Dir(e.Path())
}
