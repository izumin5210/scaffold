package repo

import (
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/pkg/errors"

	"github.com/izumin5210/scaffold/domain/scaffold"
	repotesting "github.com/izumin5210/scaffold/infra/scaffold/repo/testing"
)

func Test_Create(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)

	cases := []struct {
		entry         scaffold.Entry
		created       bool
		parentCreated bool
	}{
		{
			entry:         scaffold.NewEntry("/foo/bar.go", "package main", false),
			created:       true,
			parentCreated: true,
		},
		{
			entry:         scaffold.NewEntry("/foo/bar", "", true),
			created:       true,
			parentCreated: true,
		},
		{
			entry:         scaffold.NewEntry("/foo.go", "package main", false),
			created:       true,
			parentCreated: false,
		},
		{
			entry:         scaffold.NewEntry("/foo", "", true),
			created:       true,
			parentCreated: false,
		},
		{
			entry:         scaffold.NewEntry("/foo", "", true),
			created:       false,
			parentCreated: false,
		},
	}

	for _, c := range cases {
		ctx.FS.EXPECT().CreateDir(filepath.Dir(c.entry.Path())).Return(c.parentCreated, nil)
		if c.entry.IsDir() {
			ctx.FS.EXPECT().CreateDir(c.entry.Path()).Return(c.created, nil)
		} else {
			ctx.FS.EXPECT().CreateFile(c.entry.Path(), c.entry.Content()).Return(nil)
		}
		created, parentCreated, err := repo.Create(c.entry)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if actual, expected := created, c.created; actual != expected {
			t.Errorf("Create() returns %t at 1st, want %t", actual, expected)
		}

		if actual, expected := parentCreated, c.parentCreated; actual != expected {
			t.Errorf("Create() returns %t at 2nd, want %t", actual, expected)
		}
	}
}

func Test_Create_WhenFailedToCreateParent(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)

	entry := scaffold.NewEntry("/app/foo.go", "package main", false)
	ctx.FS.EXPECT().CreateDir("/app").Return(false, errors.New("error"))

	created, parentCreated, err := repo.Create(entry)

	if err == nil {
		t.Error("Should return an error")
	}

	if actual, expected := created, false; actual != expected {
		t.Errorf("Create() returns %t at 1st, want %t", actual, expected)
	}

	if actual, expected := parentCreated, false; actual != expected {
		t.Errorf("Create() returns %t at 2nd, want %t", actual, expected)
	}
}

func Test_Create_WhenFailedToCreateDir(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)

	entry := scaffold.NewEntry("/app/foo", "", true)
	ctx.FS.EXPECT().CreateDir("/app").Return(false, nil)
	ctx.FS.EXPECT().CreateDir(entry.Path()).Return(false, errors.New("error"))

	created, parentCreated, err := repo.Create(entry)

	if err == nil {
		t.Error("Should return an error")
	}

	if actual, expected := created, false; actual != expected {
		t.Errorf("Create() returns %t at 1st, want %t", actual, expected)
	}

	if actual, expected := parentCreated, false; actual != expected {
		t.Errorf("Create() returns %t at 2nd, want %t", actual, expected)
	}
}

func Test_Create_WhenFailedToCreateFile(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)

	entry := scaffold.NewEntry("/app/foo.go", "package main", false)
	ctx.FS.EXPECT().CreateDir("/app").Return(true, nil)
	ctx.FS.EXPECT().CreateFile(entry.Path(), gomock.Any()).Return(errors.New("error"))

	created, parentCreated, err := repo.Create(entry)

	if err == nil {
		t.Error("Should return an error")
	}

	if actual, expected := created, false; actual != expected {
		t.Errorf("Create() returns %t at 1st, want %t", actual, expected)
	}

	if actual, expected := parentCreated, true; actual != expected {
		t.Errorf("Create() returns %t at 2nd, want %t", actual, expected)
	}
}
