package usecase

import (
	"github.com/izumin5210/scaffold/domain/scaffold"
)

// GetScaffoldsUseCase is an use-case for loading scaffolds
type GetScaffoldsUseCase interface {
	Perform() ([]scaffold.Scaffold, error)
}

type getScaffoldsUseCase struct {
	repo scaffold.Repository
}

// NewGetScaffoldsUseCase creates a GetUseCase instance
func NewGetScaffoldsUseCase(
	repo scaffold.Repository,
) GetScaffoldsUseCase {
	return &getScaffoldsUseCase{repo: repo}
}

func (u *getScaffoldsUseCase) Perform() ([]scaffold.Scaffold, error) {
	return u.repo.GetAll()
}
