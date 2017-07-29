package cmd

import (
	"fmt"

	"bytes"

	"github.com/izumin5210/scaffold/app/ui"
	"github.com/mitchellh/cli"
)

type emptyCommand struct {
	name     string
	synopsis string
	ui       ui.UI
}

// NewEmptyCommandFactory creates a command factory for ...
func NewEmptyCommandFactory(name, synopsis string, ui ui.UI) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &emptyCommand{
			name:     name,
			synopsis: synopsis,
			ui:       ui,
		}, nil
	}
}

// Synopsis returns a short synopsis of the command.
// It is an implementation of mitchellh/cli.Command#Synopsis()
func (c *emptyCommand) Synopsis() string {
	return c.synopsis
}

// Help returns a long-term help text of the command.
// It is an implementation of mitchellh/cli.Command#Help()
func (c *emptyCommand) Help() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Usage: scaffold %s [--help] <name> [<args>]", c.name))
	return buf.String()
}

// Run runs the actual command behavior
// It is an implementation of mitchellh/cli.Command#Run()
func (c *emptyCommand) Run(args []string) int {
	if len(args) > 0 {
		c.ui.Error(fmt.Sprintf("Scaffold %q is not found", args[0]))
		return ui.ExitCodeScffoldNotFoundError
	}
	c.ui.Error("Require scaffold name")
	return ui.ExitCodeInvalidArgumentListLengthError
}
