package cmd

import (
	"fmt"

	"bytes"

	"github.com/izumin5210/scaffold/app/ui"
	"github.com/mitchellh/cli"
)

type generateCommand struct {
	ui ui.UI
}

// NewGenerateCommandFactory creates a command factory for ...
func NewGenerateCommandFactory(ui ui.UI) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &generateCommand{ui: ui}, nil
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
	buf.WriteString("Usage: scaffold generate [--help] <name> [<args>]")
	return buf.String()
}

// Run runs the actual command behavior
// It is an implementation of mitchellh/cli.Command#Run()
func (c *generateCommand) Run(args []string) int {
	if len(args) > 0 {
		c.ui.Error(fmt.Sprintf("Scaffold %q is not found", args[0]))
		return ui.ExitCodeScffoldNotFoundError
	}
	c.ui.Error("Require scaffold name")
	return ui.ExitCodeInvalidArgumentListLengthError
}
