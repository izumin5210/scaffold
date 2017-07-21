package entity

// Scaffold represents an executable subcommand
type Scaffold struct {
	*Directory
	*ScaffoldMeta
}

// NewScaffold reeturns a new scaffold object
func NewScaffold(path string, meta *ScaffoldMeta) *Scaffold {
	return &Scaffold{
		Directory:    EmptyDirectory(path),
		ScaffoldMeta: meta,
	}
}

// Name returns this scaffold name
func (s *Scaffold) Name() string {
	if len(s.ScaffoldMeta.Name) > 0 {
		return s.ScaffoldMeta.Name
	}
	return s.Directory.Name()
}
