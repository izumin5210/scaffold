package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/izumin5210/scaffold/app/ui"
	"github.com/izumin5210/scaffold/domain/scaffold"

	"github.com/izumin5210/scaffold/app/usecase"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/app"
)

type cliTestContext struct {
	ctrl           *gomock.Controller
	ctx            *app.MockContext
	in             *bytes.Buffer
	out            *bytes.Buffer
	err            *bytes.Buffer
	ui             *ui.MockUI
	getScaffolds   *usecase.MockGetScaffoldsUseCase
	createScaffold *usecase.MockCreateScaffoldUseCase
	name           string
	dir            string
	version        string
	revision       string
}

func getCLITestContext(t *testing.T) *cliTestContext {
	ctrl := gomock.NewController(t)
	ctx := app.NewMockContext(ctrl)
	inBuf := bytes.NewBuffer([]byte{})
	outBuf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	ui := ui.NewMockUI(ctrl)
	getScaffolds := usecase.NewMockGetScaffoldsUseCase(ctrl)
	createScaffold := usecase.NewMockCreateScaffoldUseCase(ctrl)
	ctx.EXPECT().InReader().Return(inBuf).AnyTimes()
	ctx.EXPECT().OutWriter().Return(outBuf).AnyTimes()
	ctx.EXPECT().ErrWriter().Return(errBuf).AnyTimes()
	ctx.EXPECT().GetScaffoldsUseCase().Return(getScaffolds).AnyTimes()
	ctx.EXPECT().CreateScaffoldUseCase().Return(createScaffold).AnyTimes()
	ctx.EXPECT().UI().Return(ui).AnyTimes()
	name := "scaffold"
	dir := "/app/.scaffold"
	version := "0.1.1"
	revision := "aaaaaaa"
	ctx.EXPECT().TemplatesPath().Return(dir).AnyTimes()
	return &cliTestContext{
		ctrl:           ctrl,
		ctx:            ctx,
		in:             inBuf,
		out:            outBuf,
		err:            errBuf,
		ui:             ui,
		getScaffolds:   getScaffolds,
		createScaffold: createScaffold,
		name:           name,
		dir:            "/app/.scaffold",
		version:        version,
		revision:       revision,
	}
}

func getTestCLI(ctx *cliTestContext) CLI {
	return NewCLI(ctx.ctx, ctx.name, ctx.version, ctx.revision)
}

func Test_CLI_Run_WithVersion(t *testing.T) {
	ctx := getCLITestContext(t)
	defer ctx.ctrl.Finish()

	for _, args := range [][]string{{"-v"}, {"--version"}} {
		cli := getTestCLI(ctx)
		ctx.getScaffolds.EXPECT().Perform(ctx.dir).Return([]scaffold.Scaffold{}, nil).Times(1)

		if actual, expected := cli.Run(args), ui.ExitCodeOK; actual != expected {
			t.Errorf("Run(%v) returns %d, want %d", args, actual, expected)
		}

		if actual, expected := ctx.err.String(), ctx.version; !strings.Contains(actual, expected) {
			t.Errorf("Run(%v) outputs %q to error stream, want to contain %q", args, actual, expected)
		}

		if actual, expected := ctx.err.String(), ctx.revision; !strings.Contains(actual, expected) {
			t.Errorf("Run(%v) outputs %q to error stream, want to contain %q", args, actual, expected)
		}

		if actual := ctx.out.String(); len(actual) != 0 {
			t.Errorf("Unexpected outputs to stdout %q", actual)
		}
	}
}

func Test_CLI_Run_WhenGetScaffoldsFailed(t *testing.T) {
	ctx := getCLITestContext(t)
	defer ctx.ctrl.Finish()
	cli := getTestCLI(ctx)

	ctx.getScaffolds.EXPECT().Perform(ctx.dir).Return(nil, errors.New("error"))
	ctx.ui.EXPECT().Error(gomock.Any())

	if actual, expected := cli.Run([]string{"g"}), ui.ExitCodeFailedToGetScaffoldsError; actual != expected {
		t.Errorf("Run() returns %d, want %d", actual, expected)
	}

	if actual := ctx.out.String(); len(actual) != 0 {
		t.Errorf("Unexpected outputs to stdout %q", actual)
	}
}

func Test_CLI_Run_WhenGetScaffoldsFailed_WithVersionOrHelp(t *testing.T) {
	ctx := getCLITestContext(t)
	defer ctx.ctrl.Finish()

	for _, args := range [][]string{{"-v"}, {"--version"}, {"-h"}, {"--help"}} {
		ctx.getScaffolds.EXPECT().Perform(ctx.dir).Return(nil, errors.New("error"))
		cli := getTestCLI(ctx)

		if actual, expected := cli.Run(args), ui.ExitCodeOK; actual != expected {
			t.Errorf("Run(%v) returns %d, want %d", args, actual, expected)
		}
	}
}

func Test_CLI_Run_WhenGetScaffoldsFailed_WithNoArgs(t *testing.T) {
	ctx := getCLITestContext(t)
	defer ctx.ctrl.Finish()

	ctx.getScaffolds.EXPECT().Perform(ctx.dir).Return(nil, errors.New("error"))
	cli := getTestCLI(ctx)
	cli.Run([]string{})
}
