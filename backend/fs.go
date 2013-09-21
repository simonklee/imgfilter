package backend

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

// Dir implementation the ImageBackend
type Dir string

func (d Dir) ReadFile(name string) ([]byte, error) {
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 || strings.Contains(name, "\x00") {
		return nil, errors.New("http: invalid character in file path")
	}

	dir := string(d)

	if dir == "" {
		dir = "."
	}

	return ioutil.ReadFile(filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name))))
}
