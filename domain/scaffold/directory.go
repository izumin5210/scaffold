package scaffold

// Directory represents a directory entry of filesystem
type Directory struct {
	Entry
	children map[string]Entry
}

// NewDirectory returns a directory object
func NewDirectory(path string, entries []Entry) *Directory {
	entry := &entry{path: path}
	children := map[string]Entry{}
	descendants := []Entry{}
	for _, e := range entries {
		if entry.IsParentOf(e) {
			children[e.Name()] = e
		} else {
			descendants = append(descendants, e)
		}
	}
	for n, e := range children {
		if e.IsDir() {
			children[n] = NewDirectory(e.Path(), descendants)
		} else {
			children[n] = e
		}
	}
	return &Directory{Entry: entry, children: children}
}

// EmptyDirectory returns an empty directory object
func EmptyDirectory(path string) *Directory {
	return &Directory{Entry: &entry{path: path}}
}

// IsDir returns true if this entry is a directory
func (d *Directory) IsDir() bool {
	return true
}

// Children returns child entries of this directory
func (d *Directory) Children() map[string]Entry {
	return d.children
}
