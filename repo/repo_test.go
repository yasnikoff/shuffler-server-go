package repo_test

import (
	"testing"
	"os"
	"path/filepath"
	"bitbucket.org/yasnikoff/shuffler/repo"
	"bitbucket.org/yasnikoff/shuffler/tests"
	"bitbucket.org/yasnikoff/shuffler/repo/repo_fixtures"
)

func TestDefaultRepo(t *testing.T) {
	f := tests.NewDirFx(t)
	defer f.Cleanup()

	f2 :=tests.NewDirFx(t)
	defer f2.Cleanup()

	subdirs := []string{
		"1/2/3",
		"1",
		"level 1/level2",
		"",
	}
	setup := func(path string) []string {
		paths := make([]string, len(subdirs))
		for i, sd := range subdirs {
			paths[i] = filepath.Join(path, sd)
			if err := os.MkdirAll(paths[i], os.ModeDir); err != nil {
				t.Fatal(err)
			}
		}
		return paths
	}

	paths_positive := setup(f.Path())
	paths_negative := setup(f2.Path())

	expected := repo.New()
	err := expected.SaveTo(f.Path())
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range paths_positive {
		defaultRepo, err := repo.Discover(path)
		if err != nil {
			t.Errorf("Error finding default repo for %q: %v", path, err)
			continue
		}
		if defaultRepo.ID() != expected.ID() {
			t.Errorf("result of repo.Discover(%q) is wrong. Got %v. Want %v", path, defaultRepo, expected)
		}
		if defaultRepo.Path() != f.Path() {
			t.Errorf("repo.Discover(%q).Path() is wrong.\nGot\t %q.\nWant\t %q", path, defaultRepo.Path(), f.Path())
		}
	}
	for _, path := range paths_negative {
		defaultRepo, err := repo.Discover(path)
		if err == nil {
			t.Errorf("Erroneously found default repo for %q: %v.", path, defaultRepo)
		}
	}
}

func TestUpdateIndex(t *testing.T) {
	f := repo_fixtures.NewDefault(t)
	defer f.Cleanup()

	funcUnderTest := "When called on normal path with projects, Repo.UpdateIndex() "

	if f.Repo.Index.Len() != len(f.Projects) {
		t.Errorf(funcUnderTest + "didn't populate Repo.Index with correct number of Projects. Got: %d, Want: %d", f.Repo.Index.Len(), len(f.Projects))
	}

	for _, pr := range f.Projects {
		if pr == nil {
			t.Errorf("repo_fixture.Projects[%d]==nil. It's a bug in the fixture setup.")
			continue
		}
		path,err:=filepath.Rel(f.Repo.Path(),pr.Path())
		if err!=nil{
			t.Fatal(err)
		}
		if pathGot, ok := f.Repo.Index.Get(pr.ID()); filepath.ToSlash(pathGot) != filepath.ToSlash(path) || !ok {
			t.Errorf(funcUnderTest + "populated repo's index with wrong path for id %q:\nGot %q\nWant %q\n", pr.ID(), pathGot, path)
			continue
		}
	}

}

func TestDocByID(t *testing.T) {
	f := repo_fixtures.NewDefault(t)
	defer f.Cleanup()

	for path, pr := range f.Projects {
		if pr == nil {
			t.Fatalf("repo_fixture.Projects[%s]==nil. It's bug in the repo fixture setup.", path)
		}
		result, err := f.Repo.GetByID(pr.ID())
		if err != nil {
			t.Errorf("Repo.DocByID(%s) called with an existing project's ID returned error:\n%s.\nWant it to return without errors.", pr.ID(), err)
			continue
		}
		if result == nil {
			t.Errorf("Repo.DocByID(%s) returned nil. Want %s", pr.ID(), pr)
			continue
		}

		if result.ID() != pr.ID() {
			t.Errorf("Repo.DocByID(%s).ID()=%s", pr.ID(), result.ID())
		}

		if result.Repo()!=f.Repo{
			t.Errorf("Repo.DocByID().Repo()==%s, want %s", result.Repo(), f.Repo)
		}
	}
}
