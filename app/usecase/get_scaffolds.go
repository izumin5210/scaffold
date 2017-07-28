package usecase

import (
	"github.com/izumin5210/scaffold/domain/scaffold"
)

// GetScaffoldsUseCase is an use-case for loading scaffolds
type GetScaffoldsUseCase interface {
	Perform() ([]scaffold.Scaffold, error)
}
