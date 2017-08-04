package cmd

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/izumin5210/scaffold/app/ui"
	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/izumin5210/scaffold/domain/scaffold"

	"github.com/golang/mock/gomock"
)

type generateTestContext struct {
	ctrl           *gomock.Controller
	createScaffold *usecase.MockCreateScaffoldUseCase
	getScaffolds   *usecase.MockGetScaffoldsUseCase
	ui             *ui.MockUI
	rootPath       string
	tmplsPath      string
}

func getGenerateTestContext(t *testing.T) *generateTestContext {
	ctrl := gomock.NewController(t)
	createScaffold := usecase.NewMockCreateScaffoldUseCase(ctrl)
	getScaffolds := usecase.NewMockGetScaffoldsUseCase(ctrl)
	ui := ui.NewMockUI(ctrl)
	return &generateTestContext{
		ctrl:           ctrl,
		createScaffold: createScaffold,
		getScaffolds:   getScaffolds,
		ui:             ui,
		rootPath:       "/app",
		tmplsPath:      "/app/.scaffold",
	}
}

func getGenerateTestCommand(
	ctx *generateTestContext,
	scff scaffold.Scaffold,
) *generateCommand {
	return &generateCommand{
		rootPath:       ctx.rootPath,
		scaffold:       scff,
		createScaffold: ctx.createScaffold,
		ui:             ctx.ui,
	}
}

func Test_NewGenerateCommandFactories(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	scffs := []scaffold.Scaffold{}
	for _, n := range []string{"model", "controller", "view"} {
		s := scaffold.NewMockScaffold(ctx.ctrl)
		s.EXPECT().Name().Return(n).MinTimes(1)
		s.EXPECT().Synopsis().Return(fmt.Sprintf("%s synopsis", n)).Times(2)
		s.EXPECT().Help().Return(fmt.Sprintf("%s help", n)).Times(2)
		scffs = append(scffs, s)
	}

	ctx.getScaffolds.EXPECT().Perform(ctx.tmplsPath).Return(scffs, nil)

	factories, err := NewGenerateCommandFactories(
		ctx.getScaffolds,
		ctx.createScaffold,
		ctx.ui,
		ctx.rootPath,
		ctx.tmplsPath,
	)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if actual, expected := len(factories), len(scffs); actual != expected {
		t.Errorf("Returned value has %d factories, want %d", actual, expected)
	}

	for _, s := range scffs {
		if f, ok := factories[s.Name()]; !ok {
			t.Errorf("factories[%s] should present", s.Name())
		} else if cmd, err := f(); err != nil {
			t.Errorf("Unexpected error %v", err)
		} else {
			if actual, expected := cmd.Synopsis(), s.Synopsis(); actual != expected {
				t.Errorf("Synopsis() returns %q, want %q", actual, expected)
			}
			if actual, expected := cmd.Help(), s.Help(); actual != expected {
				t.Errorf("Help() returns %q, want %q", actual, expected)
			}
		}
	}
}

func Test_Generate_Run(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	scff := scaffold.NewMockScaffold(ctx.ctrl)
	cmd := getGenerateTestCommand(ctx, scff)
	name := "foo"

	ctx.createScaffold.EXPECT().Perform(scff, ctx.rootPath, name).Return(nil)

	code := cmd.Run([]string{name})

	if actual, expected := code, ui.ExitCodeOK; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}

func Test_Generate_Run_WithoutName(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	scff := scaffold.NewMockScaffold(ctx.ctrl)
	cmd := getGenerateTestCommand(ctx, scff)
	ctx.ui.EXPECT().Error(gomock.Any())

	code := cmd.Run([]string{})

	if actual, expected := code, ui.ExitCodeInvalidArgumentListLengthError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}

func Test_Generate_Run_WhenCreateScffoldFaild(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	scff := scaffold.NewMockScaffold(ctx.ctrl)
	cmd := getGenerateTestCommand(ctx, scff)
	name := "foo"

	ctx.createScaffold.EXPECT().Perform(scff, ctx.rootPath, name).Return(errors.New("error"))
	ctx.ui.EXPECT().Error(gomock.Any())

	code := cmd.Run([]string{name})

	if actual, expected := code, ui.ExitCodeFailedToCreatetScaffoldsError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}
