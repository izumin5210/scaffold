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

type constructCallbackCall struct {
	dir        bool
	conflicted bool
	status     scaffold.ConstructStatus
}

type constructTestEntry struct {
	dir             bool
	template        string
	outputPath      string
	outputContent   string
	existing        bool
	existingContent string
	overwriting     bool
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

func setupConstructTest(
	ctx *constructTestContext,
	entriesByPath map[string]*constructTestEntry,
) {
	// Stubbing fs.Walk()
	ctx.fs.EXPECT().Walk(ctx.scffPath, gomock.Any()).
		Do(func(_ string, cb func(path string, dir bool, err error) error) error {
			for path, entry := range entriesByPath {
				cb(filepath.Join(ctx.scffPath, path), entry.dir, nil)
			}
			return nil
		}).
		Times(1)

	for path, entry := range entriesByPath {
		// Should skip meta.toml
		if path == "meta.toml" {
			continue
		}

		templateAbsPath := filepath.Join(ctx.scffPath, path)
		outputAbsPath := filepath.Join(ctx.rootPath, entry.outputPath)

		// Stubbing fs.Exists() and fs.DirExists()
		if entry.dir {
			ctx.fs.EXPECT().DirExists(outputAbsPath).
				Return(entry.existing, nil).
				AnyTimes()
		} else {
			ctx.fs.EXPECT().Exists(outputAbsPath).
				Return(entry.existing, nil).
				AnyTimes()
		}

		// Stubbing fs.ReadFile() for templates
		if !entry.dir {
			ctx.fs.EXPECT().ReadFile(templateAbsPath).
				Return([]byte(entry.template), nil).
				AnyTimes()
			if entry.existing {
				ctx.fs.EXPECT().ReadFile(outputAbsPath).
					Return([]byte(entry.existingContent), nil).
					AnyTimes()
			}
		}

		// Stubbing fs.CreateDir()
		if entry.dir && !entry.existing {
			ctx.fs.EXPECT().CreateDir(outputAbsPath).
				Return(nil).
				Times(1)
		}

		// Stubbingg fs.CreateFile()
		if !entry.dir && (!entry.existing || (entry.existing && entry.overwriting)) {
			ctx.fs.EXPECT().CreateFile(outputAbsPath, entry.outputContent).
				Return(nil).
				Times(1)
		}
	}
}

func Test_Construct(t *testing.T) {
	ctx := getConstructTestContext(t)
	defer ctx.ctrl.Finish()

	entries := map[string]*constructTestEntry{
		"bar": {
			dir:        true,
			outputPath: "bar",
			existing:   false,
		},
		"bar/baz": {
			dir:           false,
			template:      "{{name}} baz",
			outputPath:    "bar/baz",
			outputContent: fmt.Sprintf("%s baz", ctx.name),
			existing:      false,
		},
		"bar/qux": {
			dir:        true,
			outputPath: "bar/qux",
			existing:   false,
		},
		"bar/qux/quux": {
			dir:           false,
			template:      "{{name}} quux",
			outputPath:    "bar/qux/quux",
			outputContent: fmt.Sprintf("%s quux", ctx.name),
			existing:      false,
		},
		"bar/qux/{{name}}": {
			dir:        true,
			outputPath: fmt.Sprintf("bar/qux/%s", ctx.name),
			existing:   false,
		},
		"bar/qux/{{name}}/{{name}}_type.go": {
			dir:           false,
			template:      "package {{name}}\n\n type {{name}}Type []string\n",
			outputPath:    fmt.Sprintf("bar/qux/%s/%s_type.go", ctx.name, ctx.name),
			outputContent: fmt.Sprintf("package %s\n\n type %sType []string\n", ctx.name, ctx.name),
			existing:      false,
		},
		"corge": {
			dir:           false,
			template:      "",
			outputPath:    "corge",
			outputContent: "",
			existing:      false,
		},
		"meta.toml": {
			dir: false,
		},
	}

	setupConstructTest(ctx, entries)

	callbackCalls := map[string]*constructCallbackCall{}

	err := ctx.repo.Construct(
		ctx.scaffold,
		ctx.name,
		func(path string, dir, conflicted bool, status scaffold.ConstructStatus) {
			if conflicted {
				t.Errorf("Unexpected conflict detected: %s", path)
			}
			callbackCalls[path] = &constructCallbackCall{dir: dir, conflicted: conflicted, status: status}
		},
		func(path, oldContent, newContent string) bool {
			t.Errorf("Unexpected conflict detected: %s", path)
			return false
		},
	)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if _, ok := callbackCalls["meta.toml"]; ok {
		t.Error("meta.toml should be ignored")
	}

	for path, entry := range entries {
		if path == "meta.toml" {
			continue
		}
		p := filepath.Join(ctx.rootPath, entry.outputPath)
		if c, ok := callbackCalls[p]; !ok {
			t.Errorf("ConstructCallback(%s, %t, %s) should be called", p, entry.dir, scaffold.ConstructSuccess)
		} else if !c.status.IsSuccess() {
			t.Errorf("ConstructCallback(%s, %t, %s) was called, want (%s, %t, %s)", p, c.dir, c.status, p, entry.dir, scaffold.ConstructSuccess)
		}
	}
}

func Test_Construct_FileExists(t *testing.T) {
	ctx := getConstructTestContext(t)
	defer ctx.ctrl.Finish()

	entries := map[string]*constructTestEntry{
		"bar": {
			dir:             false,
			template:        "{{name}} bar",
			outputPath:      "bar",
			outputContent:   fmt.Sprintf("%s bar", ctx.name),
			existing:        true,
			existingContent: fmt.Sprintf("%s bar", ctx.name),
		},
		"baz": {
			dir:             false,
			template:        "{{name}} baz",
			outputPath:      "baz",
			outputContent:   fmt.Sprintf("%s baz", ctx.name),
			existing:        true,
			existingContent: fmt.Sprintf("%s baz", ctx.name),
		},
	}

	setupConstructTest(ctx, entries)

	callbackCalls := map[string]*constructCallbackCall{}

	err := ctx.repo.Construct(
		ctx.scaffold,
		ctx.name,
		func(path string, dir, conflicted bool, status scaffold.ConstructStatus) {
			if conflicted {
				t.Errorf("Unexpected conflict detected: %s", path)
			}
			callbackCalls[path] = &constructCallbackCall{dir: dir, status: status}
		},
		func(path, oldContent, newContent string) bool {
			t.Errorf("Unexpected conflict detected: %s", path)
			return false
		},
	)

	for _, entry := range entries {
		p := filepath.Join(ctx.rootPath, entry.outputPath)
		if c, ok := callbackCalls[p]; ok {
			expected := scaffold.ConstructSuccess
			if entry.existing {
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

	entries := map[string]*constructTestEntry{
		"bar": {
			dir:        true,
			outputPath: "bar",
			existing:   true,
		},
		"bar/baz": {
			dir:           false,
			template:      "{{name}} baz",
			outputPath:    "bar/baz",
			outputContent: fmt.Sprintf("%s baz", ctx.name),
			existing:      false,
		},
		"qux": {
			dir:        true,
			outputPath: "qux",
			existing:   false,
		},
		"qux/quux": {
			dir:           false,
			template:      "{{name}} quux",
			outputPath:    "qux/quux",
			outputContent: fmt.Sprintf("%s quux", ctx.name),
			existing:      false,
		},
	}

	callbackCalls := map[string]*constructCallbackCall{}

	setupConstructTest(ctx, entries)

	err := ctx.repo.Construct(
		ctx.scaffold,
		ctx.name,
		func(path string, dir, conflicted bool, status scaffold.ConstructStatus) {
			if conflicted {
				t.Errorf("Unexpected conflict detected: %s", path)
			}
			callbackCalls[path] = &constructCallbackCall{dir: dir, status: status}
		},
		func(path, oldContent, newContent string) bool {
			t.Errorf("Unexpected conflict detected: %s", path)
			return false
		},
	)

	for _, entry := range entries {
		p := filepath.Join(ctx.rootPath, entry.outputPath)
		if c, ok := callbackCalls[p]; ok {
			expected := scaffold.ConstructSuccess
			if entry.existing {
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

func Test_Construct_WhenDetectConflicts(t *testing.T) {
	ctx := getConstructTestContext(t)
	defer ctx.ctrl.Finish()

	entries := map[string]*constructTestEntry{
		"bar": {
			dir:             false,
			template:        "{{name}} bar",
			outputPath:      "bar",
			outputContent:   fmt.Sprintf("%s bar", ctx.name),
			existing:        true,
			existingContent: fmt.Sprintf("%s bar", ctx.name),
		},
		"baz": {
			dir:             false,
			template:        "{{name}} baz",
			outputPath:      "baz",
			outputContent:   fmt.Sprintf("%s baz", ctx.name),
			existing:        true,
			existingContent: "bazbaz",
			overwriting:     true,
		},
		"qux": {
			dir:             false,
			template:        "{{name}} qux",
			outputPath:      "qux",
			outputContent:   fmt.Sprintf("%s qux", ctx.name),
			existing:        true,
			existingContent: "quxqux",
			overwriting:     false,
		},
	}

	setupConstructTest(ctx, entries)

	callbackCalls := map[string]*constructCallbackCall{}

	err := ctx.repo.Construct(
		ctx.scaffold,
		ctx.name,
		func(path string, dir, conflicted bool, status scaffold.ConstructStatus) {
			callbackCalls[path] = &constructCallbackCall{dir: dir, conflicted: conflicted, status: status}
		},
		func(path, oldContent, newContent string) bool {
			relpath, err := filepath.Rel(ctx.rootPath, path)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}
			if entry, ok := entries[relpath]; ok {
				if oldContent != entry.existingContent {
					t.Errorf("2nd argument was %q, want %q", oldContent, entry.existingContent)
				}
				if newContent != entry.outputContent {
					t.Errorf("2nd argument was %q, want %q", newContent, entry.outputContent)
				}
				if entry.outputContent == entry.existingContent {
					t.Errorf("Conflicted content should be different with old content: %q", path)
				}
				return entry.overwriting
			}
			return true
		},
	)

	for _, entry := range entries {
		p := filepath.Join(ctx.rootPath, entry.outputPath)
		if c, ok := callbackCalls[p]; ok {
			expected := scaffold.ConstructSuccess
			if entry.existing && !entry.overwriting {
				expected = scaffold.ConstructSkipped
			}
			if entry.outputContent != entry.existingContent && !c.conflicted {
				t.Errorf("3rd argument ConstructCallback(%s, bool, bool, ConstructStatus) was %t, want %t", p, true, false)
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
