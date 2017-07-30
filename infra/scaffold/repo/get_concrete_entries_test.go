package repo

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"github.com/izumin5210/scaffold/domain/scaffold"
	repotesting "github.com/izumin5210/scaffold/infra/scaffold/repo/testing"
)

func Test_GetConcreteEntries(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	name := "gopher"
	v := struct{ Name string }{Name: name}

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	scff := scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "foo"), &scaffold.Meta{})

	toTS := func(str string) scaffold.TemplateString {
		return scaffold.TemplateString(str)
	}
	getTmplPath := func(relpath string) scaffold.TemplateString {
		return toTS(filepath.Join(scff.Path(), relpath))
	}

	// app/.scaffold/foo
	// ├── bar
	// │   ├── qux
	// │   │   └── corge.go
	// │   └── qux
	// │       └── grault.go
	// ├── {{name}}.go
	// └── baz
	//     ├── {{name}}
	//     │   └──garply
	//     ├── {{name}}.go
	//     └── {{name}}_test.go
	cases := []struct {
		in    []scaffold.TemplateEntry
		out   map[string]scaffold.ConcreteEntry
		files map[string]string
		dirs  map[string]struct{}
	}{
		{},
		{
			in: []scaffold.TemplateEntry{
				scaffold.NewTemplateFile(getTmplPath("{{name}}.go"), ""),
				scaffold.NewTemplateFile(getTmplPath("bar.go"), ""),
			},
			out: map[string]scaffold.ConcreteEntry{},
		},
		{
			in: []scaffold.TemplateEntry{
				scaffold.NewTemplateFile(getTmplPath("{{name}}.go"), ""),
				scaffold.NewTemplateFile(getTmplPath("bar/quux/grault.go"), ""),
				scaffold.NewTemplateFile(getTmplPath("bar/qux/corge.go"), ""),
				scaffold.NewTemplateFile(getTmplPath("baz/{{name}}.go"), ""),
				scaffold.NewTemplateDir(getTmplPath("baz/{{name}}/garply")),
				scaffold.NewTemplateFile(getTmplPath("baz/{{name}}_test.go"), ""),
			},
			out: map[string]scaffold.ConcreteEntry{
				fmt.Sprintf("%s.go", name): scaffold.NewConcreteEntry(
					filepath.Join(ctx.RootPath, fmt.Sprintf("%s.go", name)),
					"package main",
					false,
					string(getTmplPath("{{name}}.go")),
				),
				"bar/qux/corge.go": scaffold.NewConcreteEntry(
					filepath.Join(ctx.RootPath, "bar/qux/corge.go"),
					"package qux",
					false,
					string(getTmplPath("bar/qux/corge.go")),
				),
				fmt.Sprintf("baz/%s.go", name): scaffold.NewConcreteEntry(
					filepath.Join(ctx.RootPath, fmt.Sprintf("baz/%s.go", name)),
					"package baz\n\ntype Gopher interface{}",
					false,
					string(getTmplPath("baz/{{name}}.go")),
				),
				fmt.Sprintf("baz/%s/garply", name): scaffold.NewConcreteEntry(
					filepath.Join(ctx.RootPath, fmt.Sprintf("baz/%s/garply", name)),
					"",
					true,
					string(getTmplPath("baz/{{name}}/garply")),
				),
			},
			files: map[string]string{
				fmt.Sprintf("%s.go", name):     "package main",
				"bar/qux/corge.go":             "package qux",
				fmt.Sprintf("baz/%s.go", name): "package baz\n\ntype Gopher interface{}",
			},
			dirs: map[string]struct{}{
				fmt.Sprintf("baz/%s/garply", name): {},
			},
		},
	}

	for _, c := range cases {
		for _, tmpl := range c.in {
			compiledPath, _ := scaffold.TemplateString(tmpl.Path()).Compile(tmpl.Path(), v)
			relpath, _ := filepath.Rel(scff.Path(), compiledPath)
			abspath := filepath.Join(ctx.RootPath, relpath)
			if tmpl.IsDir() {
				_, ok := c.dirs[relpath]
				ctx.FS.EXPECT().DirExists(abspath).Return(ok, nil)
			} else {
				if content, ok := c.files[relpath]; ok {
					ctx.FS.EXPECT().Exists(abspath).Return(true, nil)
					ctx.FS.EXPECT().ReadFile(abspath).Return([]byte(content), nil)
				} else {
					ctx.FS.EXPECT().Exists(abspath).Return(false, nil)
				}
			}
		}
		entries, err := repo.GetConcreteEntries(scff, c.in, v)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if actual, expected := len(entries), len(c.out); actual != expected {
			t.Errorf("GetConcreteEntries() returns %d items, but expected %d items", actual, expected)
		}

		for path, entry := range entries {
			relpath, _ := filepath.Rel(ctx.RootPath, entry.Path())
			if actual, expected := entry.Path(), c.out[relpath].Path(); actual != expected {
				t.Errorf("GetConcreteEntries()[%q].Path() is %q, but expected %q", path, actual, expected)
			}
			if actual, expected := entry.Content(), c.out[relpath].Content(); actual != expected {
				t.Errorf("GetConcreteEntries()[%q].Content() is %q, but expected %q", path, actual, expected)
			}
			if actual, expected := entry.IsDir(), c.out[relpath].IsDir(); actual != expected {
				t.Errorf("GetConcreteEntries()[%q].IsDir() is %t, but expected %t", path, actual, expected)
			}
			if actual, expected := entry.TemplatePath(), c.out[relpath].TemplatePath(); actual != expected {
				t.Errorf("GetConcreteEntries()[%q].Path() is %q, but expected %q", path, actual, expected)
			}
		}
	}
}

func Test_GetConcreteEntries_WhenFailedToCompilePath(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	name := "gopher"
	v := struct{ Name string }{Name: name}

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	scff := scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "foo"), &scaffold.Meta{})

	tmpls := []scaffold.TemplateEntry{
		scaffold.NewTemplateDir(scaffold.TemplateString(filepath.Join(scff.Path(), "{{name}_test"))),
	}
	entries, err := repo.GetConcreteEntries(scff, tmpls, v)

	if err == nil {
		t.Error("Should return an error")
	}

	if entries != nil {
		t.Errorf("GetConcreteEntries() returns %v, want nil", entries)
	}
}

func Test_GetConcreteEntries_WhenTemplateOutsideScaffold(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	name := "gopher"
	v := struct{ Name string }{Name: name}

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	scff := scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "foo"), &scaffold.Meta{})

	tmpls := []scaffold.TemplateEntry{
		scaffold.NewTemplateDir(scaffold.TemplateString("/app/{{name}}_test")),
	}
	entries, err := repo.GetConcreteEntries(scff, tmpls, v)

	if err == nil {
		t.Error("Should return an error")
	}

	if entries != nil {
		t.Errorf("GetConcreteEntries() returns %v, want nil", entries)
	}
}

func Test_GetConcreteEntries_WhenFailedToCheckDirExistence(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	name := "gopher"
	v := struct{ Name string }{Name: name}

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	scff := scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "foo"), &scaffold.Meta{})

	tmpls := []scaffold.TemplateEntry{
		scaffold.NewTemplateDir(scaffold.TemplateString(filepath.Join(scff.Path(), "{{name}}_test"))),
	}
	ctx.FS.EXPECT().DirExists(gomock.Any()).Return(false, errors.New("error"))
	entries, err := repo.GetConcreteEntries(scff, tmpls, v)

	if err == nil {
		t.Error("Should return an error")
	}

	if entries != nil {
		t.Errorf("GetConcreteEntries() returns %v, want nil", entries)
	}
}

func Test_GetConcreteEntries_WhenFailedToCheckFileExistence(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	name := "gopher"
	v := struct{ Name string }{Name: name}

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	scff := scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "foo"), &scaffold.Meta{})

	tmpls := []scaffold.TemplateEntry{
		scaffold.NewTemplateFile(
			scaffold.TemplateString(filepath.Join(scff.Path(), "{{name}}_test")),
			scaffold.TemplateString("package main"),
		),
	}
	ctx.FS.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
	entries, err := repo.GetConcreteEntries(scff, tmpls, v)

	if err == nil {
		t.Error("Should return an error")
	}

	if entries != nil {
		t.Errorf("GetConcreteEntries() returns %v, want nil", entries)
	}
}

func Test_GetConcreteEntries_WhenFailedToCreadFile(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	name := "gopher"
	v := struct{ Name string }{Name: name}

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	scff := scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "foo"), &scaffold.Meta{})

	tmpls := []scaffold.TemplateEntry{
		scaffold.NewTemplateFile(
			scaffold.TemplateString(filepath.Join(scff.Path(), "{{name}}_test")),
			scaffold.TemplateString("package main"),
		),
	}
	ctx.FS.EXPECT().Exists(gomock.Any()).Return(true, nil)
	ctx.FS.EXPECT().ReadFile(gomock.Any()).Return(nil, errors.New("error"))
	entries, err := repo.GetConcreteEntries(scff, tmpls, v)

	if err == nil {
		t.Error("Should return an error")
	}

	if entries != nil {
		t.Errorf("GetConcreteEntries() returns %v, want nil", entries)
	}
}
