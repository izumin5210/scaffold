package main

import (
	"bytes"
	"strings"
	"testing"

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
	getScaffolds   *usecase.MockGetScaffoldsUseCase
	createScaffold *usecase.MockCreateScaffoldUseCase
	name           string
	version        string
	revision       string
	cli            CLI
}

func getCLITestContext(t *testing.T) *cliTestContext {
	ctrl := gomock.NewController(t)
	ctx := app.NewMockContext(ctrl)
	inBuf := bytes.NewBuffer([]byte{})
	outBuf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	getScaffolds := usecase.NewMockGetScaffoldsUseCase(ctrl)
	createScaffold := usecase.NewMockCreateScaffoldUseCase(ctrl)
	ctx.EXPECT().InReader().Return(inBuf).AnyTimes()
	ctx.EXPECT().OutWriter().Return(outBuf).AnyTimes()
	ctx.EXPECT().ErrWriter().Return(errBuf).AnyTimes()
	ctx.EXPECT().GetScaffoldsUseCase().Return(getScaffolds).AnyTimes()
	ctx.EXPECT().CreateScaffoldUseCase().Return(createScaffold).AnyTimes()
	ctx.EXPECT().UI().Return(ui.NewMockUI(ctrl)).AnyTimes()
	name := "scaffold"
	version := "0.1.1"
	revision := "aaaaaaa"
	return &cliTestContext{
		ctrl:           ctrl,
		ctx:            ctx,
		in:             inBuf,
		out:            outBuf,
		err:            errBuf,
		getScaffolds:   getScaffolds,
		createScaffold: createScaffold,
		name:           name,
		version:        version,
		revision:       revision,
		cli:            NewCLI(ctx, name, version, revision),
	}
}

func Test_CLI_Run_WithVersion(t *testing.T) {
	ctx := getCLITestContext(t)
	defer ctx.ctrl.Finish()

	ctx.getScaffolds.EXPECT().Perform().Return([]scaffold.Scaffold{}, nil).AnyTimes()

	for _, args := range [][]string{{"-v"}, {"--version"}} {
		if actual, expected := ctx.cli.Run(args), 0; actual != expected {
			t.Errorf("Run() returns %d, want %d", actual, expected)
		}

		if actual, expected := ctx.err.String(), ctx.version; !strings.Contains(actual, expected) {
			t.Errorf("Run() outputs %q to error stream, want to contain %q", actual, expected)
		}

		if actual, expected := ctx.err.String(), ctx.revision; !strings.Contains(actual, expected) {
			t.Errorf("Run() outputs %q to error stream, want to contain %q", actual, expected)
		}

		if actual := ctx.out.String(); len(actual) != 0 {
			t.Errorf("Unexpected outputs to stdout %q", actual)
		}
	}
}
