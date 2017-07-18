package scaffolds

import (
	"github.com/izumin5210/scaffold/entity"
)

// Repository is a repository for scaffolds
type Repository interface {
	GetAll() ([]*entity.Scaffold, error)
}

type repo struct {
	context *entity.Context
}

// NewRepository returns a Repository implementation
func NewRepository(context *entity.Context) Repository {
	return &repo{context: context}
}
