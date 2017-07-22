package repo

import (
	"path"

	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/pkg/errors"
)

func (r *repo) GetAll() ([]scaffold.Scaffold, error) {
	dirs, err := r.fs.GetDirs(r.tmplsPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get directories from %q", err)
	}
	var scaffolds []scaffold.Scaffold
	for _, dir := range dirs {
		scPath := filepath.Join(r.tmplsPath, dir)
		var meta scaffold.Meta
		data, err := r.fs.ReadFile(path.Join(scPath, "meta.toml"))
		if err == nil {
			toml.Decode(string(data), &meta)
		}
		scaffolds = append(scaffolds, scaffold.NewScaffold(scPath, &meta))
	}
	return scaffolds, nil
}
