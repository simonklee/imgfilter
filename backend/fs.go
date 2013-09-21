package backend

import (
	"io/ioutil"
	"path/filepath"
)

// FileSystem is a file system implementation of the ImageBackend
type FileSystem struct {
	basepath string
}

func NewFileSystem(basepath string) *FileSystem {
	return &FileSystem {
		basepath: basepath,
	}
}

func (fs *FileSystem) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(fs.basepath, filename))
}
