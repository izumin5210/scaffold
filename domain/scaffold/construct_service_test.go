package scaffold

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

type constructServiceTestContext struct {
	ctrl     *gomock.Controller
	repo     *MockRepository
	root     string
	name     string
	scaffold Scaffold
	service  ConstructService
}

type constructCallbackCall struct {
	dir        bool
	conflicted bool
	status     ConstructStatus
}

func newConstructServiceTestContext(t *testing.T) *constructServiceTestContext {
	ctrl := gomock.NewController(t)
	repo := NewMockRepository(ctrl)
	root := "/app"
	sc := NewScaffold(filepath.Join(root, ".scaffold", "golang"), &Meta{})
	name := "gopher"
	return &constructServiceTestContext{
		ctrl:     ctrl,
		repo:     repo,
		root:     root,
		scaffold: sc,
		name:     name,
		service:  NewConstructService(repo),
	}
}

func Test_ConstructService(t *testing.T) {
	ctx := newConstructServiceTestContext(t)
	defer ctx.ctrl.Finish()

	getTmplPath := func(relpath string) TemplateString {
		return TemplateString(filepath.Join(ctx.scaffold.Path(), relpath))
	}
	getConcPath := func(relpath string) string {
		return filepath.Join(ctx.root, relpath)
	}

	cases := []struct {
		tmpls             []TemplateEntry
		values            interface{}
		existings         map[string]ConcreteEntry
		createdWithParent map[string]bool
		calls             map[string]*constructCallbackCall
		overwriting       map[string]struct{}
	}{
		{},
		{
			tmpls: []TemplateEntry{
				NewTemplateFile(getTmplPath("foo.go"), "package main", ctx.scaffold.Path()),
			},
			createdWithParent: map[string]bool{
				getConcPath("foo.go"): false,
			},
			values: struct{ Name string }{Name: ctx.name},
			calls: map[string]*constructCallbackCall{
				getConcPath("foo.go"): &constructCallbackCall{status: ConstructSuccess},
			},
		},
		// ├── foo
		// │   ├── qux
		// │   │   └── corge.go
		// │   └── quux
		// │       └── grault.go
		// ├── bar
		// └── baz
		// │   ├── {{name}}
		// │   │   ├── garlpy.go
		// │   │   └── {{name}}.go
		// │   ├── {{name}}.go
		// │   └── {{name}}_test.go
		// └── meta.toml
		// └── waldo.go
		{
			tmpls: []TemplateEntry{
				NewTemplateDir(getTmplPath("bar"), ctx.scaffold.Path()),
				NewTemplateFile(getTmplPath("baz/{{name}}.go"), "package baz", ctx.scaffold.Path()),
				NewTemplateFile(getTmplPath("baz/{{name}}_test.go"), "package baz", ctx.scaffold.Path()),
				NewTemplateFile(getTmplPath("baz/{{name}}/garlpy.go"), "package {{name}}", ctx.scaffold.Path()),
				NewTemplateFile(getTmplPath("baz/{{name}}/{{name}}.go"), "package {{name}}", ctx.scaffold.Path()),
				NewTemplateFile(getTmplPath("foo/quux/grault.go"), "package quux", ctx.scaffold.Path()),
				NewTemplateFile(getTmplPath("foo/qux/corge.go"), "package qux", ctx.scaffold.Path()),
				NewTemplateFile(getTmplPath("waldo.go"), "package main", ctx.scaffold.Path()),
			},
			createdWithParent: map[string]bool{
				getConcPath("bar"):                                           false,
				getConcPath(fmt.Sprintf("baz/%s_test.go", ctx.name)):         false,
				getConcPath(fmt.Sprintf("baz/%s/garlpy.go", ctx.name)):       true,
				getConcPath(fmt.Sprintf("baz/%s/%s.go", ctx.name, ctx.name)): false,
				getConcPath("foo/quux/grault.go"):                            true,
				getConcPath("waldo.go"):                                      false,
			},
			values: struct{ Name string }{Name: ctx.name},
			existings: map[string]ConcreteEntry{
				// existed
				getConcPath("foo/qux/corge.go"): NewConcreteFile(
					getConcPath("foo/qux/corge.go"),
					"package qux",
					string(getTmplPath("foo/qux/corge.go")),
				),
				// conflicted
				getConcPath(fmt.Sprintf("baz/%s.go", ctx.name)): NewConcreteFile(
					getConcPath(fmt.Sprintf("baz/%s.go", ctx.name)),
					"",
					string(getTmplPath("baz/{{name}}.go")),
				),
				// conflicted and will overwrite
				getConcPath(fmt.Sprintf("baz/%s_test.go", ctx.name)): NewConcreteFile(
					getConcPath(fmt.Sprintf("baz/%s_test.go", ctx.name)),
					"",
					string(getTmplPath("baz/{{name}}_test.go")),
				),
			},
			calls: map[string]*constructCallbackCall{
				getConcPath("bar"):                                           &constructCallbackCall{dir: true, status: ConstructSuccess},
				getConcPath("baz"):                                           &constructCallbackCall{dir: true, status: ConstructSkipped},
				getConcPath(fmt.Sprintf("baz/%s.go", ctx.name)):              &constructCallbackCall{conflicted: true, status: ConstructSkipped},
				getConcPath(fmt.Sprintf("baz/%s_test.go", ctx.name)):         &constructCallbackCall{conflicted: true, status: ConstructSuccess},
				getConcPath(fmt.Sprintf("baz/%s", ctx.name)):                 &constructCallbackCall{dir: true, status: ConstructSuccess},
				getConcPath(fmt.Sprintf("baz/%s/garlpy.go", ctx.name)):       &constructCallbackCall{status: ConstructSuccess},
				getConcPath(fmt.Sprintf("baz/%s/%s.go", ctx.name, ctx.name)): &constructCallbackCall{status: ConstructSuccess},
				getConcPath("foo/quux"):                                      &constructCallbackCall{dir: true, status: ConstructSuccess},
				getConcPath("foo/quux/grault.go"):                            &constructCallbackCall{status: ConstructSuccess},
				getConcPath("foo/qux/corge.go"):                              &constructCallbackCall{status: ConstructSkipped},
				getConcPath("waldo.go"):                                      &constructCallbackCall{status: ConstructSuccess},
			},
			overwriting: map[string]struct{}{
				getConcPath(fmt.Sprintf("baz/%s_test.go", ctx.name)): struct{}{},
			},
		},
	}

	for _, c := range cases {
		ctx.repo.EXPECT().GetTemplates(ctx.scaffold).
			Return(c.tmpls, nil)
		ctx.repo.EXPECT().
			GetConcreteEntries(ctx.scaffold, c.tmpls, c.values).
			Return(c.existings, nil)

		for _, tmpl := range c.tmpls {
			conc, _ := tmpl.Compile(ctx.root, c.values)
			if dirCreated, ok := c.createdWithParent[conc.Path()]; ok {
				ctx.repo.EXPECT().Create(conc).Return(true, dirCreated, nil)
			}
		}

		calls := map[string]*constructCallbackCall{}

		err := ctx.service.Perform(
			ctx.root,
			ctx.scaffold,
			c.values,
			func(path string, dir, conflicted bool, status ConstructStatus) {
				calls[path] = &constructCallbackCall{dir: dir, conflicted: conflicted, status: status}
			},
			func(path, oldContent, newContent string) bool {
				if oldContent == newContent {
					t.Errorf("%s is conflicted but old one and new one has same content", path)
				}
				_, ok := c.overwriting[path]
				return ok
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
