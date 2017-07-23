//go:generate mockgen -source=repo.go -package scaffold -destination=repo_mock.go

package scaffold

type ConstructStatus int

const (
	ConstructSuccess ConstructStatus = iota + 1
	ConstructSkipped
	ConstructFailed
)

func (s ConstructStatus) IsSuccess() bool {
	return s == ConstructSuccess
}

func (s ConstructStatus) IsSkipped() bool {
	return s == ConstructSkipped
}

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
