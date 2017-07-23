package usecase

import (
	"path/filepath"

	"github.com/izumin5210/scaffold/app"
	"github.com/izumin5210/scaffold/domain/scaffold"
)

// CreateScaffoldUseCase is an use-case for loading scaffolds
type CreateScaffoldUseCase interface {
	Perform(scff scaffold.Scaffold, name string) error
}

type createScaffoldUseCase struct {
	rootPath string
	repo     scaffold.Repository
	ui       app.UI
}

// NewCreateScaffoldUseCase creates a CreateScaffoldUseCase instance
func NewCreateScaffoldUseCase(rootPath string, repo scaffold.Repository, ui app.UI) CreateScaffoldUseCase {
	return &createScaffoldUseCase{rootPath: rootPath, repo: repo, ui: ui}
}

func (u *createScaffoldUseCase) Perform(scff scaffold.Scaffold, name string) error {
	return u.repo.Construct(scff, name, u.constructCallback)
}

func (u *createScaffoldUseCase) constructCallback(path string, dir bool, status scaffold.ConstructStatus) {
	relpath, _ := filepath.Rel(u.rootPath, path)
	if status.IsSuccess() {
		u.ui.Status("create", relpath, app.UIColorGreen)
	} else if status.IsSkipped() {
		if dir {
			u.ui.Status("exist", relpath, app.UIColorBlue)
		} else {
			u.ui.Status("identical", relpath, app.UIColorBlue)
		}
	}
}
