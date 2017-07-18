package scaffolds

import (
	"path"

	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/izumin5210/scaffold/entity"
	"github.com/pkg/errors"
)

func (r *repo) GetAll() ([]*entity.Scaffold, error) {
	dirs, err := r.context.FS.GetDirs(r.context.ScaffoldsPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get directories from %q", err)
	}
	var scaffolds []*entity.Scaffold
	for _, dir := range dirs {
		scPath, err := filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
		var meta entity.ScaffoldMeta
		data, err := r.context.FS.ReadFile(path.Join(scPath, "meta.toml"))
		if err == nil {
			toml.Decode(string(data), &meta)
		}
		scaffolds = append(scaffolds, entity.NewScaffold(scPath, &meta))
	}
	return scaffolds, nil
}
