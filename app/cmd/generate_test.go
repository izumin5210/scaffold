package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/app/ui"
	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/pkg/errors"
)

type generateCommandTestContext struct {
	ctrl           *gomock.Controller
	getScaffolds   *usecase.MockGetScaffoldsUseCase
	createScaffold *usecase.MockCreateScaffoldUseCase
	ui             *ui.MockUI
}

func getGenerateTestContext(t *testing.T) *generateCommandTestContext {
	ctrl := gomock.NewController(t)
	return &generateCommandTestContext{
		ctrl:           ctrl,
		getScaffolds:   usecase.NewMockGetScaffoldsUseCase(ctrl),
		createScaffold: usecase.NewMockCreateScaffoldUseCase(ctrl),
		ui:             ui.NewMockUI(ctrl),
	}
}

func getGenerateTestCommand(ctx *generateCommandTestContext) *generateCommand {
	return &generateCommand{
		createScaffold: ctx.createScaffold,
		ui:             ctx.ui,
	}
}

func Test_NewGenerateCommandFactory(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	scffs := []scaffold.Scaffold{
		scaffold.NewScaffold("/app/.scaffold/foo", &scaffold.Meta{}),
		scaffold.NewScaffold("/app/.scaffold/bar", &scaffold.Meta{}),
		scaffold.NewScaffold("/app/.scaffold/baz", &scaffold.Meta{}),
	}

	ctx.getScaffolds.EXPECT().
		Perform().
		Return(scffs, nil)

	f := NewGenerateCommandFactory(ctx.getScaffolds, ctx.createScaffold, ctx.ui)
	cmd, err := f()

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if len(cmd.Synopsis()) == 0 {
		t.Error("GenerateCommand#Synopsis() should be present")
	}

	if len(cmd.Help()) == 0 {
		t.Error("GenerateCommand#Help() should be present")
	}
}

func Test_NewGenerateCommandFactory_WhenFailToGetScaffolds(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	ctx.getScaffolds.EXPECT().
		Perform().
		Return(nil, errors.New("Failed to get scaffolds"))

	ctx.ui.EXPECT().Error(gomock.Any()).MinTimes(1)

	f := NewGenerateCommandFactory(ctx.getScaffolds, ctx.createScaffold, ctx.ui)
	cmd, err := f()

	if err == nil {
		t.Error("NewGenerateCommandFactory() Should return error")
	}

	if cmd != nil {
		t.Error("NewGenerateCommandFactory() Should not return commands")
	}
}

func Test_Run(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	cmd := getGenerateTestCommand(ctx)
	cmd.scaffolds = []scaffold.Scaffold{
		scaffold.NewScaffold("/app/.scaffold/foo", &scaffold.Meta{}),
		scaffold.NewScaffold("/app/.scaffold/bar", &scaffold.Meta{}),
		scaffold.NewScaffold("/app/.scaffold/baz", &scaffold.Meta{}),
	}
	scffName, name := "bar", "qux"

	ctx.createScaffold.EXPECT().
		Perform(cmd.scaffolds[1], name).
		Return(nil)

	code := cmd.Run([]string{scffName, name})

	if actual, expected := code, ui.ExitCodeOK; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}

func Test_Run_WhenArgListToShort(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()
	cmd := getGenerateTestCommand(ctx)

	ctx.ui.EXPECT().Error(gomock.Any()).MinTimes(1)

	code := cmd.Run([]string{"foo"})

	if actual, expected := code, ui.ExitCodeInvalidArgumentListLengthError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}

func Test_Run_WhenArgumentListTooLong(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()
	cmd := getGenerateTestCommand(ctx)

	ctx.ui.EXPECT().Error(gomock.Any()).MinTimes(1)

	code := cmd.Run([]string{"foo", "bar", "baz"})

	if actual, expected := code, ui.ExitCodeInvalidArgumentListLengthError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}

func Test_Run_WhenScaffoldsDoesNotExist(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	cmd := getGenerateTestCommand(ctx)
	cmd.scaffolds = []scaffold.Scaffold{
		scaffold.NewScaffold("/app/.scaffold/foo", &scaffold.Meta{}),
	}
	scffName, name := "bar", "qux"

	ctx.ui.EXPECT().Error(gomock.Any()).MinTimes(1)

	code := cmd.Run([]string{scffName, name})

	if actual, expected := code, ui.ExitCodeScffoldNotFoundError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}

func Test_Run_WhenFailToCreateScaffold(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	cmd := getGenerateTestCommand(ctx)
	cmd.scaffolds = []scaffold.Scaffold{
		scaffold.NewScaffold("/app/.scaffold/foo", &scaffold.Meta{}),
		scaffold.NewScaffold("/app/.scaffold/bar", &scaffold.Meta{}),
		scaffold.NewScaffold("/app/.scaffold/baz", &scaffold.Meta{}),
	}
	scffName, name := "bar", "qux"

	ctx.createScaffold.EXPECT().
		Perform(cmd.scaffolds[1], name).
		Return(errors.New("Failed to create scaffolds"))

	ctx.ui.EXPECT().Error(gomock.Any()).MinTimes(1)

	code := cmd.Run([]string{scffName, name})

	if actual, expected := code, ui.ExitCodeFailedToCreatetScaffoldsError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}
