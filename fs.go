package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/blang/vfs"
)

type Filesystem struct {
	vfs.Filesystem
}

func OSFileSystem() *Filesystem {
	return NewFilesystem(vfs.OS())
}

func NewFilesystem(vfs vfs.Filesystem) *Filesystem {
	return &Filesystem{vfs}
}

func (f *Filesystem) Create(name string) (vfs.File, error) {
	return f.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
}

func (f *Filesystem) Open(name string) (vfs.File, error) {
	return f.OpenFile(name, os.O_RDWR, 0)
}

func (f *Filesystem) ReadFile(name string) ([]byte, error) {
	r, err := f.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	defer r.Close()
	return ioutil.ReadAll(r)
}

func RelativePath(base, rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(base, rel)
}
