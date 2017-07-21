package scaffold

// File represents a file entry of filesystem
type File struct {
	*Entry
}

// NewFile returns a file object
func NewFile(path string) *File {
	return &File{Entry: &Entry{path: path}}
}
