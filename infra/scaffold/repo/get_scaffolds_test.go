package repo

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/izumin5210/scaffold/domain/scaffold"
	repotesting "github.com/izumin5210/scaffold/infra/scaffold/repo/testing"
	"github.com/pkg/errors"
)

func Test_GetScaffolds(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)

	cases := []struct {
		dirs  []string
		metas map[string]string
		out   []scaffold.Scaffold
	}{
		{dirs: []string{}},
		{
			dirs: []string{"foo", "bar", "baz"},
			metas: map[string]string{
				"foo": "synopsis = \"\"\"\nGenerates foo\n\"\"\"",
				"baz": "synopsis = \"\"\"\nGenerates baz\n\"\"\"",
			},
			out: []scaffold.Scaffold{
				scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "foo"), &scaffold.Meta{Synopsis: "Geenrates foo"}),
				scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "bar"), &scaffold.Meta{Synopsis: "Geenrates bar"}),
				scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "baz"), &scaffold.Meta{}),
			},
		},
	}

	for _, c := range cases {
		ctx.FS.EXPECT().GetDirs(ctx.TmplsPath).Return(c.dirs, nil).Times(1)
		for _, dir := range c.dirs {
			tomlpath := filepath.Join(ctx.TmplsPath, dir, "meta.toml")
			if meta, ok := c.metas[dir]; ok {
				ctx.FS.EXPECT().ReadFile(tomlpath).Return([]byte(meta), nil).Times(1)
			} else {
				ctx.FS.EXPECT().ReadFile(tomlpath).Return(nil, errors.Errorf("%s does not exist.", tomlpath)).Times(1)
			}
		}
		scffs, err := repo.GetScaffolds(ctx.TmplsPath)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if actual, expected := len(scffs), len(c.dirs); actual != expected {
			t.Errorf("GetScaffolds() returns %d items, but expected %d items", actual, expected)
		}

		for i, s := range scffs {
			if actual, expected := s, c.out[i]; reflect.DeepEqual(actual, expected) {
				t.Errorf("GetScaffolds()[%d] is %v, but expected %v", i, actual, expected)
			}
		}
	}
}

func Test_GetScaffolds_WhenFailToGetDirs(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)

	ctx.FS.EXPECT().GetDirs(ctx.TmplsPath).
		Return(nil, errors.New("error"))

	scffs, err := repo.GetScaffolds(ctx.TmplsPath)

	if scffs != nil {
		t.Error("Should not return a scaffold list")
	}

	if err == nil {
		t.Error("Should return an error")
	}
}
