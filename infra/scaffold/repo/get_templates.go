package repo

import (
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/pkg/errors"
)

func (r *repo) GetTemplates(s scaffold.Scaffold) ([]scaffold.TemplateEntry, error) {
	entries, err := r.fs.GetEntries(s.Path(), true)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get templates under %q", s.Path())
	}
	var tmpls []scaffold.TemplateEntry
	for _, e := range entries {
		tmplPath := scaffold.TemplateString(e.Path())
		if e.IsDir() {
			tmpls = append(tmpls, scaffold.NewTemplateDir(tmplPath))
		} else {
			content, err := r.fs.ReadFile(e.Path())
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to read template %q", e.Path())
			}
			tmpl := scaffold.NewTemplateFile(tmplPath, scaffold.TemplateString(content))
			tmpls = append(tmpls, tmpl)
		}
	}
	return tmpls, nil
}
