package dependency_fixtures

import (
	"bitbucket.org/yasnikoff/shuffler/tests/project_fixtures"
	"bitbucket.org/yasnikoff/shuffler/repo/repo_fixtures"
	"bitbucket.org/yasnikoff/shuffler/tests"
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"testing"
	"bitbucket.org/yasnikoff/shuffler"
)

// Fixture with one main project User and one default dependency project Dep.
// Additional dependency projects can be added by calling AddDep.
// default dependency project is not added to User as dependency during fixture creation.
// To add default dependency project as dependency to User, use AddDefaultDep.
type DepFixture struct {
	repo_fixtures.RepoFixture
	Dep  *project_fixtures.ProjectFixture
	User *project_fixtures.ProjectFixture
}

func New(t *testing.T) *DepFixture {
	return NewFromDirFx(tests.NewDirFx(t))
}

func NewFromDirFx(dirFx *tests.DirFx) *DepFixture {
	f := &DepFixture{RepoFixture:*repo_fixtures.NewFromDirFx(dirFx)}

	f.Dep = project_fixtures.NewFromDirFx(f.RepoFixture.NewDirFx("dep"))
	f.User = project_fixtures.NewFromDirFx(f.RepoFixture.NewDirFx("user"))

	f.User.Project.SetRepo(f.Repo)

	if err := f.Repo.UpdateIndex(); err != nil {
		f.Fatal(err)
	}

	return f
}

func (f *DepFixture)AddDep(target *project_fixtures.ProjectFixture, path string, dep *project_fixtures.ProjectFixture, params *parameters.Collection) *shuffler.DepDesc {
	desc := target.AddDep(path, dep, params)
	f.UpdateIndex()
	return desc
}

func (f *DepFixture)AddDefaultDep(path string) *shuffler.DepDesc {
	return f.AddDep(f.User, path, f.Dep, nil)
}

