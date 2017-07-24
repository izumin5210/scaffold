package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/app/ui"
)

type {{name | camelize}}TestContext struct {
	ctrl *gomock.Controller
}

func get{{name | pascalize}}TestContext(t *testing.T) *{{name | camelize}}TestContext {
	ctrl := gomock.NewController(t)
	return &{{name | camelize}}TestContext{
		ctrl: ctrl,
	}
}

func get{{name | pascalize}}TestCommand(ctx *{{name | camelize}}TestContext) *{{name | camelize}}Command {
	return &{{name | camelize}}Command{}
}

func Test_New{{name | pascalize}}CommandFactory(t *testing.T) {
	ctx := get{{name | pascalize}}TestContext(t)
	defer ctx.ctrl.Finish()

	f := New{{name | pascalize}}CommandFactory()
	cmd, err := f()

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if len(cmd.Synopsis()) == 0 {
		t.Error("Synopsis() should be present")
	}

	if len(cmd.Help()) == 0 {
		t.Error("Help() should be present")
	}
}

func Test_{{name | pascalize}}_Run(t *testing.T) {
	ctx := get{{name | pascalize}}TestContext(t)
	defer ctx.ctrl.Finish()

	cmd := get{{name | pascalize}}TestCommand(ctx)

	code := cmd.Run([]string{})

	if actual, expected := code, ui.ExitCodeOK; actual != expected {
		t.Errorf("Unexpected exit code %d, want %d", actual, expected)
	}
}
