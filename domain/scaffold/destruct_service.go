//go:generate mockgen -source=destruct_service.go -package scaffold -destination=destruct_service_mock.go

package scaffold

import (
	"github.com/pkg/errors"
)

// DestructService removes directories and files from scaffold templates
type DestructService interface {
	Perform(
		rootPath string,
		sc Scaffold,
		v interface{},
	) ([]string, error)
}

type destructService struct {
	repo Repository
}

// NewDestructService creates DestructService implementation instance
func NewDestructService(repo Repository) DestructService {
	return &destructService{
		repo: repo,
	}
}

func (s *destructService) Perform(
	rootPath string,
	sc Scaffold,
	v interface{},
) ([]string, error) {
	return []string{}, errors.New("Not yet implemented")
}
