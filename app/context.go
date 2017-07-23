package app

import (
	"io"

	"github.com/izumin5210/scaffold/domain/scaffold"
	scaffoldrepo "github.com/izumin5210/scaffold/domain/scaffold/repo"
	"github.com/izumin5210/scaffold/infra/fs"
)

// Context is container storing configurations
type Context interface {
	RootPath() string
	TemplatesPath() string
	Repository() scaffold.Repository
	UI() UI
}

type context struct {
	rootPath  string
	tmplsPath string
	fs        fs.FS
	repo      scaffold.Repository
	ui        UI
}

// NewContext creates a new context object
func NewContext(
	inStream io.Reader,
	outStream, errStream io.Writer,
	rootPath, tmplsPath string,
	fs fs.FS,
) Context {
	return &context{
		rootPath:  rootPath,
		tmplsPath: tmplsPath,
		fs:        fs,
		ui:        NewUI(inStream, outStream, errStream),
	}
}

func (c *context) RootPath() string {
	return c.rootPath
}

func (c *context) TemplatesPath() string {
	return c.tmplsPath
}

func (c *context) Repository() scaffold.Repository {
	if c.repo == nil {
		c.repo = scaffoldrepo.New(c.rootPath, c.tmplsPath, c.fs)
	}
	return c.repo
}

func (c *context) UI() UI {
	return c.ui
}
