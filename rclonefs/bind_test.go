package rclonefs

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"testing"
)


func walkPrintFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	var size int64
	if !info.IsDir() {
		size = info.Size()
	}
	fmt.Println(path, info.Name(), size, info.IsDir(), err)
	return nil
}

func TestNewRCloneFs(t *testing.T) {
	fs := NewRCloneFs("ibm")

	_ = afero.Walk(fs, "/", walkPrintFn)
}

func TestNewBindPathFs(t *testing.T) {
	fsRoot := afero.NewMemMapFs()
	fsHome := afero.NewMemMapFs()

	_ = fsRoot.MkdirAll("/tmp/test", os.ModeDir)
	_ = afero.WriteFile(fsRoot, "/tmp/test/hello.go", []byte("import os"), os.ModePerm)

	bindFs := NewBindPathFs(map[string]afero.Fs{
		"/": fsRoot,
		"/home": fsHome,
	})

	fmt.Println("list all files inside bind filesystem")
	_ = afero.Walk(bindFs, "/", walkPrintFn)
}