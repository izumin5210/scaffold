package factory

import (
	"github.com/izumin5210/scaffold/cmd"
	"github.com/izumin5210/scaffold/entity"
	"github.com/mitchellh/cli"
)

type factory struct {
}

// New reutnrs a new object implemented cmd.Factory
func New() cmd.Factory {
	return &factory{}
}

func (f *factory) CreateCreateScaffoldCommandFactory(sc *entity.Scaffold) cli.CommandFactory {
	return func() (cli.Command, error) {
		return cmd.NewCreateScaffoldCommand(sc), nil
	}
}