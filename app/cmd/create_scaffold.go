package cmd

import (
	"github.com/izumin5210/scaffold/app"
	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/mitchellh/cli"
)

// CreateScaffold represents a command object for scaffolding templates
// It can be treated as an mitchellh/cli.Command
type createScaffold struct {
	rootPath string
	repo     scaffold.Repository
	ui       app.UI
	scaffold scaffold.Scaffold
}

// NewCreateScaffoldCommandFactory creates a command for creating scaffold
func NewCreateScaffoldCommandFactory(
	rootPath string,
	repo scaffold.Repository,
	ui app.UI,
	sc scaffold.Scaffold,
) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &createScaffold{rootPath: rootPath, repo: repo, ui: ui, scaffold: sc}, nil
	}
}

// NewCreateScaffoldCommandFactories creates comands for creating scaffold
func NewCreateScaffoldCommandFactories(
	rootPath string,
	repo scaffold.Repository,
	ui app.UI,
	scaffolds []scaffold.Scaffold,
) CommandFactories {
	factories := CommandFactories{}
	for _, sc := range scaffolds {
		factories[sc.Name()] = NewCreateScaffoldCommandFactory(rootPath, repo, ui, sc)
	}
	return factories
}

// Synopsis returns a short synopsis of the command.
// It is an implementation of mitchellh/cli.Command#Synopsis()
func (sc *createScaffold) Synopsis() string {
	return sc.scaffold.Synopsis()
}

// Help returns a long-term help text of the command.
// It is an implementation of mitchellh/cli.Command#Help()
func (sc *createScaffold) Help() string {
	return sc.scaffold.Help()
}

// Run runs the actual command behavior
// It is an implementation of mitchellh/cli.Command#Run()
func (sc *createScaffold) Run(args []string) int {
	u := usecase.NewCreateScaffoldUseCase(sc.rootPath, sc.repo, sc.ui)
	err := u.Perform(sc.scaffold, args[0])
	if err != nil {
		// TODO: Should constantize exit codes
		return 1
	}
	return 0
}
