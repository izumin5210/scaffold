package cmd

import (
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/mitchellh/cli"
)

// CreateScaffold represents a command object for scaffolding templates
// It can be treated as an mitchellh/cli.Command
type createScaffold struct {
	repo     scaffold.Repository
	scaffold scaffold.Scaffold
}

// NewCreateScaffoldCommandFactory creates a command for creating scaffold
func NewCreateScaffoldCommandFactory(
	repo scaffold.Repository,
	sc scaffold.Scaffold,
) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &createScaffold{repo: repo, scaffold: sc}, nil
	}
}

// NewCreateScaffoldCommandFactories creates comands for creating scaffold
func NewCreateScaffoldCommandFactories(
	repo scaffold.Repository,
	scaffolds []scaffold.Scaffold,
) CommandFactories {
	factories := CommandFactories{}
	for _, sc := range scaffolds {
		factories[sc.Name()] = NewCreateScaffoldCommandFactory(repo, sc)
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
	// TODO: Not yet implemented.
	return 0
}
