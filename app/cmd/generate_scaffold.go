package cmd

import (
	"fmt"

	"github.com/izumin5210/scaffold/app/ui"
	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/mitchellh/cli"
	"github.com/pkg/errors"
)

type generateScaffoldCommand struct {
	scaffold       scaffold.Scaffold
	createScaffold usecase.CreateScaffoldUseCase
	ui             ui.UI
}

// NewGenerateScaffoldCommandFactories creates a command factory for ...
func NewGenerateScaffoldCommandFactories(
	getScaffolds usecase.GetScaffoldsUseCase,
	createScaffold usecase.CreateScaffoldUseCase,
	ui ui.UI,
) (map[string]cli.CommandFactory, error) {
	factories := map[string]cli.CommandFactory{}
	scffs, err := getScaffolds.Perform()
	if err != nil {
		return nil, errors.Cause(err)
	}
	for _, s := range scffs {
		factories[s.Name()] = newGenerateScaffoldCommandFactory(s, createScaffold, ui)
	}
	return factories, nil
}

func newGenerateScaffoldCommandFactory(
	s scaffold.Scaffold,
	createScaffold usecase.CreateScaffoldUseCase,
	ui ui.UI,
) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &generateScaffoldCommand{
			scaffold:       s,
			createScaffold: createScaffold,
			ui:             ui,
		}, nil
	}
}

// Synopsis returns a short synopsis of the command.
// It is an implementation of mitchellh/cli.Command#Synopsis()
func (c *generateScaffoldCommand) Synopsis() string {
	return c.scaffold.Synopsis()
}

// Help returns a long-term help text of the command.
// It is an implementation of mitchellh/cli.Command#Help()
func (c *generateScaffoldCommand) Help() string {
	return c.scaffold.Help()
}

// Run runs the actual command behavior
// It is an implementation of mitchellh/cli.Command#Run()
func (c *generateScaffoldCommand) Run(args []string) int {
	if len(args) != 1 {
		c.ui.Error(fmt.Sprintf("Invalid arguments: %v", args))
		return ui.ExitCodeInvalidArgumentListLengthError
	}
	if err := c.createScaffold.Perform(c.scaffold, args[0]); err != nil {
		c.ui.Error(fmt.Sprintf("Error: %s", err.Error()))
		return ui.ExitCodeFailedToCreatetScaffoldsError
	}
	return ui.ExitCodeOK
}
