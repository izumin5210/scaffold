package usecase

import (
	"fmt"
	"path/filepath"

	"github.com/izumin5210/scaffold/app/ui"
	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/izumin5210/scaffold/domain/scaffold"
)

type createScaffoldUseCase struct {
	constructSvc scaffold.ConstructService
	ui           ui.UI
}

type createScaffoldParams struct {
	Name string
}

// NewCreateScaffoldUseCase creates a CreateScaffoldUseCase instance
func NewCreateScaffoldUseCase(
	constructSvc scaffold.ConstructService,
	ui ui.UI,
) usecase.CreateScaffoldUseCase {
	return &createScaffoldUseCase{constructSvc: constructSvc, ui: ui}
}

func (u *createScaffoldUseCase) Perform(scff scaffold.Scaffold, rootPath, name string) error {
	return u.constructSvc.Perform(
		rootPath,
		scff,
		&createScaffoldParams{Name: name},
		u.getConstructCallback(rootPath),
		u.getConstructConflictedCallback(rootPath),
	)
}

func (u *createScaffoldUseCase) getConstructCallback(rootPath string) scaffold.ConstructCallback {
	return func(path string, dir, conflicted bool, status scaffold.ConstructStatus) {
		relpath, _ := filepath.Rel(rootPath, path)
		if status.IsSuccess() {
			if conflicted {
				u.ui.Status("force", relpath, ui.ColorYellow)
			} else {
				u.ui.Status("create", relpath, ui.ColorGreen)
			}
		} else if status.IsSkipped() {
			if dir {
				u.ui.Status("exist", relpath, ui.ColorBlue)
			} else if conflicted {
				u.ui.Status("skip", relpath, ui.ColorYellow)
			} else {
				u.ui.Status("identical", relpath, ui.ColorBlue)
			}
		}
	}

}

func (u *createScaffoldUseCase) getConstructConflictedCallback(
	rootPath string,
) scaffold.ConstructConflictedCallback {
	return func(path, oldContent, newContent string) bool {
		relpath, _ := filepath.Rel(rootPath, path)
		u.ui.Status("conflicted", relpath, ui.ColorRed)
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
}
