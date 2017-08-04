package repo

import (
	"path/filepath"
	"testing"

	"github.com/pkg/errors"

	"github.com/izumin5210/scaffold/infra/fs"

	"github.com/izumin5210/scaffold/domain/scaffold"
	repotesting "github.com/izumin5210/scaffold/infra/scaffold/repo/testing"
)

func Test_GetTemplates(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	scff := scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "foo"), &scaffold.Meta{})

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
		entries      []fs.Entry
		contents     map[string]string
		excludeCount int
	}{
		{},
		{
			entries: []fs.Entry{
				repotesting.NewFSEntry(filepath.Join(scff.Path(), "{{name}}.go"), false),
				repotesting.NewFSEntry(filepath.Join(scff.Path(), "bar/quux/grault.go"), false),
				repotesting.NewFSEntry(filepath.Join(scff.Path(), "bar/qux/corge.go"), false),
				repotesting.NewFSEntry(filepath.Join(scff.Path(), "baz/{{name}}.go"), false),
				repotesting.NewFSEntry(filepath.Join(scff.Path(), "baz/{{name}}/garply"), true),
				repotesting.NewFSEntry(filepath.Join(scff.Path(), "baz/{{name}}_test.go"), false),
				repotesting.NewFSEntry(filepath.Join(scff.Path(), "meta/meta.toml"), false),
				repotesting.NewFSEntry(filepath.Join(scff.Path(), "meta.toml"), false),
			},
			contents: map[string]string{
				"{{name}}.go":          "package main",
				"bar/quux/grault.go":   "package quux",
				"bar/qux/corge.go":     "package qux",
				"baz/{{name}}.go":      "package baz\n\ntype {{name | pascalize}} interface{}",
				"baz/{{name}}_test.go": "package baz",
				"meta/meta.toml":       "",
			},
			excludeCount: 1,
		},
	}

	for _, c := range cases {
		ctx.FS.EXPECT().GetEntries(scff.Path(), true).Return(c.entries, nil).Times(1)
		for path, content := range c.contents {
			ctx.FS.EXPECT().ReadFile(filepath.Join(scff.Path(), path)).Return([]byte(content), nil)
		}
		tmpls, err := repo.GetTemplates(scff)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if actual, expected := len(tmpls), len(c.entries)-c.excludeCount; actual != expected {
			t.Errorf("GetTemplates() returns %d items, but expected %d items", actual, expected)
		}

		for i, tmpl := range tmpls {
			if actual, expected := tmpl.Path(), c.entries[i].Path(); actual != expected {
				t.Errorf("GetTemplates()[%d].Path() is %q, but expected %q", i, actual, expected)
			}
			if actual, expected := tmpl.IsDir(), c.entries[i].IsDir(); actual != expected {
				t.Errorf("GetTemplates()[%d].IsDir() is %t, but expected %t", i, actual, expected)
			}
		}
	}
}

func Test_GetTemplates_WhenGetEntriesFailed(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	scff := scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "foo"), &scaffold.Meta{})

	ctx.FS.EXPECT().GetEntries(scff.Path(), true).Return(nil, errors.New("error"))

	tmpls, err := repo.GetTemplates(scff)

	if err == nil {
		t.Error("Should return an error")
	}

	if tmpls != nil {
		t.Errorf("GetTemplates() returns %v, want nil", tmpls)
	}
}

func Test_GetTemplates_WhenReadFileFailed(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	scff := scaffold.NewScaffold(filepath.Join(ctx.TmplsPath, "foo"), &scaffold.Meta{})

	entries := []fs.Entry{
		repotesting.NewFSEntry(filepath.Join(scff.Path(), "bar.go"), false),
	}
	ctx.FS.EXPECT().GetEntries(scff.Path(), true).Return(entries, nil)
	ctx.FS.EXPECT().ReadFile(entries[0].Path()).Return(nil, errors.New("error"))

	tmpls, err := repo.GetTemplates(scff)

	if err == nil {
		t.Error("Should return an error")
	}

	if tmpls != nil {
		t.Errorf("GetTemplates() returns %v, want nil", tmpls)
	}
}
