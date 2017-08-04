package usecase

import (
	"github.com/izumin5210/scaffold/domain/scaffold"
)

// CreateScaffoldUseCase is an use-case for loading scaffolds
type CreateScaffoldUseCase interface {
	Perform(scff scaffold.Scaffold, rootPath, name string) error
}
