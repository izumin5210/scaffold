package usecase

import (
	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/izumin5210/scaffold/domain/scaffold"
)

type getScaffoldsUseCase struct {
	repo scaffold.Repository
}

// NewGetScaffoldsUseCase creates a GetUseCase instance
func NewGetScaffoldsUseCase(
	repo scaffold.Repository,
) usecase.GetScaffoldsUseCase {
	return &getScaffoldsUseCase{repo: repo}
}

func (u *getScaffoldsUseCase) Perform(dir string) ([]scaffold.Scaffold, error) {
	return u.repo.GetScaffolds(dir)
}
