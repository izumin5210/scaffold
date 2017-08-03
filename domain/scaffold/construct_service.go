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
	tmpls, err := s.repo.GetTemplates(sc)
	if err != nil {
		return errors.Wrapf(err, "Failed to get templates under %s", sc.Path())
	}

	exConcs, err := s.repo.GetConcreteEntries(sc, tmpls, v)
	if err != nil {
		return errors.Wrapf(err, "Failed to get existing entries of %s", sc.Path())

	}

	createdDirs := map[string]struct{}{}

	for _, tmpl := range tmpls {
		conc, err := tmpl.Compile(rootPath, v)
		if err != nil {
			return errors.Wrapf(err, "Failed to compile %s", sc.Path())
		}
		conflicted := false

		if exConc, ok := exConcs[conc.Path()]; ok {
			if conc.IsDir() {
				cb(conc.Path(), true, false, ConstructSkipped)
				continue
			}
			if got, want := exConc.Content(), conc.Content(); got == want {
				cb(conc.Path(), false, false, ConstructSkipped)
				continue
			} else {
				if cCb(conc.Path(), got, want) {
					conflicted = true
				} else {
					cb(conc.Path(), false, true, ConstructSkipped)
					continue
				}
			}
		}

		created, dirCreated, err := s.repo.Create(conc)
		if err != nil {
			return errors.Wrapf(err, "Failed to create new entry %s", conc.Path())
		}
		if _, alreadyCreated := createdDirs[conc.Dir()]; conc.Dir() != rootPath && !alreadyCreated {
			status := ConstructSuccess
			if !dirCreated {
				status = ConstructSkipped
			}
			cb(conc.Dir(), true, false, status)
			createdDirs[conc.Dir()] = struct{}{}
		}
		status := ConstructSuccess
		if !created {
			status = ConstructSkipped
		}
		cb(conc.Path(), conc.IsDir(), conflicted, status)
	}
	return nil
}
