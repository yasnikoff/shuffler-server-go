package shuffler

import (
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"fmt"
	"bitbucket.org/yasnikoff/shuffler/fs"
	"bitbucket.org/yasnikoff/goutil/errutil/errorlist"
	"path/filepath"
	"sync"
)

func (pr *Project)Setup(params *parameters.Collection) (err error) {
	return pr.SetupPath(pr.Path(), params)
}

func (pr *Project)SetupPath(path string, params *parameters.Collection) (err error) {

	var repo Repo
	if len(pr.Deps.items) > 0 {
		repo = pr.Repo()
		if repo == nil {
			return fmt.Errorf("Project %s has dependencies but no repo", pr)
		}
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("Project.SetupPath() only accepts absolute paths. Got %q", path)
	}
	targetPath, err := fs.NormalizePath(path)
	if err != nil {
		return err
	}

	errors := make(chan error)
	errLst := errorlist.ErrorList{} // TODO: replace with errors channel?
	defer func() {
		if errLst.HasErrors() {
			err = errLst // FIXME: errors are not returned. need tests
		}
	}()

	go func() {
		for err := range errors {
			errLst.AddNotNil(err)
		}
	}()

	wg := sync.WaitGroup{}
	defer wg.Wait()

	fileMap, err := pr.FileMapWithDeps(params, "")
	if err != nil {
		return err
	}
	for dst, src := range fileMap {
		wg.Add(1)
		go func(dst, src string) {
			defer wg.Done()
			if err := fs.Copy(dst, src); err != nil {
				errors <- err
			}
		}(filepath.Join(targetPath, dst),  src)
	}

	close(errors)

	return nil
}

