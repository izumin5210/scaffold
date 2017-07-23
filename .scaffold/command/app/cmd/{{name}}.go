package cmd

import (
	"github.com/mitchellh/cli"
)

type {{name | camelize}}Command struct {
}

// New{{name | pascalize}}CommandFactory creates a command factory for ...
func New{{name | pascalize}}CommandFactory() cli.CommandFactory {
	return func() (cli.Command, error) {
		return &{{name | camelize}}Command{}, nil
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
	return 0
}
