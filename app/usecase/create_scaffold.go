package usecase

import (
	"github.com/izumin5210/scaffold/domain/scaffold"
)

// CreateScaffoldUseCase is an use-case for loading scaffolds
type CreateScaffoldUseCase interface {
	Perform(scff scaffold.Scaffold, name string) error
}

type createScaffoldUseCase struct {
	repo scaffold.Repository
}

// NewCreateScaffoldUseCase creates a CreateScaffoldUseCase instance
func NewCreateScaffoldUseCase(repo scaffold.Repository) CreateScaffoldUseCase {
	return &createScaffoldUseCase{repo: repo}
}

func (u *createScaffoldUseCase) Perform(scff scaffold.Scaffold, name string) error {
	return u.repo.Construct(scff, name)
}
