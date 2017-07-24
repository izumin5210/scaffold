package usecase

import (
	"testing"

	"path/filepath"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/app/ui"
	"github.com/izumin5210/scaffold/domain/scaffold"
)

type createScaffoldTestContext struct {
	ctrl     *gomock.Controller
	rootPath string
	repo     *scaffold.MockRepository
	ui       *ui.MockUI
}

func getCreateScaffoldTestContext(t *testing.T) *createScaffoldTestContext {
	ctrl := gomock.NewController(t)
	return &createScaffoldTestContext{
		ctrl:     ctrl,
		rootPath: "/app/.scaffold",
		repo:     scaffold.NewMockRepository(ctrl),
		ui:       ui.NewMockUI(ctrl),
	}
}

func getCreateScaffoldTestUseCase(ctx *createScaffoldTestContext) CreateScaffoldUseCase {
	return &createScaffoldUseCase{
		rootPath: ctx.rootPath,
		repo:     ctx.repo,
		ui:       ctx.ui,
	}
}

func Test_CreateScaffold_Perform(t *testing.T) {
	ctx := getCreateScaffoldTestContext(t)
	defer ctx.ctrl.Finish()

	type call struct {
		dir        bool
		conflicted bool
		status     scaffold.ConstructStatus
		prefix     string
		color      ui.ColorAttrs
		oldContent string
		newContent string
		overwrite  bool
	}

	calls := map[string]*call{
		"bar": {
			dir:        true,
			conflicted: false,
			status:     scaffold.ConstructSuccess,
			prefix:     "create",
			color:      ui.ColorGreen,
		},
		"bar/baz": {
			dir:        false,
			conflicted: false,
			status:     scaffold.ConstructSuccess,
			prefix:     "create",
			color:      ui.ColorGreen,
		},
		"bar/qux": {
			dir:        false,
			conflicted: true,
			status:     scaffold.ConstructSuccess,
			prefix:     "force",
			color:      ui.ColorYellow,
			oldContent: "barqux",
			newContent: "bar\nqux",
			overwrite:  false,
		},
		"quux": {
			dir:        true,
			conflicted: false,
			status:     scaffold.ConstructSkipped,
			prefix:     "exist",
			color:      ui.ColorBlue,
		},
		"corge": {
			dir:        false,
			conflicted: false,
			status:     scaffold.ConstructSkipped,
			prefix:     "identical",
			color:      ui.ColorBlue,
		},
		"bar/grault": {
			dir:        false,
			conflicted: true,
			status:     scaffold.ConstructSkipped,
			prefix:     "skip",
			color:      ui.ColorYellow,
			oldContent: "bargrault",
			newContent: "bar\ngrault",
			overwrite:  true,
		},
	}

	name := "bar"
	scff := scaffold.NewScaffold(filepath.Join(ctx.rootPath, "foo"), &scaffold.Meta{})
	ctx.repo.EXPECT().Construct(scff, name, gomock.Any(), gomock.Any()).
		Do(func(
			_ scaffold.Scaffold,
			_ string,
			cb scaffold.ConstructCallback,
			conflictedCb scaffold.ConstructConflictedCallback,
		) error {
			for p, c := range calls {
				abspath := filepath.Join(ctx.rootPath, p)
				if c.conflicted {
					ctx.ui.EXPECT().Status("conflicted", p, ui.ColorRed).Times(1)
					conflictedCb(abspath, c.oldContent, c.newContent)
				}
				ctx.ui.EXPECT().Status(c.prefix, p, c.color)
				cb(abspath, c.dir, c.conflicted, c.status)
			}
			return nil
		}).
		Times(1)

	ctx.ui.EXPECT().Ask(gomock.Any()).Return("Y", nil).Times(1).
		After(ctx.ui.EXPECT().Ask(gomock.Any()).Return("n", nil).Times(1))

	u := getCreateScaffoldTestUseCase(ctx)
	err := u.Perform(scff, name)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
