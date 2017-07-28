package cmd

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/izumin5210/scaffold/app/ui"
	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/izumin5210/scaffold/domain/scaffold"

	"github.com/golang/mock/gomock"
)

type generateScaffoldTestContext struct {
	ctrl           *gomock.Controller
	createScaffold *usecase.MockCreateScaffoldUseCase
	getScaffolds   *usecase.MockGetScaffoldsUseCase
	ui             *ui.MockUI
}

func getGenerateScaffoldTestContext(t *testing.T) *generateScaffoldTestContext {
	ctrl := gomock.NewController(t)
	createScaffold := usecase.NewMockCreateScaffoldUseCase(ctrl)
	getScaffolds := usecase.NewMockGetScaffoldsUseCase(ctrl)
	ui := ui.NewMockUI(ctrl)
	return &generateScaffoldTestContext{
		ctrl:           ctrl,
		createScaffold: createScaffold,
		getScaffolds:   getScaffolds,
		ui:             ui,
	}
}

func getGenerateScaffoldTestCommand(
	ctx *generateScaffoldTestContext,
	scff scaffold.Scaffold,
) *generateScaffoldCommand {
	return &generateScaffoldCommand{
		scaffold:       scff,
		createScaffold: ctx.createScaffold,
		ui:             ctx.ui,
	}
}

func Test_NewGenerateScaffoldCommandFactories(t *testing.T) {
	ctx := getGenerateScaffoldTestContext(t)
	defer ctx.ctrl.Finish()

	scff0 := scaffold.NewMockScaffold(ctx.ctrl)
	scff1 := scaffold.NewMockScaffold(ctx.ctrl)
	scff0.EXPECT().Name().Return("model").AnyTimes()
	scff1.EXPECT().Name().Return("controller").AnyTimes()
	scffs := []scaffold.Scaffold{scff0, scff1}
	ctx.getScaffolds.EXPECT().Perform().Return(scffs, nil)

	factories, err := NewGenerateScaffoldCommandFactories(
		ctx.getScaffolds,
		ctx.createScaffold,
		ctx.ui,
	)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if actual, expected := len(factories), len(scffs); actual != expected {
		t.Errorf("Returned value has %d factories, want %d", actual, expected)
	}

	for _, s := range scffs {
		if _, ok := factories[s.Name()]; !ok {
			t.Errorf("factories[%s] should present", s.Name())
		}
	}
}

func Test_GenerateScaffold_Run(t *testing.T) {
	ctx := getGenerateScaffoldTestContext(t)
	defer ctx.ctrl.Finish()

	scff := scaffold.NewMockScaffold(ctx.ctrl)
	cmd := getGenerateScaffoldTestCommand(ctx, scff)
	name := "foo"

	ctx.createScaffold.EXPECT().Perform(scff, name).Return(nil)

	code := cmd.Run([]string{name})

	if actual, expected := code, ui.ExitCodeOK; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}

func Test_GenerateScaffold_Run_WithoutName(t *testing.T) {
	ctx := getGenerateScaffoldTestContext(t)
	defer ctx.ctrl.Finish()

	scff := scaffold.NewMockScaffold(ctx.ctrl)
	cmd := getGenerateScaffoldTestCommand(ctx, scff)
	ctx.ui.EXPECT().Error(gomock.Any())

	code := cmd.Run([]string{})

	if actual, expected := code, ui.ExitCodeInvalidArgumentListLengthError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}

func Test_GenerateScaffold_Run_WhenCreateScffoldFaild(t *testing.T) {
	ctx := getGenerateScaffoldTestContext(t)
	defer ctx.ctrl.Finish()

	scff := scaffold.NewMockScaffold(ctx.ctrl)
	cmd := getGenerateScaffoldTestCommand(ctx, scff)
	name := "foo"

	ctx.createScaffold.EXPECT().Perform(scff, name).Return(errors.New("error"))
	ctx.ui.EXPECT().Error(gomock.Any())

	code := cmd.Run([]string{name})

	if actual, expected := code, ui.ExitCodeFailedToCreatetScaffoldsError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}
