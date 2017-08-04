package usecase

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/domain/scaffold"
)

func Test_GetScaffoldsUseCase_Perform(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := scaffold.NewMockRepository(ctrl)
	u := &getScaffoldsUseCase{repo: repo}

	testcases := []struct {
		scaffolds []scaffold.Scaffold
	}{
		{scaffolds: []scaffold.Scaffold{}},
		{scaffolds: []scaffold.Scaffold{
			scaffold.NewScaffold("/app/foo", &scaffold.Meta{}),
			scaffold.NewScaffold("/app/bar", &scaffold.Meta{}),
			scaffold.NewScaffold("/app/baz", &scaffold.Meta{}),
		}},
	}

	for _, testcase := range testcases {
		repo.EXPECT().GetScaffolds("/app").Return(testcase.scaffolds, nil).Times(1)
		scaffolds, err := u.Perform("/app")

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if actual, expected := len(scaffolds), len(testcase.scaffolds); actual != expected {
			t.Errorf("Return %d scaffolds, want %d", actual, expected)
		}
	}
}
