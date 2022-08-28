package app

import (
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"path/filepath"
	"github.com/skratchdot/open-golang/open"
	"bitbucket.org/yasnikoff/shuffler"
	"encoding/json"
	"bitbucket.org/yasnikoff/shuffler/fs"
	"os"
	"github.com/pkg/errors"
	"strings"
)

func (app *app)NewProject(path string) (*shuffler.Project, error) {

	p := shuffler.NewProject()
	p.SetRepo(app.Repo())

	name := filepath.Base(path)
	if name == "." {
		name = ""
	}
	p.SetName(name)
	p.Signature().SetTypesRegistry(app.Types())

	if err := p.SaveTo(path); err != nil {
		return nil, err
	}
	if err := app.Repo().UpdateIndexForPath(p.Path()); err != nil {
		return nil, err
	}
	app.repo.AddToRecent(p.ID())
	app.Logger().Printf("Initialized %q with %s", p.Path(), p)
	return p, nil
}

func (app *app)GetProjectByIDString(ID string) (*shuffler.Project, error) {
	prID, err := id.Parse(ID)
	if err != nil {
		return nil, err
	}
	return app.GetProject(prID)
}

func (app *app)GetProject(ID id.ID) (*shuffler.Project, error) {
	pr, err := app.Repo().GetByID(ID)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (app *app)LoadProject(path string) (*shuffler.Project, error) {
	pr, err := shuffler.LoadProject(path)
	if err != nil {
		return nil, err
	}
	pr.SetRepo(app.Repo())
	pr.Signature().SetTypesRegistry(app.Types())
	return pr, nil
}

func (app *app)OpenProjectFolder(PrID string) error {
	// TODO: allow only for local clients
	ID, err := id.Parse(PrID)
	if err != nil {
		app.Logger().Println(err)
		return err
	}
	doc, err := app.Repo().GetByID(ID)

	if err != nil {
		app.Logger().Println(err)
	}
	pr := doc
	return open.Start(pr.Path())
}

func (app *app)CurrentProject() (*shuffler.Project, error) {
	return app.LoadProject(".") // TODO: move CurrentProject to cmd interface. It has little to no sense in gui mode
}

func (app *app)DuplicateProject(ID string, newName string, newPath string) (*shuffler.Project, error) {

	var err error
	pr, err := app.GetProjectByIDString(ID)
	if err != nil {
		return nil, err
	}

	// TODO: test for spaces in name and path
	newPath=strings.TrimSpace(newPath)
	newName=strings.TrimSpace(newName)

	// first, try to copy files
	var filesToCopy []string

	oldPath := pr.Path()
	filesToCopy, err = fs.MatchFilesFromPath(pr.Path(), shuffler.ProjectDuplicateFilesMatcher)
	if err != nil {
		return nil, err
	}

	for _, file := range filesToCopy {
		fs.Copy(filepath.Join(newPath, file), filepath.Join(oldPath, file))
	}


	// create new project and copy data from the origin

	data, err := json.Marshal(pr)
	if err != nil {
		return nil, err
	}

	dup, ok := pr.Type().New().(*shuffler.Project)
	if !ok {
		return nil, &shuffler.WrongDocTypeError{dup, "Project"}
	}

	newId := dup.ID()

	// fill in all the data from origin
	if err := json.Unmarshal(data, dup); err != nil {
		return nil, err
	}

	// fill in changed data
	dup.SetID(newId)
	dup.SetPath(newPath)
	dup.SetName(newName)
	dup.Origin = pr.ID()


	// save
	if err := dup.Save(); err != nil {
		return nil, err
	}

	// update repo's index

	if err := app.Repo().UpdateIndexForPath(dup.Path()); err != nil {
		return dup, err
	}

	// add to recent
	//if err:=app.Repo().AddToRecent(newId);err!=nil{
	//	return nil, errors.Wrap("can't add duplicated project to recent", err)
	//}
	return dup, nil
}
func (app *app)MoveProject(ID id.ID, newPath string) (ok bool, err error) {

	if fs.PathExists(newPath) {
		files, err := fs.FilesFromPath(newPath)
		if err != nil {
			return false, errors.Wrap(err, "can't get files list from target path")
		}
		if len(files) > 0 {
			return false, errors.New("target dir is not empty")
		}
	}

	pr, err := app.GetProject(ID)
	if err != nil {
		return false, err
	}

	// first, try to copy files
	var filesToCopy []string

	oldPath := pr.Path()
	filesToCopy, err = fs.MatchFilesFromPath(pr.Path(), shuffler.ProjectMoveFileMatcher)
	if err != nil {
		return false, err
	}

	for _, file := range filesToCopy {
		if err := fs.Copy(filepath.Join(newPath, file), filepath.Join(oldPath, file)); err != nil {
			return false, err
		}
	}

	// remove old path

	if err := os.RemoveAll(oldPath); err != nil {
		return true, err
	}

	// update repo's index (only after deleting the old copy. Otherwise, "same ID is also found at..." occurs)

	if err := app.Repo().UpdateIndexForPath(newPath); err != nil {
		return true, err
	}
	return true, nil
}

