//go:generate mockgen -source=factory.go -package factory -destination=factory/mock.go

package cmd

import (
	"github.com/izumin5210/scaffold/entity"
	"github.com/mitchellh/cli"
)

// Factory provides factory functions for generating cli.CommandFactory
type Factory interface {
	CreateCreateScaffoldCommandFactory(sc *entity.Scaffold) cli.CommandFactory
}
