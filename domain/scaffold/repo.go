//go:generate mockgen -source=repo.go -package scaffold -destination=repo_mock.go

package scaffold

// ConstructStatus represents a status of constructing processes
type ConstructStatus int

const (
	// ConstructSuccess is returned when constructing is succeeded
	ConstructSuccess ConstructStatus = iota + 1
	// ConstructSkipped is returned when constructing is skipped
	ConstructSkipped
	// ConstructFailed is returned when constructing is failed
	ConstructFailed
)

// IsSuccess returns true if a status represents success
func (s ConstructStatus) IsSuccess() bool {
	return s == ConstructSuccess
}

// IsSkipped returns true if a status represents skipped
func (s ConstructStatus) IsSkipped() bool {
	return s == ConstructSkipped
}

// IsFailed returns true if a status represents failed
func (s ConstructStatus) IsFailed() bool {
	return s == ConstructFailed
}

func (s ConstructStatus) String() string {
	switch s {
	case ConstructSuccess:
		return "ConstructSuccess"
	case ConstructSkipped:
		return "ConstructSkipped"
	case ConstructFailed:
		return "ConstructFailed"
	default:
		return "UnknownStatus"
	}
}

// ConstructCallback is called after files and directories is created
type ConstructCallback func(path string, dir bool, status ConstructStatus)

// Repository is a repository for scaffolds
type Repository interface {
	GetAll() ([]Scaffold, error)
	Construct(
		scff Scaffold,
		name string,
		cb ConstructCallback,
	) error
}
