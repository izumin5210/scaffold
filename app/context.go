package app

import (
	"io"

	"github.com/izumin5210/scaffold/app/ui"
	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/izumin5210/scaffold/domain/scaffold"
	usecaseImpl "github.com/izumin5210/scaffold/domain/usecase"
	"github.com/izumin5210/scaffold/infra/fs"
	scaffoldrepo "github.com/izumin5210/scaffold/infra/scaffold/repo"
)

// Context is container storing configurations
type Context interface {
	RootPath() string
	TemplatesPath() string
	Repository() scaffold.Repository
	UI() ui.UI
	GetScaffoldsUseCase() usecase.GetScaffoldsUseCase
	CreateScaffoldUseCase() usecase.CreateScaffoldUseCase
}

type context struct {
	rootPath  string
	tmplsPath string
	fs        fs.FS
	repo      scaffold.Repository
	ui        ui.UI
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
		ui:        ui.NewUI(inStream, outStream, errStream),
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

func (c *context) UI() ui.UI {
	return c.ui
}

func (c *context) GetScaffoldsUseCase() usecase.GetScaffoldsUseCase {
	return usecaseImpl.NewGetScaffoldsUseCase(c.Repository())
}

func (c *context) CreateScaffoldUseCase() usecase.CreateScaffoldUseCase {
	return usecaseImpl.NewCreateScaffoldUseCase(c.RootPath(), c.Repository(), c.UI())
}
