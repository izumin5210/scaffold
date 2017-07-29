package usecase

import (
	"testing"
)

func Test_{{name | pascalize}}UseCase_Perform(t *testing.T) {
	ctx := get{{name | pascalize}}TestContext(t)
	defer ctx.ctrl.Finish()

	u := get{{name | pascalize}}TestUseCase(ctx)
	err := u.Perform()

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
