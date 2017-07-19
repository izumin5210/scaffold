package cmd

import (
	"github.com/mitchellh/cli"
)

// CommandFactories is a type alias for map[string]mitchellh/cli.CommandFactory
type CommandFactories map[string]cli.CommandFactory
