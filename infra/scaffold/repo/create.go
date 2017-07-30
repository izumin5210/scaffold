package repo

import (
	"path/filepath"

	"github.com/izumin5210/scaffold/domain/scaffold"
	"github.com/pkg/errors"
)

func (r *repo) Create(e scaffold.Entry) (bool, bool, error) {
	parent := filepath.Dir(e.Path())
	parentCreated, err := r.fs.CreateDir(parent)
	if err != nil {
		return false, parentCreated, errors.Wrapf(err, "Failed to create directory %q", parent)
	}
	created := false
	if e.IsDir() {
		created, err = r.fs.CreateDir(e.Path())
		if err != nil {
			return created, parentCreated, errors.Wrapf(err, "Failed to create directory %q", e.Path())
		}
	} else {
		err = r.fs.CreateFile(e.Path(), e.Content())
		if err != nil {
			return created, parentCreated, errors.Wrapf(err, "Failed to create file %q", e.Path())
		}
		created = true
	}
	return created, parentCreated, nil
}
