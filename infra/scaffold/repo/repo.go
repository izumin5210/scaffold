package repo

import (
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/izumin5210/scaffold/infra/fs"
)

type repo struct {
	rootPath  string
	tmplsPath string
	fs        fs.FS
}

// New returns a Repository implementation
func New(rootPath, tmplsPath string, fs fs.FS) scaffold.Repository {
	return &repo{rootPath: rootPath, tmplsPath: tmplsPath, fs: fs}
}
