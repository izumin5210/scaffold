package cmd

import (
	"fmt"

	"bytes"

	"github.com/izumin5210/scaffold/app/ui"
	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/mitchellh/cli"
	"github.com/pkg/errors"
)

type generateCommand struct {
	scaffolds      []scaffold.Scaffold
	createScaffold usecase.CreateScaffoldUseCase
	ui             ui.UI
}

// NewGenerateCommandFactory creates a command factory for ...
func NewGenerateCommandFactory(
	getScaffolds usecase.GetScaffoldsUseCase,
	createScaffold usecase.CreateScaffoldUseCase,
	ui ui.UI,
) cli.CommandFactory {
	return func() (cli.Command, error) {
		scffs, err := getScaffolds.Perform()
		if err != nil {
			ui.Error("Cloud not load scaffolds")
			return nil, errors.Wrap(err, "")
		}
		return &generateCommand{
			scaffolds:      scffs,
			createScaffold: createScaffold,
			ui:             ui,
		}, nil
	}
}

// Synopsis returns a short synopsis of the command.
// It is an implementation of mitchellh/cli.Command#Synopsis()
func (c *generateCommand) Synopsis() string {
	return "Generate new code"
}

// Help returns a long-term help text of the command.
// It is an implementation of mitchellh/cli.Command#Help()
func (c *generateCommand) Help() string {
	var buf bytes.Buffer
	buf.WriteString("Usage: scaffold generate [--help] <name> [<args>]\n\n")
	buf.WriteString("Available scaffolds are:\n")
	for _, scff := range c.scaffolds {
		buf.WriteString(fmt.Sprintf("    %s    %s\n", scff.Name(), scff.Synopsis()))
	}
	return buf.String()
}

// Run runs the actual command behavior
// It is an implementation of mitchellh/cli.Command#Run()
func (c *generateCommand) Run(args []string) int {
	if len(args) != 2 {
		c.ui.Error(fmt.Sprintf("Invalid arguments: %v", args))
		return ui.ExitCodeInvalidArgumentListLengthError
	}
	scffName, name := args[0], args[1]
	scff := c.getScaffoldByName(scffName)
	if scff == nil {
		c.ui.Error(fmt.Sprintf("Could not found scaffold %q", scffName))
		return ui.ExitCodeScffoldNotFoundError
	}
	if err := c.createScaffold.Perform(scff, name); err != nil {
		c.ui.Error(fmt.Sprintf("Error: %s", err.Error()))
		return ui.ExitCodeFailedToCreatetScaffoldsError
	}
	return ui.ExitCodeOK
}

func (c *generateCommand) getScaffoldByName(name string) scaffold.Scaffold {
	for _, s := range c.scaffolds {
		if s.Name() == name {
			return s
		}
	}
	return nil
}
