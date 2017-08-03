package main

import (
	"fmt"

	"github.com/izumin5210/scaffold/app/ui"

	"github.com/izumin5210/scaffold/app"
	"github.com/izumin5210/scaffold/app/cmd"
	mcli "github.com/mitchellh/cli"
)

// CLI parses arguments and runs commands
type CLI interface {
	Run(args []string) int
}

type cli struct {
	ctx app.Context
	cli *mcli.CLI
}

// NewCLI returns a new CLI instance
func NewCLI(ctx app.Context, name, version, revision string) CLI {
	mcli := mcli.NewCLI(name, fmt.Sprintf("%s (%s)", version, revision))
	mcli.HelpWriter = ctx.ErrWriter()
	return &cli{
		ctx: ctx,
		cli: mcli,
	}
}

func (c *cli) Run(args []string) int {
	c.cli.Args = args
	cmds, err := c.getCommands()
	c.cli.Commands = cmds
	if !c.cli.IsVersion() && !c.cli.IsHelp() && len(args) != 0 && err != nil {
		c.ctx.UI().Error(err.Error())
		return ui.ExitCodeFailedToGetScaffoldsError
	}
	exitStatus, err := c.cli.Run()
	if err != nil {
		c.ctx.UI().Error(err.Error())
	}
	return exitStatus
}

func (c *cli) getCommands() (cmd.CommandFactories, error) {
	factories := cmd.CommandFactories{}
	genScffFactories, err := cmd.NewGenerateCommandFactories(
		c.ctx.GetScaffoldsUseCase(),
		c.ctx.CreateScaffoldUseCase(),
		c.ctx.UI(),
		c.ctx.TemplatesPath(),
	)
	for n, f := range genScffFactories {
		factories[fmt.Sprintf("%s %s", ui.CommandGenerate, n)] = f
		factories[fmt.Sprintf("%s %s", ui.CommandGenerateShort, n)] = f
	}
	factories[ui.CommandGenerate] = cmd.NewEmptyCommandFactory(
		ui.CommandGenerate,
		ui.CommandGenerateSynopsis,
		c.ctx.UI(),
	)
	factories[ui.CommandGenerateShort] = factories[ui.CommandGenerate]
	return factories, err
}
