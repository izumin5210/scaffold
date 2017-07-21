//go:generate mockgen -source=scaffold.go -package scaffold -destination=scaffold_mock.go

package scaffold

// Scaffold represents an executable subcommand
type Scaffold interface {
	Path() string
	Name() string
	Synopsis() string
	Help() string
}

type scaffold struct {
	dir  *Directory
	meta *Meta
}

// NewScaffold reeturns a new scaffold object
func NewScaffold(path string, meta *Meta) Scaffold {
	return &scaffold{
		dir:  EmptyDirectory(path),
		meta: meta,
	}
}

func (s *scaffold) Path() string {
	return s.dir.Path()
}

func (s *scaffold) Name() string {
	if len(s.meta.Name) > 0 {
		return s.meta.Name
	}
	return s.dir.Name()
}

func (s *scaffold) Synopsis() string {
	return s.meta.Synopsis
}

func (s *scaffold) Help() string {
	return s.meta.Help
}
