package repo

import (
	"path/filepath"

	"github.com/izumin5210/scaffold/domain/scaffold"
)

func (r *repo) Construct(
	scff scaffold.Scaffold,
	name string,
	cb scaffold.ConstructCallback,
	conflctedCb scaffold.ConstructConflictedCallback,
) error {
	tmpl := scaffold.NewTemplate(name)
	metaPath := filepath.Join(scff.Path(), "meta.toml")
	err := r.fs.Walk(scff.Path(), func(path string, dir bool, err error) error {
		if err != nil {
			return err
		}
		if path == metaPath {
			return nil
		}

		relpath, err := filepath.Rel(scff.Path(), path)
		if err != nil {
			return err
		}
		outputPath, err := tmpl.Compile(filepath.Join(r.rootPath, relpath))
		if err != nil {
			return err
		}

		if dir {
			if exists, err := r.fs.DirExists(outputPath); exists || err != nil {
				if exists && err == nil {
					cb(outputPath, true, false, scaffold.ConstructSkipped)
				}
				return err
			}
			err = r.fs.CreateDir(outputPath)
			if err == nil {
				cb(outputPath, true, false, scaffold.ConstructSuccess)
			}
			return err
		}

		raw, err := r.fs.ReadFile(path)
		if err != nil {
			return err
		}
		compiled, err := tmpl.Compile(string(raw))
		if err != nil {
			return err
		}
		content := string(compiled)
		conflicted := false

		if exists, err := r.fs.Exists(outputPath); exists || err != nil {
			if err != nil {
				return err
			}
			if exists {
				existing, err := r.fs.ReadFile(outputPath)
				if err != nil {
					return err
				}
				existingContent := string(existing)
				conflicted = content != existingContent
				if !conflicted || !conflctedCb(outputPath, existingContent, content) {
					cb(outputPath, false, conflicted, scaffold.ConstructSkipped)
					return nil
				}
			}
		}

		err = r.fs.CreateFile(outputPath, content)
		if err == nil {
			cb(outputPath, false, conflicted, scaffold.ConstructSuccess)
		}
		return err
	})
	return err
}
