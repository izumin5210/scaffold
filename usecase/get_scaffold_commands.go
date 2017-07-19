package usecase

import (
	"github.com/izumin5210/scaffold/cmd"
	"github.com/izumin5210/scaffold/repo/scaffolds"
)

// GetScaffoldCommandUseCase is an use-case for loading scaffolds
type GetScaffoldCommandUseCase interface {
	Perform() (cmd.CommandFactories, error)
}

type getScaffoldCommandUseCase struct {
	repo    scaffolds.Repository
	factory cmd.Factory
}

// NewGetScaffoldCommandUseCase creates a GetUseCase instance
func NewGetScaffoldCommandUseCase(
	repo scaffolds.Repository,
	factory cmd.Factory,
) GetScaffoldCommandUseCase {
	return &getScaffoldCommandUseCase{repo: repo, factory: factory}
}

func (u *getScaffoldCommandUseCase) Perform() (cmd.CommandFactories, error) {
	scaffolds, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}

	factories := cmd.CommandFactories{}
	for _, sc := range scaffolds {
		factories[sc.Name()] = u.factory.CreateCreateScaffoldCommandFactory(sc)
	}

	return factories, nil
}
