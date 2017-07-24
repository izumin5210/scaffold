package usecase

// {{name | pascalize}}UseCase is ...
type {{name | pascalize}}UseCase interface {
	Perform() error
}

type {{name | camelize}}UseCase struct {
}

// New{{name | pascalize}}UseCase creates a {{name | pascalize}} instance
func New{{name | pascalize}}UseCase() {{name | pascalize}} {
	return &{{name | camelize}}{}
}

func (u *{{name | camelize}}) Perform() error {
	return nil
}

