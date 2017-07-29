package usecase

import (
	"testing"

	"github.com/izumin5210/scaffold/app/usecase"
	"github.com/golang/mock/gomock"
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

func get{{name | pascalize}}TestUseCase(ctx *{{name | camelize}}TestContext) usecase.{{name | pascalize}}UseCase {
	return &{{name | camelize}}{}
}
