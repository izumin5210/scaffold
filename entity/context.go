package entity

import infra "github.com/izumin5210/scaffold/infra/fs"

// Context is container storing configurations
type Context struct {
	ScaffoldsPath string
	FS            infra.FS
}
