package rclonefs

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"testing"
)

func TestNewRCloneFs(t *testing.T) {
	fs := NewRCloneFs("ibm")

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Error("walkFn err:", err)
		}
		var size int64
		if !info.IsDir() {
			size = info.Size()
		}
		fmt.Println(path, info.Name(), size, info.IsDir(), err)
		return nil
	}

	afero.Walk(fs, "/", walkFn)
}