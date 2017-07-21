package scaffold

// Scaffold represents an executable subcommand
type Scaffold struct {
	*Directory
	*Meta
}

// NewScaffold reeturns a new scaffold object
func NewScaffold(path string, meta *Meta) *Scaffold {
	return &Scaffold{
		Directory: EmptyDirectory(path),
		Meta:      meta,
	}
}

// Name returns this scaffold name
func (s *Scaffold) Name() string {
	if len(s.Meta.Name) > 0 {
		return s.Meta.Name
	}
	return s.Directory.Name()
}
