//go:generate mockgen -source=repo.go -package scaffold -destination=repo_mock.go

package scaffold

// Repository is a repository for scaffolds
type Repository interface {
	GetAll() ([]Scaffold, error)
}
