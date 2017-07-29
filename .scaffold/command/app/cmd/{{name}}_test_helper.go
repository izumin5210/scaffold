package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/app/ui"
)

type {{name | camelize}}TestContext struct {
	ctrl *gomock.Controller
	ui   *ui.MockUI
}

func get{{name | pascalize}}TestContext(t *testing.T) *{{name | camelize}}TestContext {
	ctrl := gomock.NewController(t)
	return &{{name | camelize}}TestContext{
		ctrl: ctrl,
		ui:   ui.NewMockUI(ctrl),
	}
}

func get{{name | pascalize}}TestCommand(ctx *{{name | camelize}}TestContext) *{{name | camelize}}Command {
	return &{{name | camelize}}Command{
		ui: ctx.ui,
	}
}
