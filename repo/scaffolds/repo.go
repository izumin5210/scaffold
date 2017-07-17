package repo

import (
	"github.com/izumin5210/scaffold/entity"
)

// ScaffoldsRepository is a repository for scaffolds
type ScaffoldsRepository interface {
	GetAll() ([]*entity.Scaffold, error)
}

type repo struct {
	context entity.Context
}

// NewScaffoldsRepository returns a ScaffoldsRepository implementation
func NewScaffoldsRepository(context entity.Context) ScaffoldsRepository {
	return &repo{context: context}
}
