package shuffler

import (
	"path/filepath"
	"os"
	"bitbucket.org/yasnikoff/shuffler/internal/md"
)

type DocHandler func(path string, d Doc, err error) error

// TODO: walk is inefficient for very large trees because it uses sorting for the result to be deterministic.
func Walk(root string, handler DocHandler) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if info.Name() == md.DirName {
				return filepath.SkipDir
			}

			d, err := LoadDoc(path)

			if err == md.ErrPathNotInitialized {
				return nil
			}
			if err != nil {
				return handler(path, nil, err)
			}
			return handler(path, d, nil)
		}
		return nil
	})
}
