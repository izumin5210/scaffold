package entity

// File represents a file entry of filesystem
type File struct {
	*Entry
}

// NewFile returns a file object
func NewFile(path string, existing bool) *File {
	return &File{Entry: &Entry{path: path, existing: existing}}
}
