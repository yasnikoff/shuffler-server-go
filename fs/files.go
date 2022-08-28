package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// NewTextFile is like NewFile but with string as content.
func NewTextFile(path string, content string) error {
	return NewFile(path, []byte(content))
}

// NewFile creates file with specified path and content.
func NewFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModeDir); err != nil {
		return err
	}
	return ioutil.WriteFile(path, content, 0777)
}

