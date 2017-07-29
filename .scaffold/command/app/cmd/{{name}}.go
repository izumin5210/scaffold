package cmd

import (
	"github.com/izumin5210/scaffold/app/ui"
	"github.com/mitchellh/cli"
	"github.com/pkg/errors"
)

type {{name | camelize}}Command struct {
	ui ui.UI
}

// New{{name | pascalize}}CommandFactory creates a command factory for ...
func New{{name | pascalize}}CommandFactory(
	ui ui.UI,
) cli.CommandFactory {
	return func() (cli.Command, error) {
		// TODO: Not yet implemented.
		return &{{name | camelize}}Command{
			ui: ui,
		}, errors.New("Not yet implemented")
	}
}

// Synopsis returns a short synopsis of the command.
// It is an implementation of mitchellh/cli.Command#Synopsis()
func (c *{{name | camelize}}Command) Synopsis() string {
	// TODO: Not yet implemented.
	return ""
}

// Help returns a long-term help text of the command.
// It is an implementation of mitchellh/cli.Command#Help()
func (c *{{name | camelize}}Command) Help() string {
	// TODO: Not yet implemented.
	return ""
}

// Run runs the actual command behavior
// It is an implementation of mitchellh/cli.Command#Run()
func (c *{{name | camelize}}Command) Run(args []string) int {
	// TODO: Not yet implemented.
	return ui.ExitCodeUnknownError
}
