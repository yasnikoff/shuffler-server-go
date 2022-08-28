package shuffler

import (
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"fmt"
)

type ProjectMatcher func(*Project) bool

var MatchAll = func(pr *Project) bool {
	return true
}

type Repo interface {
	Name() string
	SetName(string)
	GetByID(ID id.ID) (*Project, error)
	GetManyByID(ids []id.ID) ([]*Project, error)
	//GetByQuery(ProjectQuery)([]Project, error)

	// absolute path to the highest directory that the repo has access
	Root() string

	// Path that all paths visible to users are relative to.
	// The path itself is relative to the Root
	Lib() string

	// Convert between absolute and library-relative path.
	ToLibPath(string) (string, error)
	FromLibPath(string) (string, error)

	UpdateIndexForPath(path string) error
	UpdateIndex() error
	AddToIndex(pr *Project) error

	//MoveProject(id id.ID, targetLibPath string)error

	AddToRecent(id.ID) error
	Recent() []id.ID
	Projects(ProjectMatcher) []id.ID

	SaveQuery(name, query string) error
	SavedQueries() map[string]string // TODO: change to map[string]ProjectQuery
}

type RepoClient interface {
	Repo() Repo
	SetRepo(Repo)
}

type repoClient struct {
	repo Repo
}

func (cl *repoClient)SetRepo(r Repo) {
	cl.repo = r
}

func (cl *repoClient)Repo() Repo {
	if cl == nil {
		return nil
	}
	return cl.repo
}

func (cl *repoClient)GetById(ID id.ID) (*Project, error) {
	if cl == nil {
		return nil, fmt.Errorf("can not get repo on nil client.")
	}
	if cl.repo == nil {
		return nil, fmt.Errorf("repo does not set for %s.", cl)
	}
	return cl.repo.GetByID(ID)
}

