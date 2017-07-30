package scaffold

// Entry represents a filesystem file or directory
type Entry interface {
	Path() string
	IsDir() bool
	Content() string
}

type entry struct {
	path    string
	content string
	dir     bool
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
