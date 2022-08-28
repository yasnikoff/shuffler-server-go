package app

import (
	"log"
	"os"
	"bitbucket.org/yasnikoff/shuffler/logging"
	"bitbucket.org/yasnikoff/shuffler"
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"bitbucket.org/yasnikoff/shuffler/repo"
	"bitbucket.org/yasnikoff/shuffler/parameters"
)

type App interface {
	//Start()error
	Fatal(error)
	Logger() *log.Logger
	AddDep(projectID, path, dependencyID string) error
	RemoveDep(projectID, depPath string) error
	UpdateDep(projectID, depPath string, recursive bool) error
	UpdateDepAll(projectID string, recursive bool) error
	ReplaceDep(projectID, subPath, newDepID string) error
	NewProject(path string) (*shuffler.Project, error)
	GetProjectByIDString(ID string) (*shuffler.Project, error)
	GetProject(ID id.ID) (*shuffler.Project, error)
	LoadProject(path string) (*shuffler.Project, error)
	OpenProjectFolder(PrID string) error
	CurrentProject() (*shuffler.Project, error)
	//TODO: change ID type from string to ID
	DuplicateProject(ID string, newName string, newPath string) (*shuffler.Project, error)
	MoveProject(ID id.ID, newPath string) (ok bool, err error)
	Repo() shuffler.Repo
	SetRepo(r shuffler.Repo)
	SaveQuery(name, query string) error
	AddTags(projectID id.ID, tags ...string) error
	Types() parameters.TypesRegistry
}

func New() App {
	a := &app{typeReg:parameters.NewTypesRegistry()}
	a.SetLogger(log.New(os.Stdout, "", log.LstdFlags))
	return a
}

// Bootstrap loads repo from given path and starts app with the repo.
func Bootstrap(path string) (App, error) {
	r, err := repo.LoadRepo(path)
	if err != nil {
		return nil, err
	}
	a := New()
	a.SetRepo(r)
	return a, nil
}

type app struct {
	typeReg parameters.TypesRegistry
	repo    shuffler.Repo
	logging.Log

}

func (app *app)Fatal(err error) {
	app.Logger().Fatalf("App is shutting down because of the fatal erorr: %s", err)
	os.Exit(1)
}
