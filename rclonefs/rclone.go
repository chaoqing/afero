// Copyright © 2021 Nicolas Wang <cqwang@uw.edu>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rclonefs

import (
	"context"

	"os"
	"time"

	_ "github.com/rclone/rclone/backend/all"
	"github.com/rclone/rclone/fs"
	_ "github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/config/configfile"
	"github.com/rclone/rclone/vfs"

	"github.com/spf13/afero"
)

type Fs struct {
	*vfs.VFS
}

func New(client fs.Fs) afero.Fs {
	return &Fs{vfs.New(client, nil)}
}

func NewRCloneFs(config string) afero.Fs {
	if fs, err := fs.NewFs(context.Background(), config+":"); err == nil {
		return New(fs)
	} else {
		panic(err)
	}
}

func (s Fs) Name() string { return "rclonefs" }

func (s Fs) Create(name string) (afero.File, error) {
	return s.VFS.Create(name)
}

func (s Fs) MkdirAll(path string, perm os.FileMode) error {
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := s.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return err
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent
		err = s.MkdirAll(path[0:j-1], perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = s.Mkdir(path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := s.Stat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}

func (s Fs) Open(name string) (afero.File, error) {
	return s.VFS.Open(name)
}

func (s Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return s.VFS.OpenFile(name, flag, perm)
}

func (s Fs) Remove(name string) error {
	return s.VFS.Remove(name)
}

func (s Fs) RemoveAll(path string) error {

	// TODO have a look at os.RemoveAll
	// https://github.com/golang/go/blob/master/src/os/path.go#L66
	return nil
}

func (s Fs) Rename(oldname, newname string) error {
	return s.VFS.Rename(oldname, newname)
}

func (s Fs) Stat(name string) (os.FileInfo, error) {
	return s.VFS.Stat(name)
}

func (s Fs) Chmod(name string, mode os.FileMode) error {
	//return s.VFS.Chmod(name, mode)
	//todo
	 return nil
}

func (s Fs) Chown(name string, uid, gid int) error {
	//return s.VFS.Chown(name, uid, gid)
	//todo
	return nil
}

func (s Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return s.VFS.Chtimes(name, atime, mtime)
}

func init(){
	//config.SetConfigPath("~/.config/rclone/rclone.conf")
	//config.ClearConfigPassword()
	configfile.Install()
}
