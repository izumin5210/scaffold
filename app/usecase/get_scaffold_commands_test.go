package usecase

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/app/cmd/factory"
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/mitchellh/cli"
)

func Test_GetScaffoldCommandUseCase_Perform(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := scaffold.NewMockRepository(ctrl)
	factory := factory.NewMockFactory(ctrl)
	u := &getScaffoldCommandUseCase{repo: repo, factory: factory}

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
		repo.EXPECT().GetAll().Return(testcase.scaffolds, nil).Times(1)
		for _, sc := range testcase.scaffolds {
			factory.EXPECT().CreateCreateScaffoldCommandFactory(sc).
				Return(func() (cli.Command, error) { return &cli.MockCommand{}, nil }).
				Times(1)
		}
		factories, err := u.Perform()

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		for _, sc := range testcase.scaffolds {
			if _, ok := factories[sc.Name()]; !ok {
				t.Errorf("Returned command factories should include %q", sc.Name())
			}
		}
	}
}
