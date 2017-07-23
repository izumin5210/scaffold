package repo

import (
	"testing"

	"path/filepath"

	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/izumin5210/scaffold/infra/fs"
)

type constructTestContext struct {
	ctrl      *gomock.Controller
	fs        *fs.MockFS
	rootPath  string
	tmplsPath string
	scffPath  string
	name      string
	repo      scaffold.Repository
	scaffold  scaffold.Scaffold
}

func getConstructTestContext(t *testing.T) *constructTestContext {
	ctrl := gomock.NewController(t)
	fs := fs.NewMockFS(ctrl)
	rootPath := "/app"
	tmplsPath := filepath.Join(rootPath, ".scaffold")
	scffPath := filepath.Join(tmplsPath, "foo")

	return &constructTestContext{
		ctrl:      ctrl,
		fs:        fs,
		rootPath:  rootPath,
		tmplsPath: tmplsPath,
		scffPath:  scffPath,
		name:      "gopher",
		repo:      New(rootPath, tmplsPath, fs),
		scaffold:  scaffold.NewScaffold(scffPath, &scaffold.Meta{}),
	}
}

func Test_Construct(t *testing.T) {
	ctx := getConstructTestContext(t)
	defer ctx.ctrl.Finish()

	entries := []struct {
		path string
		dir  bool
	}{
		{path: "bar", dir: true},
		{path: "bar/baz", dir: false},
		{path: "bar/qux", dir: true},
		{path: "bar/qux/quux", dir: false},
		{path: "bar/qux/{{name}}", dir: true},
		{path: "bar/qux/{{name}}/{{name}}_type.go", dir: false},
		{path: "corge", dir: false},
		{path: "meta.toml", dir: false},
	}
	contents := map[string]string{
		"bar/baz":                           "{{name}} baz",
		"bar/qux/quux":                      "{{name}} quux",
		"bar/qux/{{name}}/{{name}}_type.go": "package {{name}}\n\n type {{name}}Type []string\n",
		"corge": "",
	}
	dirs := []string{
		"bar",
		"bar/qux",
		fmt.Sprintf("bar/qux/%s", ctx.name),
	}
	outputs := map[string]string{
		"bar/baz":      fmt.Sprintf("%s baz", ctx.name),
		"bar/qux/quux": fmt.Sprintf("%s quux", ctx.name),
		fmt.Sprintf("bar/qux/%s/%s_type.go", ctx.name, ctx.name): fmt.Sprintf("package %s\n\n type %sType []string\n", ctx.name, ctx.name),
		"corge": "",
	}

	ctx.fs.EXPECT().Walk(ctx.scffPath, gomock.Any()).
		Do(func(_ string, cb func(path string, dir bool, err error) error) error {
			for _, entry := range entries {
				cb(filepath.Join(ctx.scffPath, entry.path), entry.dir, nil)
			}
			return nil
		}).
		Times(1)
	ctx.fs.EXPECT().DirExists(gomock.Any()).Return(false, nil).AnyTimes()
	ctx.fs.EXPECT().Exists(gomock.Any()).Return(false, nil).AnyTimes()
	for path, content := range contents {
		ctx.fs.EXPECT().ReadFile(filepath.Join(ctx.scffPath, path)).Return([]byte(content), nil)
	}
	for _, dir := range dirs {
		ctx.fs.EXPECT().CreateDir(filepath.Join(ctx.rootPath, dir)).Return(nil)
	}
	for path, content := range outputs {
		ctx.fs.EXPECT().CreateFile(filepath.Join(ctx.rootPath, path), content).Return(nil)
	}

	type callbackCall struct {
		dir    bool
		status scaffold.ConstructStatus
	}
	callbackCalls := map[string]*callbackCall{}

	err := ctx.repo.Construct(
		ctx.scaffold,
		ctx.name,
		func(path string, dir bool, status scaffold.ConstructStatus) {
			callbackCalls[path] = &callbackCall{dir: dir, status: status}
		},
	)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	for _, dir := range dirs {
		p := filepath.Join(ctx.rootPath, dir)
		if c, ok := callbackCalls[p]; !ok {
			t.Errorf("ConstructCallback(%s, %t, %s) should be called", p, true, scaffold.ConstructSuccess)
		} else if !c.dir || !c.status.IsSuccess() {
			t.Errorf("ConstructCallback(%s, %t, %s) was called, want (%s, %t, %s)", p, c.dir, c.status, p, true, scaffold.ConstructSuccess)
		}
	}

	for relpath := range outputs {
		p := filepath.Join(ctx.rootPath, relpath)
		if c, ok := callbackCalls[p]; !ok {
			t.Errorf("ConstructCallback(%s, %t, %s) should be called", p, false, scaffold.ConstructSuccess)
		} else if c.dir || !c.status.IsSuccess() {
			t.Errorf("ConstructCallback(%s, %t, %s) was called, want (%s, %t, %s)", p, c.dir, c.status, p, false, scaffold.ConstructSuccess)
		}
	}
}

func Test_Construct_FileExists(t *testing.T) {
	ctx := getConstructTestContext(t)
	defer ctx.ctrl.Finish()

	entries := []struct {
		path   string
		dir    bool
		exists bool
	}{
		{path: "bar", dir: false, exists: true},
		{path: "baz", dir: false, exists: false},
	}
	contents := map[string]string{
		"baz": "{{name}} baz",
	}
	outputs := map[string]string{
		"baz": fmt.Sprintf("%s baz", ctx.name),
	}

	ctx.fs.EXPECT().Walk(ctx.scffPath, gomock.Any()).
		Do(func(_ string, cb func(path string, dir bool, err error) error) error {
			for _, entry := range entries {
				cb(filepath.Join(ctx.scffPath, entry.path), entry.dir, nil)
			}
			return nil
		}).
		Times(1)
	for _, entry := range entries {
		path := filepath.Join(ctx.rootPath, entry.path)
		if entry.dir {
			ctx.fs.EXPECT().DirExists(path).Return(entry.exists, nil).Times(1)
		} else {
			ctx.fs.EXPECT().Exists(path).Return(entry.exists, nil).Times(1)
		}
	}
	for path, content := range contents {
		ctx.fs.EXPECT().ReadFile(filepath.Join(ctx.scffPath, path)).
			Return([]byte(content), nil).
			Times(1)
	}
	for path, content := range outputs {
		ctx.fs.EXPECT().CreateFile(filepath.Join(ctx.rootPath, path), content).
			Return(nil).
			Times(1)
	}

	type callbackCall struct {
		dir    bool
		status scaffold.ConstructStatus
	}
	callbackCalls := map[string]*callbackCall{}

	err := ctx.repo.Construct(
		ctx.scaffold,
		ctx.name,
		func(path string, dir bool, status scaffold.ConstructStatus) {
			callbackCalls[path] = &callbackCall{dir: dir, status: status}
		},
	)

	for _, entry := range entries {
		p := filepath.Join(ctx.rootPath, entry.path)
		if c, ok := callbackCalls[p]; ok {
			expected := scaffold.ConstructSuccess
			if entry.exists {
				expected = scaffold.ConstructSkipped
			}
			if actual := c.status; actual != expected {
				t.Errorf("ConstructCallback(%s, %t, %s) was called, want (%s, %t, %s)", p, c.dir, actual, p, entry.dir, expected)
			}
		} else {
			t.Errorf("ConstructCallback(%s, %t, ConstructStatus) should be called", p, entry.dir)
		}
	}

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func Test_Construct_DirExists(t *testing.T) {
	ctx := getConstructTestContext(t)
	defer ctx.ctrl.Finish()

	entries := []struct {
		path   string
		dir    bool
		exists bool
	}{
		{path: "bar", dir: true, exists: true},
		{path: "bar/baz", dir: false, exists: false},
		{path: "qux", dir: true, exists: false},
		{path: "qux/quux", dir: false, exists: false},
	}
	contents := map[string]string{
		"bar/baz":  "{{name}} baz",
		"qux/quux": "{{name}} quux",
	}
	dirs := []string{
		"qux",
	}
	outputs := map[string]string{
		"bar/baz":  fmt.Sprintf("%s baz", ctx.name),
		"qux/quux": fmt.Sprintf("%s quux", ctx.name),
	}

	ctx.fs.EXPECT().Walk(ctx.scffPath, gomock.Any()).
		Do(func(_ string, cb func(path string, dir bool, err error) error) error {
			for _, entry := range entries {
				cb(filepath.Join(ctx.scffPath, entry.path), entry.dir, nil)
			}
			return nil
		}).
		Times(1)
	for _, entry := range entries {
		path := filepath.Join(ctx.rootPath, entry.path)
		if entry.dir {
			ctx.fs.EXPECT().DirExists(path).Return(entry.exists, nil).Times(1)
		} else {
			ctx.fs.EXPECT().Exists(path).Return(entry.exists, nil).Times(1)
		}
	}
	for path, content := range contents {
		ctx.fs.EXPECT().ReadFile(filepath.Join(ctx.scffPath, path)).
			Return([]byte(content), nil).
			Times(1)
	}
	for _, dir := range dirs {
		ctx.fs.EXPECT().CreateDir(filepath.Join(ctx.rootPath, dir)).
			Return(nil).
			Times(1)
	}
	for path, content := range outputs {
		ctx.fs.EXPECT().CreateFile(filepath.Join(ctx.rootPath, path), content).
			Return(nil).
			Times(1)
	}

	type callbackCall struct {
		dir    bool
		status scaffold.ConstructStatus
	}
	callbackCalls := map[string]*callbackCall{}

	err := ctx.repo.Construct(
		ctx.scaffold,
		ctx.name,
		func(path string, dir bool, status scaffold.ConstructStatus) {
			callbackCalls[path] = &callbackCall{dir: dir, status: status}
		},
	)

	for _, entry := range entries {
		p := filepath.Join(ctx.rootPath, entry.path)
		if c, ok := callbackCalls[p]; ok {
			expected := scaffold.ConstructSuccess
			if entry.exists {
				expected = scaffold.ConstructSkipped
			}
			if actual := c.status; actual != expected {
				t.Errorf("ConstructCallback(%s, %t, %s) was called, want (%s, %t, %s)", p, c.dir, actual, p, entry.dir, expected)
			}
		} else {
			t.Errorf("ConstructCallback(%s, %t, ConstructStatus) should be called", p, entry.dir)
		}
	}

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}
