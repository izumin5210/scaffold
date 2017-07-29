package cmd

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/app/ui"
)

type emptyCommandTestContext struct {
	ctrl     *gomock.Controller
	name     string
	synopsis string
	ui       *ui.MockUI
}

func getEmptyTestContext(t *testing.T) *emptyCommandTestContext {
	ctrl := gomock.NewController(t)
	name := "awesome"
	synopsis := "this is awesome command"
	return &emptyCommandTestContext{
		ctrl:     ctrl,
		name:     name,
		synopsis: synopsis,
		ui:       ui.NewMockUI(ctrl),
	}
}

func getEmptyTestCommand(ctx *emptyCommandTestContext) *emptyCommand {
	return &emptyCommand{
		name:     ctx.name,
		synopsis: ctx.synopsis,
		ui:       ctx.ui,
	}
}

func Test_NewEmptyCommandFactory(t *testing.T) {
	ctx := getEmptyTestContext(t)
	defer ctx.ctrl.Finish()

	f := NewEmptyCommandFactory(ctx.name, ctx.synopsis, ctx.ui)
	cmd, err := f()

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if actual, expected := cmd.Synopsis(), ctx.synopsis; actual != expected {
		t.Errorf("Synopsis() returns %q, want %q", actual, expected)
	}

	if !strings.Contains(cmd.Help(), ctx.name) {
		t.Errorf("Help() returns %q, want to contain %q", cmd.Help(), ctx.name)
	}
}

func Test_Empty_Run(t *testing.T) {
	ctx := getEmptyTestContext(t)
	defer ctx.ctrl.Finish()

	cmd := getEmptyTestCommand(ctx)
	ctx.ui.EXPECT().Error(gomock.Any())

	code := cmd.Run([]string{})

	if actual, expected := code, ui.ExitCodeInvalidArgumentListLengthError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}

func Test_Empty_Run_WithArgs(t *testing.T) {
	ctx := getEmptyTestContext(t)
	defer ctx.ctrl.Finish()

	cmd := getEmptyTestCommand(ctx)
	ctx.ui.EXPECT().Error(gomock.Any())

	code := cmd.Run([]string{"qux"})

	if actual, expected := code, ui.ExitCodeScffoldNotFoundError; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}
