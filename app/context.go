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
	InReader() io.Reader
	OutWriter() io.Writer
	ErrWriter() io.Writer
	Repository() scaffold.Repository
	UI() ui.UI
	GetScaffoldsUseCase() usecase.GetScaffoldsUseCase
	CreateScaffoldUseCase() usecase.CreateScaffoldUseCase
}

type context struct {
	rootPath  string
	tmplsPath string
	inReader  io.Reader
	outWriter io.Writer
	errWriter io.Writer
	fs        fs.FS
	repo      scaffold.Repository
	ui        ui.UI
}

// NewContext creates a new context object
func NewContext(
	inReader io.Reader,
	outWriter, errWriter io.Writer,
	rootPath, tmplsPath string,
	fs fs.FS,
) Context {
	return &context{
		rootPath:  rootPath,
		tmplsPath: tmplsPath,
		inReader:  inReader,
		outWriter: outWriter,
		errWriter: errWriter,
		fs:        fs,
		ui:        ui.NewUI(inReader, outWriter, errWriter),
	}
}

func (c *context) RootPath() string {
	return c.rootPath
}

func (c *context) TemplatesPath() string {
	return c.tmplsPath
}

func (c *context) InReader() io.Reader {
	return c.inReader
}

func (c *context) OutWriter() io.Writer {
	return c.outWriter
}

func (c *context) ErrWriter() io.Writer {
	return c.errWriter
}

func (c *context) ConstructService() scaffold.ConstructService {
	return scaffold.NewConstructService(c.Repository())
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
	return usecaseImpl.NewCreateScaffoldUseCase(c.ConstructService(), c.UI())
}
