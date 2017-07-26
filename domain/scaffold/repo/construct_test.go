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
	t         *testing.T
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
	edge            bool
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
	scffPath := filepath.Join(tmplsPath, "tmpl")

	return &constructTestContext{
		t:         t,
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
	entries := []fs.Entry{}
	for path, attrs := range entriesByPath {
		if !attrs.edge {
			entry, err := fs.NewEntry(filepath.Join(ctx.scffPath, path), attrs.dir)
			if err != nil {
				ctx.t.Fatalf("Unexpected error %v", err)
			}
			entries = append(entries, entry)
		}
	}
	ctx.fs.EXPECT().GetEntries(ctx.scffPath, true).
		Return(entries, nil)

	for path, entry := range entriesByPath {
		// Should skip meta.toml
		if path == "meta.toml" {
			continue
		}

		templateAbsPath := filepath.Join(ctx.scffPath, path)
		outputAbsPath := filepath.Join(ctx.rootPath, entry.outputPath)

		// Stubbing fs.Exists()
		if !entry.dir {
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
		if entry.dir {
			ctx.fs.EXPECT().CreateDir(outputAbsPath).
				Return(!entry.existing, nil).
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

	// ├── foo
	// │   ├── qux
	// │   │   └── corge.go
	// │   └── qux
	// │       └── grault.go
	// ├── bar
	// └── baz
	// │   ├── {{name}}
	// │   │   └── garlpy.go
	// │   ├── {{name}}.go
	// │   └── {{name}}_test.go
	// └── meta.toml
	// └── waldo.go
	entries := map[string]*constructTestEntry{
		"bar": {
			dir:        true,
			outputPath: "bar",
			existing:   false,
		},
		"baz": {
			dir:        true,
			edge:       true,
			outputPath: "baz",
			existing:   true,
		},
		"baz/{{name}}": {
			dir:        true,
			edge:       true,
			outputPath: fmt.Sprintf("baz/%s", ctx.name),
			existing:   false,
		},
		"baz/{{name}}.go": {
			dir:             false,
			template:        "// baz/{{name}}.go",
			outputPath:      fmt.Sprintf("baz/%s.go", ctx.name),
			outputContent:   fmt.Sprintf("// baz/%s.go", ctx.name),
			existing:        true,
			existingContent: fmt.Sprintf("// baz/%s.go", ctx.name),
		},
		"baz/{{name}}/garply.go": {
			dir:           false,
			template:      "// baz/{{name}}/garply.go",
			outputPath:    fmt.Sprintf("baz/%s/garply.go", ctx.name),
			outputContent: fmt.Sprintf("// baz/%s/garply.go", ctx.name),
			existing:      false,
		},
		"baz/{{name}}_test.go": {
			dir:             false,
			template:        "// baz/{{name}}_test.go",
			outputPath:      fmt.Sprintf("baz/%s_test.go", ctx.name),
			outputContent:   fmt.Sprintf("// baz/%s_test.go", ctx.name),
			existing:        true,
			existingContent: fmt.Sprintf("// baz/%s.go", ctx.name),
			overwriting:     true,
		},
		"meta.toml": {
			dir: false,
		},
		"foo/quux": {
			dir:        true,
			edge:       true,
			outputPath: "foo/quux",
			existing:   false,
		},
		"foo/quux/grault.go": {
			dir:           false,
			template:      "// foo/quux/grault.go",
			outputPath:    "foo/quux/grault.go",
			outputContent: "// foo/quux/grault.go",
			existing:      false,
		},
		"foo/qux": {
			dir:        true,
			edge:       true,
			outputPath: "foo/qux",
			existing:   false,
		},
		"foo/qux/corge.go": {
			dir:           false,
			template:      "// foo/qux/corge.go",
			outputPath:    "foo/qux/corge.go",
			outputContent: "// foo/qux/corge.go",
			existing:      false,
		},
		"waldo.go": {
			dir:           false,
			template:      "// waldo.go",
			outputPath:    "waldo.go",
			outputContent: "// waldo.go",
			existing:      false,
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

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	for path, entry := range entries {
		if path == "meta.toml" {
			continue
		}
		p := filepath.Join(ctx.rootPath, entry.outputPath)
		if c, ok := callbackCalls[p]; ok {
			if path == "meta.toml" {
				t.Error("meta.toml should be ignored")
			} else {
				expected := scaffold.ConstructSuccess
				if entry.existing && !entry.overwriting {
					expected = scaffold.ConstructSkipped
				}
				if entry.existing && entry.outputContent != entry.existingContent && !c.conflicted {
					t.Errorf("3rd argument ConstructCallback(%s, bool, bool, ConstructStatus) was %t, want %t", p, true, false)
				}
				if actual := c.status; actual != expected {
					t.Errorf("ConstructCallback(%s, %t, %s) was called, want (%s, %t, %s)", p, c.dir, actual, p, entry.dir, expected)
				}
			}
		} else {
			t.Errorf("ConstructCallback(%s, %t, ConstructStatus) should be called", p, entry.dir)
		}
	}
}
