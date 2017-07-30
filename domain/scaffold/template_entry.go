package scaffold

import (
	"github.com/pkg/errors"
)

// TemplateEntry represents a scaffold template entry
type TemplateEntry interface {
	Entry
	Compile(v interface{}) (Entry, error)
}

type templateEntry struct {
	path    TemplateString
	content TemplateString
	dir     bool
}

// NewTemplateFile returns a new TemplateEntry object treated as a file
func NewTemplateFile(path, content TemplateString) TemplateEntry {
	return NewTemplateEntry(path, content, false)
}

// NewTemplateDir returns a new TemplateEntry object treated as a directory
func NewTemplateDir(path TemplateString) TemplateEntry {
	return NewTemplateEntry(path, "", true)
}

// NewTemplateEntry returns a new TemplateEntry object
func NewTemplateEntry(path, content TemplateString, dir bool) TemplateEntry {
	return &templateEntry{
		path:    path,
		content: content,
		dir:     dir,
	}
}

func (e *templateEntry) Path() string {
	return string(e.path)
}

func (e *templateEntry) IsDir() bool {
	return e.dir
}

func (e *templateEntry) Content() string {
	return string(e.content)
}

func (e *templateEntry) Compile(v interface{}) (Entry, error) {
	path, err := e.path.Compile(string(e.path), v)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not compile path: %q", e.path)
	}
	content, err := e.content.Compile(string(e.path), v)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not compile content: %q", e.path)
	}
	return NewEntry(path, content, e.dir), nil
}
