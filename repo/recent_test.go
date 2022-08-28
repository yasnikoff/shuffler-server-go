package repo_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/repo/repo_fixtures"
	"reflect"
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"bitbucket.org/yasnikoff/shuffler/repo"
)

func TestNewRecent(t *testing.T) {

	f := repo_fixtures.New(t)
	defer f.Cleanup()

	recent := f.Repo.Recent()
	expected := []id.ID{}

	if !reflect.DeepEqual(recent, expected) {
		t.Errorf("Recent() method of a new repo reurned %s, want %s", recent, expected)
	}

}

func TestRepo_AddToRecent(t *testing.T) {

	f:= repo_fixtures.New(t)
	defer f.Cleanup()

	pr:=f.NewProject()

	f.Repo.AddToRecent(pr.ID())

	recent:=f.Repo.Recent()
	expected:=[]id.ID{pr.ID()}

	if !reflect.DeepEqual(recent, expected){
		t.Errorf("After calling Repo.AddToRecent(pr.ID()), repo.Recent() is %s, want %s", recent, expected)
	}

}
func TestRepo_AddToRecent_Load(t *testing.T) {

	f:= repo_fixtures.New(t)
	defer f.Cleanup()

	pr:=f.NewProject()

	f.Repo.AddToRecent(pr.ID())

	loadedRepo, err:=repo.LoadRepo(f.Repo.Path())
	if err!=nil{
		t.Fatal(err)
	}

	recent:= loadedRepo.Recent()
	expected:=[]id.ID{pr.ID()}

	if !reflect.DeepEqual(recent, expected){
		t.Errorf("After calling Repo.AddToRecent(pr.ID()), changes are not saved to disk. LoadedRepo.Recent() is %s, want %s", recent, expected)
	}

}
