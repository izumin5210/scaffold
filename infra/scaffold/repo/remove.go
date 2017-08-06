package repo

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/izumin5210/scaffold/domain/scaffold"
)

func (r *repo) Remove(e scaffold.Entry) ([]string, error) {
	if err := r.fs.Remove(e.Path()); err != nil {
		return nil, errors.Wrapf(err, "Failed to remove %s", e.Path())
	}
	removed := []string{e.Path()}
	for path := filepath.Dir(e.Path()); path != "/"; path = filepath.Dir(path) {
		if entries, err := r.fs.GetEntries(path, true); err != nil {
			return removed, errors.Wrapf(err, "Failed to get %s children", path)
		} else if len(entries) > 0 {
			break
		}
		if err := r.fs.Remove(path); err != nil {
			return removed, errors.Wrapf(err, "Failed to remove %s", path)
		}
		removed = append(removed, path)
	}
	return removed, nil
}
