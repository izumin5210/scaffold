package entity

// Scaffold represents an executable subcommand
type Scaffold struct {
	*Directory
	*ScaffoldMeta
}

func NewScaffold(path string, meta *ScaffoldMeta) *Scaffold {
	return &Scaffold{
		Directory:    EmptyDirectory(path),
		ScaffoldMeta: meta,
	}
}

func (s *Scaffold) Name() string {
	if len(s.ScaffoldMeta.Name) > 0 {
		return s.ScaffoldMeta.Name
	}
	return s.Directory.Name()
}
