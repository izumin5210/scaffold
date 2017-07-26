package repo

import (
	"path/filepath"

	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/pkg/errors"
)

func (r *repo) Construct(
	scff scaffold.Scaffold,
	name string,
	cb scaffold.ConstructCallback,
	conflctedCb scaffold.ConstructConflictedCallback,
) error {
	tmpl := scaffold.NewTemplate(name)
	metaPath := filepath.Join(scff.Path(), "meta.toml")
	entries, err := r.fs.GetEntries(scff.Path(), true)
	if err != nil {
		return errors.Wrapf(err, "Failed to get entries under %q", scff.Path())
	}

	createdDirs := map[string]struct{}{}

	for _, entry := range entries {
		if entry.Path() == metaPath {
			continue
		}
		relpath, err := filepath.Rel(scff.Path(), entry.Path())
		if err != nil {
			return errors.Cause(err)
		}
		outputPath, err := tmpl.Compile(filepath.Join(r.rootPath, relpath))
		if err != nil {
			return errors.Cause(err)
		}

		// construct directory
		if entry.IsDir() {
			if _, ok := createdDirs[outputPath]; !ok {
				if err := r.createDir(outputPath, cb); err != nil {
					return errors.Cause(err)
				}
				createdDirs[outputPath] = struct{}{}
			}
			if err != nil {
				return err
			}
			continue
		}

		// construct file
		content, err := r.getCompiledContent(tmpl, entry.Path())
		if err != nil {
			return errors.Cause(err)
		}
		conflicted := false

		parentDir := filepath.Dir(outputPath)
		if parentDir != r.rootPath {
			if _, ok := createdDirs[parentDir]; !ok {
				if err := r.createDir(parentDir, cb); err != nil {
					return errors.Cause(err)
				}
				createdDirs[parentDir] = struct{}{}
			}
		}

		if exists, err := r.fs.Exists(outputPath); err != nil {
			return errors.Cause(err)
		} else if exists {
			existing, err := r.fs.ReadFile(outputPath)
			if err != nil {
				return errors.Cause(err)
			}
			existingContent := string(existing)
			conflicted = content != existingContent
			if !conflicted || !conflctedCb(outputPath, existingContent, content) {
				cb(outputPath, false, conflicted, scaffold.ConstructSkipped)
				continue
			}
		}

		if err = r.fs.CreateFile(outputPath, content); err != nil {
			if conflicted {
				return errors.Wrapf(err, "Failed to overwrite %q", outputPath)
			}
			return errors.Wrapf(err, "Failed to create %q", outputPath)
		}
		cb(outputPath, false, conflicted, scaffold.ConstructSuccess)
	}
	return errors.Cause(err)
}

func (r *repo) createDir(path string, cb scaffold.ConstructCallback) error {
	ok, err := r.fs.CreateDir(path)
	if err != nil {
		return errors.Wrapf(err, "Failed to create directory %q", path)
	}
	if ok {
		cb(path, true, false, scaffold.ConstructSuccess)
	} else {
		cb(path, true, false, scaffold.ConstructSkipped)
	}
	return nil
}

func (r *repo) getCompiledContent(tmpl scaffold.Template, path string) (string, error) {
	raw, err := r.fs.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to read template %q", path)
	}
	compiled, err := tmpl.Compile(string(raw))
	if err != nil {
		return "", errors.Wrapf(err, "Failed to complile template %q", path)
	}
	return string(compiled), nil
}
