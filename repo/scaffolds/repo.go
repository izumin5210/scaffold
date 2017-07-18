package scaffolds

import (
	"github.com/izumin5210/scaffold/entity"
	"github.com/izumin5210/scaffold/infra/fs"
)

// Repository is a repository for scaffolds
type Repository interface {
	GetAll() ([]*entity.Scaffold, error)
}

type repo struct {
	path string
	fs   fs.FS
}

// NewRepository returns a Repository implementation
func NewRepository(path string, fs fs.FS) Repository {
	return &repo{path: path, fs: fs}
}
