package usecase

import (
	"fmt"
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
	return u.repo.Construct(scff, name, u.constructCallback, u.constructConflictCallback)
}

func (u *createScaffoldUseCase) constructCallback(path string, dir, conflicted bool, status scaffold.ConstructStatus) {
	relpath, _ := filepath.Rel(u.rootPath, path)
	if status.IsSuccess() {
		if conflicted {
			u.ui.Status("force", relpath, app.UIColorYellow)
		} else {
			u.ui.Status("create", relpath, app.UIColorGreen)
		}
	} else if status.IsSkipped() {
		if dir {
			u.ui.Status("exist", relpath, app.UIColorBlue)
		} else if conflicted {
			u.ui.Status("skip", relpath, app.UIColorYellow)
		} else {
			u.ui.Status("identical", relpath, app.UIColorBlue)
		}
	}
}

func (u *createScaffoldUseCase) constructConflictCallback(path, oldContent, newContent string) bool {
	relpath, _ := filepath.Rel(u.rootPath, path)
	u.ui.Status("conflicted", relpath, app.UIColorRed)
	q := fmt.Sprintf("Overwrite %s? [Yn]", path)
	// https://github.com/erikhuda/thor/blob/69cff50300d63b287eb89df2933ffa218f4b2e6e/lib/thor/shell/basic.rb#L339-L348
	// q := fmt.Sprintf("Overwrite %s? [Ynaqdh]", path)
	for {
		ans, _ := u.ui.Ask(q)
		switch ans {
		case "Y":
			return true
		case "n":
			return false
		}
	}
}
