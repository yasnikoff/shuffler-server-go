package repo

import (
	"path/filepath"
	"github.com/pkg/errors"
)

func (r *Repo)ToLibPath(path string) (string, error) {
	result, err:=filepath.Rel(r.libPath(), path)
	if err!=nil{
		return "", err
	}
	return filepath.ToSlash(result), nil
}

func (r *Repo)FromLibPath(path string) (string, error) {
	result:= filepath.Clean(filepath.Join(r.libPath(), path))
	if !filepath.HasPrefix(result, r.Path()){
		return "", errors.New("resulting path is not under the repo's Root")
	}
	return result, nil
}

