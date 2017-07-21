package entity

import "path"

// Entry represents a filesystem entry
type Entry struct {
	path     string
	existing bool
}

// Name returns a base name of this entry
func (e *Entry) Name() string {
	return path.Base(e.path)
}

// Path returns a path of this entry
func (e *Entry) Path() string {
	return e.path
}

// IsDir returns true if this entry is a directory
func (e *Entry) IsDir() bool {
	return false
}

// Exists returns true if this entry exists
func (e *Entry) Exists() bool {
	return e.existing
}

// IsParentOf returns true if the entry is parent of a given entry
func (e *Entry) IsParentOf(other *Entry) bool {
	return e.Path() == path.Dir(other.Path())
}

// IsChildOf returns true if the entry is child of a given entry
func (e *Entry) IsChildOf(other *Entry) bool {
	return other.IsParentOf(e)
}
