package repo_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/repo/repo_fixtures"
	"github.com/termie/go-shutil"
	"os"
	"bitbucket.org/yasnikoff/shuffler"
)

type ProjectMoveFixture struct {
	repo_fixtures.RepoFixture
	pr      *shuffler.Project
	oldPath string
	newPath string
}

func NewProjectMoveFixture(t *testing.T) *ProjectMoveFixture {
	r := repo_fixtures.New(t)
	pr := r.NewProjectWithPath("old")
	f := &ProjectMoveFixture{
		RepoFixture:*r,
		pr:pr,
		oldPath : r.Abs("old"),
		newPath : r.Abs("new"),
	}
	return f
}

func (f *ProjectMoveFixture)CloneProject() {
	if err := shutil.CopyTree(f.oldPath, f.newPath, nil); err != nil {
		f.T.Fatalf("Error cloning project: %q", err)
	}
}

func (f *ProjectMoveFixture)RemoveOldProjectCopy() {
	if err := os.RemoveAll(f.oldPath); err != nil {
		f.T.Fatalf("Error removing old project copy: %q", err)
	}

}

func (f *ProjectMoveFixture)Project() *shuffler.Project {
	return f.pr
}

func TestIndex(t *testing.T) {
	t.Parallel()
	f := NewProjectMoveFixture(t)
	defer f.Cleanup()

	pr := f.Project()
	// before moving project
	if f.Repo.Index.Len() != 1 {
		t.Fatalf("new repo has %d items in the index. Want 1", f.Repo.Index.Len())
	}

	if _, err := f.Repo.GetByID(pr.ID()); err != nil {
		t.Fatal(err)
	}

	f.CloneProject()

	// now there are two copies of the project
	if err := f.Repo.UpdateIndex(); err == nil {
		t.Fatalf("While indexing path with two copies of the project repo.UpdateIdex() didn't return an error.")
	}

	f.RemoveOldProjectCopy()

	if err := f.Repo.UpdateIndex(); err != nil {
		t.Fatalf("After moving project repo.UpdateIdex() returned error:\n%s", err)
	}

	if f.Repo.Index.Len() != 1 {
		t.Fatalf("after moving project repo has %d items in the index. Want 1. ", f.Repo.Index.Len(), f.Repo.Index)
	}

}
