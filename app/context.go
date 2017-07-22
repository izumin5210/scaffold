package app

import (
	"github.com/izumin5210/scaffold/app/cmd"
	"github.com/izumin5210/scaffold/app/cmd/factory"
	"github.com/izumin5210/scaffold/domain/scaffold"
	scaffoldrepo "github.com/izumin5210/scaffold/domain/scaffold/repo"
	"github.com/izumin5210/scaffold/infra/fs"
)

// Context is container storing configurations
type Context interface {
	Path() string
	Repository() scaffold.Repository
	CommandFactoryFactory() cmd.Factory
}

type context struct {
	path string
	fs   fs.FS
	repo scaffold.Repository
}

// NewContext creates a new context object
func NewContext(path string, fs fs.FS) Context {
	return &context{path: path, fs: fs}
}

func (c *context) Path() string {
	return c.path
}

func (c *context) Repository() scaffold.Repository {
	if c.repo == nil {
		c.repo = scaffoldrepo.New(c.path, c.fs)
	}
	return c.repo
}

func (c *context) CommandFactoryFactory() cmd.Factory {
	return factory.New()
}
