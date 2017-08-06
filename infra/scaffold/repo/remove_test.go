package repo

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/pkg/errors"

	"github.com/izumin5210/scaffold/infra/fs"

	"github.com/izumin5210/scaffold/domain/scaffold"
	repotesting "github.com/izumin5210/scaffold/infra/scaffold/repo/testing"
)

func Test_Remove(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)

	cases := []struct {
		setup func()
		in    scaffold.Entry
		out   []string
	}{
		{
			setup: func() {
				ctx.FS.EXPECT().Remove("/foo.go").Return(nil)
			},
			in:  scaffold.NewEntry("/foo.go", "package main", false),
			out: []string{"/foo.go"},
		},
		{
			setup: func() {
				ctx.FS.EXPECT().Remove("/foo/bar/baz.go").Return(nil)
				ctx.FS.EXPECT().GetEntries("/foo/bar", true).Return([]fs.Entry{}, nil)
				ctx.FS.EXPECT().Remove("/foo/bar").Return(nil)
				ctx.FS.EXPECT().GetEntries("/foo", true).Return([]fs.Entry{}, nil)
				ctx.FS.EXPECT().Remove("/foo").Return(nil)
			},
			in:  scaffold.NewEntry("/foo/bar/baz.go", "package bar", false),
			out: []string{"/foo/bar/baz.go", "/foo/bar", "/foo"},
		},
		{
			setup: func() {
				ctx.FS.EXPECT().Remove("/foo/bar.go").Return(nil)
				ctx.FS.EXPECT().GetEntries("/foo", true).Return([]fs.Entry{}, nil)
				ctx.FS.EXPECT().Remove("/foo").Return(nil)
			},
			in:  scaffold.NewEntry("/foo/bar.go", "package foo", false),
			out: []string{"/foo/bar.go", "/foo"},
		},
		{
			setup: func() {
				ctx.FS.EXPECT().Remove("/foo/bar.go").Return(nil)
				ctx.FS.EXPECT().GetEntries("/foo", true).
					Return([]fs.Entry{
						func() fs.Entry { e, _ := fs.NewEntry("/foo/baz.go", false); return e }(),
					}, nil)
			},
			in:  scaffold.NewEntry("/foo/bar.go", "package foo", false),
			out: []string{"/foo/bar.go"},
		},
	}

	for _, c := range cases {
		c.setup()

		entries, err := repo.Remove(c.in)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if got, want := len(entries), len(c.out); got != want {
			t.Errorf("Remove() returned %d items, want %d items", got, want)
		}

		for i, got := range entries {
			if want := c.out[i]; !reflect.DeepEqual(want, got) {
				t.Errorf("Remove()[%d] was %v, want %v", i, got, want)
			}
		}
	}
}

func Test_Remove_WhenRemoveFailed(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	entry := scaffold.NewEntry("/foo/bar.go", "package foo", false)

	ctx.FS.EXPECT().Remove(entry.Path()).Return(errors.New("error"))

	removed, err := repo.Remove(entry)

	if got, want := len(removed), 0; got != want {
		t.Errorf("Remove() returned %d items, want %d items", got, want)
	}

	if err == nil {
		t.Error("Remove() should return an error")
	}
}

func Test_Remove_WhenGetEntriesFailed(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	entry := scaffold.NewEntry("/foo/bar.go", "package foo", false)

	ctx.FS.EXPECT().Remove(entry.Path()).Return(nil)
	ctx.FS.EXPECT().GetEntries("/foo", gomock.Any()).Return(nil, errors.New("error"))

	removed, err := repo.Remove(entry)

	if got, want := len(removed), 1; got != want {
		t.Errorf("Remove() returned %d items, want %d items", got, want)
	}

	if got, want := removed[0], entry.Path(); got != want {
		t.Errorf("Remove()[0] returned %s, want %s", got, want)
	}

	if err == nil {
		t.Error("Remove() should return an error")
	}
}

func Test_Remove_WhenRemoveParentFailed(t *testing.T) {
	ctx := repotesting.NewRepoTestContext(t)
	defer ctx.Ctrl.Finish()

	repo := New(ctx.RootPath, ctx.TmplsPath, ctx.FS)
	entry := scaffold.NewEntry("/foo/bar.go", "package foo", false)

	ctx.FS.EXPECT().Remove(entry.Path()).Return(nil)
	ctx.FS.EXPECT().GetEntries("/foo", gomock.Any()).Return([]fs.Entry{}, nil)
	ctx.FS.EXPECT().Remove("/foo").Return(errors.New("error"))

	removed, err := repo.Remove(entry)

	if got, want := len(removed), 1; got != want {
		t.Errorf("Remove() returned %d items, want %d items", got, want)
	}

	if got, want := removed[0], entry.Path(); got != want {
		t.Errorf("Remove()[0] returned %s, want %s", got, want)
	}

	if err == nil {
		t.Error("Remove() should return an error")
	}
}
