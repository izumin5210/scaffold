package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/app/ui"
)

type generateCommandTestContext struct {
	ctrl *gomock.Controller
	ui   *ui.MockUI
}

func getGenerateTestContext(t *testing.T) *generateCommandTestContext {
	ctrl := gomock.NewController(t)
	return &generateCommandTestContext{
		ctrl: ctrl,
		ui:   ui.NewMockUI(ctrl),
	}
}

func getGenerateTestCommand(ctx *generateCommandTestContext) *generateCommand {
	return &generateCommand{
		ui: ctx.ui,
	}
}

func Test_NewGenerateCommandFactory(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	f := NewGenerateCommandFactory(ctx.ui)
	cmd, err := f()

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if len(cmd.Synopsis()) == 0 {
		t.Error("Synopsis() should present")
	}

	if len(cmd.Help()) == 0 {
		t.Error("Help() should present")
	}
}

func Test_Generate_Run(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	cmd := getGenerateTestCommand(ctx)
	ctx.ui.EXPECT().Error(gomock.Any())

	code := cmd.Run([]string{})

	if actual, expected := code, ui.ExitCodeInvalidArgumentListLengthError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}

func Test_Generate_Run_WithArgs(t *testing.T) {
	ctx := getGenerateTestContext(t)
	defer ctx.ctrl.Finish()

	cmd := getGenerateTestCommand(ctx)
	ctx.ui.EXPECT().Error(gomock.Any())

	code := cmd.Run([]string{"qux"})

	if actual, expected := code, ui.ExitCodeScffoldNotFoundError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}
