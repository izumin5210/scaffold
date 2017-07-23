package repo

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/izumin5210/scaffold/domain/scaffold"
)

func (r *repo) Construct(scff scaffold.Scaffold, name string, cb scaffold.ConstructCallback) error {
	evaluator := &templateEvaluator{funcMap: template.FuncMap{
		"name": func() string { return name },
	}}
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
		outputPath, err := evaluator.evaluate(path, filepath.Join(r.rootPath, relpath))
		if err != nil {
			return err
		}

		if dir {
			if exists, err := r.fs.DirExists(outputPath); exists || err != nil {
				if exists && err == nil {
					cb(outputPath, true, scaffold.ConstructSkipped)
				}
				return err
			}
			err = r.fs.CreateDir(outputPath)
			if err == nil {
				cb(outputPath, true, scaffold.ConstructSuccess)
			}
			return err
		}

		if exists, err := r.fs.Exists(outputPath); exists || err != nil {
			if exists && err == nil {
				cb(outputPath, false, scaffold.ConstructSkipped)
			}
			return err
		}

		raw, err := r.fs.ReadFile(path)
		if err != nil {
			return err
		}
		content, err := evaluator.evaluate(path, string(raw))
		if err != nil {
			return err
		}
		err = r.fs.CreateFile(outputPath, content)
		if err == nil {
			cb(outputPath, false, scaffold.ConstructSuccess)
		}
		return err
	})
	return err
}

type templateEvaluator struct {
	funcMap template.FuncMap
}

func (e *templateEvaluator) evaluate(name string, text string) (string, error) {
	tmpl := template.New(name)
	tmpl.Funcs(e.funcMap)
	tmpl, err := tmpl.Parse(text)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, struct{}{})
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}
