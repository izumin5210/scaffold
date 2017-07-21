package scaffold

import "path"

// Entry represents a filesystem entry
type Entry interface {
	Name() string
	Path() string
	IsDir() bool
	IsParentOf(other Entry) bool
	IsChildOf(other Entry) bool
}

type entry struct {
	path string
}

// Name returns a base name of this entry
func (e *entry) Name() string {
	return path.Base(e.path)
}

// Path returns a path of this entry
func (e *entry) Path() string {
	return e.path
}

// IsDir returns true if this entry is a directory
func (e *entry) IsDir() bool {
	return false
}

// IsParentOf returns true if the entry is parent of a given entry
func (e *entry) IsParentOf(other Entry) bool {
	return e.Path() == path.Dir(other.Path())
}

// IsChildOf returns true if the entry is child of a given entry
func (e *entry) IsChildOf(other Entry) bool {
	return other.IsParentOf(e)
}
