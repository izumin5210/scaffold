package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// TemplateEntry represents a scaffold template entry
type TemplateEntry interface {
	Entry
	Compile(root string, v interface{}) (Entry, error)
}

type templateEntry struct {
	path         TemplateString
	content      TemplateString
	templateRoot string
	dir          bool
}

// NewTemplateFile returns a new TemplateEntry object treated as a file
func NewTemplateFile(path, content TemplateString, templateRoot string) TemplateEntry {
	return NewTemplateEntry(path, content, false, templateRoot)
}

// NewTemplateDir returns a new TemplateEntry object treated as a directory
func NewTemplateDir(path TemplateString, templateRoot string) TemplateEntry {
	return NewTemplateEntry(path, "", true, templateRoot)
}

// NewTemplateEntry returns a new TemplateEntry object
func NewTemplateEntry(path, content TemplateString, dir bool, templateRoot string) TemplateEntry {
	return &templateEntry{
		path:         path,
		content:      content,
		dir:          dir,
		templateRoot: templateRoot,
	}
}

func (e *templateEntry) Path() string {
	return string(e.path)
}

func (e *templateEntry) IsDir() bool {
	return e.dir
}

func (e *templateEntry) Dir() string {
	return filepath.Dir(e.Path())
}

func (e *templateEntry) Content() string {
	return string(e.content)
}

func (e *templateEntry) Compile(root string, v interface{}) (Entry, error) {
	path, err := e.path.Compile(string(e.path), v)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not compile path: %q", e.path)
	}
	path = strings.Replace(path, e.templateRoot, root, 1)
	content, err := e.content.Compile(string(e.path), v)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not compile content: %q", e.path)
	}
	return NewEntry(path, content, e.dir), nil
}
