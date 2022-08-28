package repo_fixtures

import (
	"bitbucket.org/yasnikoff/shuffler"
	"bitbucket.org/yasnikoff/shuffler/tests"
	"bitbucket.org/yasnikoff/shuffler/repo"
	"fmt"
	"bitbucket.org/yasnikoff/shuffler/tests/project_fixtures"
	"testing"
)

type RepoFixture struct {
	tests.DirFx
	Projects []*shuffler.Project
	Repo     *repo.Repo
}

func New(t *testing.T) *RepoFixture {
	root := tests.NewDirFx(t)
	return NewFromDirFx(root)
}

func NewFromDirFx(d *tests.DirFx) *RepoFixture {
	r := &RepoFixture{
		DirFx:*d,
		Repo:NewRepo(d.T, d.Path()),
	}
	if err := r.Repo.UpdateIndex(); err != nil {
		r.T.Fatal(err)
	}
	return r
}

func NewDefault(t *testing.T) *RepoFixture {
	r := New(t)
	for i := 0; i < 3; i++ {
		r.NewProject()
	}
	return r
}

// Add project to the fixture.
func (rf *RepoFixture)NewProject() *shuffler.Project {
	return rf.NewProjectWithPath("")

}

func (rf *RepoFixture)NewProjectFx(path string) *project_fixtures.ProjectFixture {

	var dirFx *tests.DirFx
	if path == "" {
		dirFx = rf.NewRandDirFx()
	} else {

		dirFx = rf.NewDirFx(path)
	}
	return project_fixtures.NewFromDirFx(dirFx)
}

func (rf *RepoFixture)NewProjectWithPath(path string) *shuffler.Project {
	prFx := rf.NewProjectFx(path)
	pr := prFx.Project
	rf.Projects = append(rf.Projects, pr)
	if err := rf.Repo.UpdateIndex(); err != nil {
		rf.T.Fatal(err)
	}
	return pr

}

func NewRepo(t *testing.T, path string) *repo.Repo {
	r := repo.New()
	if err := r.SaveTo(path); err != nil {
		t.Fatal(err)
	}
	return r
}

func (rf *RepoFixture)UpdateIndex() {
	if err := rf.Repo.UpdateIndex(); err != nil {
		rf.Fatal(err)
	}
}

func (rf *RepoFixture)AddToIndex(prFx *project_fixtures.ProjectFixture){
	if err:=rf.Repo.AddToIndex(prFx.Project);err!=nil{
		rf.Fatal(err)
	}
	//check
	_,err:=rf.Repo.GetByID(prFx.Project.ID())
	if err!=nil{
		rf.Fatalf("Didn't add project %s to repo fixture.", prFx.Project)
	}
}

func checkProject(pr *shuffler.Project) error {
	loadedPr, err := shuffler.LoadProject(pr.Path());
	if err != nil {
		return fmt.Errorf("Error with project fixture. Can't load newly setup project from %q:\n%s", pr.Path(), err)
	}
	if loadedPr == nil {
		return fmt.Errorf("Error with project fixture. Loaded nil from newly setup project's path %q", pr.Path())
	}
	if loadedPr.ID() != pr.ID() {
		return fmt.Errorf("Error with project fixture. Project loaded  from newly setup project's path %q has ID=%q, want %q", pr.Path(), loadedPr.ID(), pr.ID())
	}
	return nil
}
