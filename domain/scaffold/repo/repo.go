package repo

import (
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/izumin5210/scaffold/infra/fs"
)

type repo struct {
	path string
	fs   fs.FS
}

// New returns a Repository implementation
func New(path string, fs fs.FS) scaffold.Repository {
	return &repo{path: path, fs: fs}
}
