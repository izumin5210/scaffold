package repo

import (
	"path/filepath"
	"strings"

	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/pkg/errors"
)

func (r *repo) GetConcreteEntries(
	s scaffold.Scaffold,
	tmpls []scaffold.TemplateEntry,
	v interface{},
) (map[string]scaffold.ConcreteEntry, error) {
	entries := map[string]scaffold.ConcreteEntry{}
	for _, tmpl := range tmpls {
		path, err := scaffold.TemplateString(tmpl.Path()).Compile(tmpl.Path(), v)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to compile path: %q", tmpl.Path())
		}
		path, _ = filepath.Rel(s.Path(), path)
		if strings.Index(path, "..") == 0 {
			return nil, errors.Errorf("%q is not contained scaffold %q", tmpl.Path(), s.Path())
		}
		path = filepath.Join(r.rootPath, path)
		if tmpl.IsDir() {
			if existing, err := r.fs.DirExists(path); err != nil {
				return nil, errors.Wrapf(err, "Failed to read %q", path)
			} else if existing {
				entries[tmpl.Path()] = scaffold.NewConcreteDir(path, tmpl.Path())
			}
		} else {
			if existing, err := r.fs.Exists(path); err != nil {
				return nil, errors.Wrapf(err, "Failed to read %q", path)
			} else if existing {
				content, err := r.fs.ReadFile(path)
				if err != nil {
					return nil, errors.Wrapf(err, "Failed to read %q", path)
				}
				entries[tmpl.Path()] = scaffold.NewConcreteFile(path, string(content), tmpl.Path())
			}
		}
	}
	return entries, nil
}
