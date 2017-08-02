package scaffold

import (
	"github.com/pkg/errors"
)

// ConstructService creates directories and files from scaffold templates
type ConstructService interface {
	Perform(
		rootPath string,
		sc Scaffold,
		v interface{},
		cb ConstructCallback,
		cCb ConstructConflictedCallback,
	) error
}

type constructService struct {
	repo Repository
}

// NewConstructService creates ConstructService implementation instance
func NewConstructService(repo Repository) ConstructService {
	return &constructService{
		repo: repo,
	}
}

func (s *constructService) Perform(
	rootPath string,
	sc Scaffold,
	v interface{},
	cb ConstructCallback,
	cCb ConstructConflictedCallback,
) error {
	return errors.New("errors")
}
