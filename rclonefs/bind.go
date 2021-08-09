package rclonefs

import (
	"context"
	"github.com/rclone/rclone/fs"
	"github.com/spf13/afero"
)

func NewRCloneFs(section string) afero.Fs {
	if newFs, err := fs.NewFs(context.Background(), section+":"); err == nil {
		return New(newFs)
	} else {
		panic(err)
	}
}