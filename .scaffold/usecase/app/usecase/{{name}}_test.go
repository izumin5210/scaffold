package usecase

import "testing"

func Test_{{name | pascalize}}UseCase_Perform(t *testing.T) {
	u := &{{name | camelize}}UseCase{}
	err := u.Perform()

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
