package scaffold

// Directory represents a directory entry of filesystem
type Directory struct {
	Entry
	entries []Entry
}

// NewDirectory returns a directory object
func NewDirectory(path string, entries []Entry) *Directory {
	return &Directory{Entry: &entry{path: path}, entries: entries}
}

// EmptyDirectory returns an empty directory object
func EmptyDirectory(path string) *Directory {
	return &Directory{Entry: &entry{path: path}}
}

// IsDir returns true if this entry is a directory
func (d *Directory) IsDir() bool {
	return true
}

// Entries returns child entries of this directory
func (d *Directory) Entries() []Entry {
	return d.entries
}
