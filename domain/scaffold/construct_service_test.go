package scaffold

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

type constructServiceTestContext struct {
	ctrl    *gomock.Controller
	repo    *MockRepository
	root    string
	service ConstructService
}

type constructCallbackCall struct {
	dir        bool
	conflicted bool
	status     ConstructStatus
}

func newConstructServiceTestContext(t *testing.T) *constructServiceTestContext {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	return &constructServiceTestContext{
		ctrl:    ctrl,
		repo:    repo,
		root:    "/app",
		service: NewConstructService(repo),
	}
}

func Test_ConstructService(t *testing.T) {
	ctx := newConstructServiceTestContext(t)
	defer ctx.ctrl.Finish()

	sc := NewScaffold(filepath.Join(ctx.root, ".scaffold", "foo"), &Meta{})
	name := "gopher"

	getTmplPath := func(relpath string) TemplateString {
		return TemplateString(filepath.Join(sc.Path(), relpath))
	}
	getConcPath := func(relpath string) string {
		return filepath.Join(ctx.root, relpath)
	}

	cases := []struct {
		tmpls     []TemplateEntry
		values    interface{}
		existings map[string]ConcreteEntry
		calls     map[string]*constructCallbackCall
	}{
		{},
		{
			tmpls: []TemplateEntry{
				NewTemplateFile(getTmplPath("bar.go"), "package main"),
			},
			values: struct{ Name string }{Name: name},
			calls: map[string]*constructCallbackCall{
				getConcPath("bar.go"): &constructCallbackCall{status: ConstructSuccess},
			},
		},
	}

	for _, c := range cases {
		ctx.repo.EXPECT().GetTemplates(sc).
			Return(c.tmpls, nil)
		ctx.repo.EXPECT().
			GetConcreteEntries(sc, c.tmpls, c.values).
			Return(c.existings, nil)

		calls := map[string]*constructCallbackCall{}

		err := ctx.service.Perform(
			sc,
			c.values,
			func(path string, dir, conflicted bool, status ConstructStatus) {
				calls[path] = &constructCallbackCall{dir: dir, conflicted: conflicted, status: status}
			},
			func(path, oldContent, newContent string) bool {
				return false
			},
		)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if got, want := len(calls), len(c.calls); got != want {
			t.Errorf("Callback called %d times, want %d times", got, want)
		}

		for path, want := range c.calls {
			if got, ok := calls[path]; !ok {
				t.Errorf("Callback for %q was missed", path)
			} else if !reflect.DeepEqual(got, want) {
				t.Errorf("Call for %q was %v, want %v", path, got, want)
			}
		}
	}
}
