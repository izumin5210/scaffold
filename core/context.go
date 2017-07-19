package core

import (
	"github.com/izumin5210/scaffold/cmd"
	"github.com/izumin5210/scaffold/cmd/factory"
	"github.com/izumin5210/scaffold/infra/fs"
	"github.com/izumin5210/scaffold/repo/scaffolds"
)

// Context is container storing configurations
type Context interface {
	Path() string
	Repository() scaffolds.Repository
	CommandFactoryFactory() cmd.Factory
}

type context struct {
	path string
	fs   fs.FS
	repo scaffolds.Repository
}

// NewContext creates a new context object
func NewContext(path string, fs fs.FS) Context {
	return &context{path: path, fs: fs}
}

func (c *context) Path() string {
	return c.path
}

func (c *context) Repository() scaffolds.Repository {
	if c.repo == nil {
		c.repo = scaffolds.NewRepository(c.path, c.fs)
	}
	return c.repo
}

func (c *context) CommandFactoryFactory() cmd.Factory {
	return factory.New()
}
