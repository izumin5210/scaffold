//go:generate mockgen -source=scaffold.go -package scaffold -destination=scaffold_mock.go

package scaffold

import "path/filepath"

// Scaffold represents an executable subcommand
type Scaffold interface {
	Path() string
	Name() string
	Synopsis() string
	Help() string
}

type scaffold struct {
	path string
	meta *Meta
}

// NewScaffold reeturns a new scaffold object
func NewScaffold(path string, meta *Meta) Scaffold {
	return &scaffold{
		path: path,
		meta: meta,
	}
}

func (s *scaffold) Path() string {
	return s.path
}

func (s *scaffold) Name() string {
	if len(s.meta.Name) > 0 {
		return s.meta.Name
	}
	return filepath.Base(s.path)
}

func (s *scaffold) Synopsis() string {
	return s.meta.Synopsis
}

func (s *scaffold) Help() string {
	return s.meta.Help
}
