package scaffolds

import (
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/entity"
	"github.com/izumin5210/scaffold/infra/fs"
	"github.com/pkg/errors"
)

func Test_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFS := fs.NewMockFS(ctrl)

	context := &entity.Context{ScaffoldsPath: "/app/.scaffold", FS: mockFS}
	repo := NewRepository(context)

	testcases := []struct {
		dirs  []string
		metas map[string]string
		out   []*entity.Scaffold
	}{
		{dirs: []string{}},
		{
			dirs: []string{"/app/.scaffold/foo", "/app/.scaffold/bar", "/app/.scaffold/baz"},
			metas: map[string]string{
				"/app/.scaffold/foo": "synopsis = \"\"\"\nGenerates foo\n\"\"\"",
				"/app/.scaffold/baz": "synopsis = \"\"\"\nGenerates baz\n\"\"\"",
			},
			out: []*entity.Scaffold{
				entity.NewScaffold("/app/.scaffold/bar", &entity.ScaffoldMeta{Synopsis: "Geenrates bar"}),
				entity.NewScaffold("/app/.scaffold/baz", &entity.ScaffoldMeta{}),
				entity.NewScaffold("/app/.scaffold/foo", &entity.ScaffoldMeta{Synopsis: "Geenrates foo"}),
			},
		},
	}

	for _, testcase := range testcases {
		mockFS.EXPECT().GetDirs("/app/.scaffold").Return(testcase.dirs, nil).Times(1)
		for _, dir := range testcase.dirs {
			tomlpath := path.Join(dir, "meta.toml")
			if meta, ok := testcase.metas[dir]; ok {
				mockFS.EXPECT().ReadFile(tomlpath).Return([]byte(meta), nil).Times(1)
			} else {
				mockFS.EXPECT().ReadFile(tomlpath).Return(nil, errors.Errorf("%s does not exist.", tomlpath)).Times(1)
			}
		}
		names, err := repo.GetAll()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if actual, expected := len(names), len(testcase.dirs); actual != expected {
			t.Errorf("GetAll() returns %d items, but expected %d items", actual, expected)
		}

		for i, s := range testcase.out {
			if actual, expected := s, testcase.out[i]; actual != expected {
				t.Errorf("GetAll()[%d] is %v, but expected %v", i, actual, expected)
			}
		}
	}
}
