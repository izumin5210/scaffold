package scaffold

// ConcreteEntry is an existing file or directory that created from the scaffold template
type ConcreteEntry interface {
	Entry
	TemplatePath() string
}

type concreteEntry struct {
	Entry
	templatePath string
}

// NewConcreteFile returns a new TemplateEntry object treated as a file
func NewConcreteFile(path, content, templatePath string) ConcreteEntry {
	return NewConcreteEntry(path, content, false, templatePath)
}

// NewConcreteDir returns a new TemplateEntry object treated as a directory
func NewConcreteDir(path, templatePath string) ConcreteEntry {
	return NewConcreteEntry(path, "", true, templatePath)
}

// NewConcreteEntry returns a new Entry object
func NewConcreteEntry(path, content string, dir bool, templatePath string) ConcreteEntry {
	return &concreteEntry{
		Entry:        NewEntry(path, content, dir),
		templatePath: templatePath,
	}
}

func (e *concreteEntry) TemplatePath() string {
	return e.templatePath
}
