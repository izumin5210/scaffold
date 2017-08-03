package scaffold

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
	// TODO: Should handle errors
	tmpls, _ := s.repo.GetTemplates(sc)
	// TODO: Should handle errors
	exConcs, _ := s.repo.GetConcreteEntries(sc, tmpls, v)

	createdDirs := map[string]struct{}{}

	for _, tmpl := range tmpls {
		// TODO: Should handle errors
		conc, _ := tmpl.Compile(rootPath, v)
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

		// TODO: Should handle errors
		created, dirCreated, _ := s.repo.Create(conc)
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
