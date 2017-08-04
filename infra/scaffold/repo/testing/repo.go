package testing

import (
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izumin5210/scaffold/infra/fs"
)

// RepoTestContext contains common objects for repository testings
type RepoTestContext struct {
	Ctrl      *gomock.Controller
	FS        *fs.MockFS
	RootPath  string
	TmplsPath string
}

// NewRepoTestContext returns common objects container for repository testings
func NewRepoTestContext(t *testing.T) *RepoTestContext {
	ctrl := gomock.NewController(t)
	fs := fs.NewMockFS(ctrl)
	rootPath := "/app"
	tmplsPath := filepath.Join(rootPath, ".scaffold")
	return &RepoTestContext{
		Ctrl:      ctrl,
		FS:        fs,
		RootPath:  rootPath,
		TmplsPath: tmplsPath,
	}
}

// NewFSEntry creates a fs.Entry object ignoring errors
func NewFSEntry(path string, dir bool) fs.Entry {
	e, _ := fs.NewEntry(path, dir)
	return e
}
