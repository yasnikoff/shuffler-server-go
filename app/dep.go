package app

import (
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"path/filepath"
	"fmt"
	"strings"
	"bitbucket.org/yasnikoff/shuffler"
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"github.com/pkg/errors"
)

func (app *app)AddDep(projectID, path, dependencyID string) (err error) {

	defer func() {
		if err != nil {
			err = errors.Wrapf(err, "can't add dep %s at %s to project %s",dependencyID, path, projectID)
		}

	}()

	prID, err := id.Parse(projectID)
	if err != nil {
		return err
	}

	depID, err := id.Parse(dependencyID)
	if err != nil {
		return err
	}

	//TODO: move checks to Dependencies.Add()


	if depID == prID {
		return errors.New("can't add project as a dependency to itself.")
		// TODO: check for cyclic deps
	}

	pr, err := app.Repo().GetByID(prID)
	if err != nil {
		return err
	}

	if _, err := pr.Deps.Get(path); err == nil {
		return errors.New("there is already a dependency at the path")
	}

	if filepath.IsAbs(path) {
		return fmt.Errorf("absolute paths are not allowed for dependencies.")
	}

	absPath, err := filepath.Abs(filepath.Join(pr.Path(), path))
	if err != nil {
		return errors.Wrap(err, "can't build absolute path")
	}
	if !strings.HasPrefix(absPath, pr.Path()) {
		return errors.New("path is not the project's subpath")
	}
	relPath, err := filepath.Rel(pr.Path(), absPath)
	if err != nil {
		return errors.New("can't create path relative to project's path")
	}

	_, err = app.repo.GetByID(depID)
	if err != nil {
		return errors.Wrapf(err, "can't add dependency %s to %s", depID, prID)
	}

	dp := shuffler.DepDesc{IDData:depID}
	subPath := filepath.ToSlash(relPath)
	err = pr.Deps.Add(subPath, &dp)
	if err != nil {
		return err
	}
	if err := pr.Save(); err != nil {
		return err
	}

	if err := pr.Deps.Setup(subPath); err != nil {
		return err
	}
	return nil
}

func (app *app)RemoveDep(projectID, depPath string) error {
	pr, err := app.GetProjectByIDString(projectID)
	if err != nil {
		return err
	}

	if err := pr.Deps.CleanUp(depPath); err != nil {
		return err
	}
	pr.Deps.Remove(depPath)

	return pr.Save()
}

func (app *app)UpdateDep(projectID, depPath string, recursive bool) error {
	// TODO: remove "recursive" parameter
	pr, err := app.GetProjectByIDString(projectID)
	if err != nil {
		return err
	}

	if err := pr.Deps.Update(depPath); err != nil {
		return err
	}
	return nil
}

func (app *app)UpdateDepAll(projectID string, recursive bool) error {
	// TODO: remove "recursive" parameter

	pr, err := app.GetProjectByIDString(projectID)
	if err != nil {
		return err
	}

	if err := pr.Deps.UpdateAll(); err != nil {
		return err
	}
	return nil
}

func (app *app)ReplaceDep(projectID, subPath, newDepID string) error {
	prID, err := id.Parse(projectID)
	if err != nil {
		return err
	}

	pr, err := app.Repo().GetByID(prID)
	if err != nil {
		return err
	}

	newProtoID, err := id.Parse(newDepID)
	if err != nil {
		return err
	}

	dep, err := pr.Deps.Get(subPath)
	if err != nil {
		return err
	}

	if dep.Descriptor().ID() == newProtoID {
		return nil
	}
	if err := pr.Deps.CleanUp(subPath); err != nil {
		return err
	}
	pr.Deps.Remove(subPath)
	pr.Deps.Add(subPath, &shuffler.DepDesc{IDData:newProtoID, Parameters:parameters.NewCollection(nil)})
	if err := pr.Deps.Setup(subPath); err != nil {
		return err
	}

	return pr.Save()
}
