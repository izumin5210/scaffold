package usecase

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/cmd/factory"
	"github.com/izumin5210/scaffold/entity"
	"github.com/izumin5210/scaffold/repo/scaffolds"
	"github.com/mitchellh/cli"
)

func Test_GetScaffoldCommandUseCase_Perform(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := scaffolds.NewMockRepository(ctrl)
	factory := factory.NewMockFactory(ctrl)
	u := &getScaffoldCommandUseCase{repo: repo, factory: factory}

	testcases := []struct {
		scaffolds []*entity.Scaffold
	}{
		{scaffolds: []*entity.Scaffold{}},
		{scaffolds: []*entity.Scaffold{
			entity.NewScaffold("/app/foo", &entity.ScaffoldMeta{}),
			entity.NewScaffold("/app/bar", &entity.ScaffoldMeta{}),
			entity.NewScaffold("/app/baz", &entity.ScaffoldMeta{}),
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
