package usecase

import (
	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/pkg/errors"
)

type {{name | camelize}} struct {
}

// New{{name | pascalize}}UseCase creates a {{name | pascalize}}UseCase implementation
func New{{name | pascalize}}UseCase() usecase.{{name | pascalize}}UseCase {
	return &{{name | camelize}}{}
}

func (u *{{name | camelize}}) Perform() error {
	return errors.New("Not yet implemented")
}
