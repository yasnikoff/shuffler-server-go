package shuffler_test

import (
	"testing"
	"io/ioutil"
	"reflect"
	"bitbucket.org/yasnikoff/shuffler/tests/dependency_fixtures"
	"os"
	"path/filepath"
	"bitbucket.org/yasnikoff/shuffler"
)

func TestAddRemove(t *testing.T) {
	t.Parallel()
	f := dependency_fixtures.New(t)
	defer f.Cleanup()

	p := f.User.Project
	//d :=f.Dep.Project
	depPath := "some path"

	// Add
	desc := f.AddDefaultDep(depPath)

	if len(p.Deps.Items()) != 1 {
		t.Errorf("Dependencies().Add() didn't cahnge dependencies count.")
	}

	items := p.Deps.Items()
	expectedItems := map[string]shuffler.Dependency{depPath:desc}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("After Dependencies().Add() Dependencies().Items() are: %q, want %q", items, expectedItems)
	}

	// Remove
	p.Deps.Remove(depPath)

	if len(p.Deps.Items()) != 0 {
		t.Fatalf("Dependencies().Remove() didn't cahnge dependencies count.")
	}
}

func TestProject(t *testing.T) {
	p := shuffler.NewProject()
	result := p.Deps.Project()
	if result != p {
		t.Errorf("Project.Dependencies().Project()==%s.\nWant %s", result, p)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	//t.Skip()
	f := dependency_fixtures.New(t)
	defer f.Cleanup()

	content := "text text text"
	filename := "file.txt"
	f.Dep.NewFile(filename, []byte(content))

	depPath := "some path"
	f.AddDefaultDep(depPath)
	//f.User.Project.Dependencies().Add(depPath, &dep.Descriptor{ProjectID:f.Dep.Project.ID()})

	if err := f.User.Project.Deps.UpdateAll(); err != nil {
		t.Fatalf("Update() returned error: %s", err)
	}
	targetPath := filepath.Join(f.User.Project.Path(), depPath, filename)
	data, err := ioutil.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("Cant read dep's file after Update(): %s", err)
	}
	readContent := string(data)
	if readContent != content {
		t.Errorf("Update() hasn't properly copied project's file.\nRead file content: %q\nWant: %q", readContent, content)
	}

}

func TestSaveLoad(t *testing.T) {
	t.Parallel()
	f := dependency_fixtures.New(t)
	defer f.Cleanup()

	loaded, err := shuffler.LoadProject(f.User.Project.Path())
	if err != nil {
		t.Fatal(err)
	}

	expDep := f.User.Project.Deps.Items()
	resDep := loaded.Deps.Items()
	if !reflect.DeepEqual(expDep, resDep) {
		t.Errorf("Loaded project have dependencies:\n%s\nWant:\n%s", resDep, expDep)
	}

}

func TestUpdate_EmptyProject(t *testing.T) {
	t.Parallel()
	f := dependency_fixtures.New(t)
	defer f.Cleanup()

	if err := f.User.Project.Deps.UpdateAll(); err != nil {
		t.Errorf("error calling Dependencies.Update() on a project without dependencies:\n%s", err)
	}
}

func TestCleanup(t *testing.T) {
	//t.Parallel()
	f := dependency_fixtures.New(t)
	defer f.Cleanup()

	depPath := "depPath"
	depPathAbs := f.User.Project.Deps.Abs(depPath)
	d := f.AddDefaultDep(depPath)

	// add file, so dep has something to copy. Will not create empty dir
	filename := f.Dep.NewFile(depPath, make([]byte, 0))
	filename, err := filepath.Rel(f.Dep.Path(), filename)
	filePathAbs:=filepath.Join(depPathAbs, filename)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.Dep.Project.SetupPath(depPathAbs, d.Parameters); err != nil {
		t.Fatal(err)
	}

	// check if dir is here
	if _, err := os.Stat(depPathAbs); err != nil {
		t.Fatal(err)
	}

	// check if the file is here
	if _, err := os.Stat(filePathAbs); err != nil {
		t.Errorf("dependency's Setup() method didn't create dep's file %s", filename)
	}

	if err := f.User.Project.Deps.CleanUp(depPath); err != nil {
		t.Fatal(err)
	}

	// check if all the dep's files are deleted
	if _, err := os.Stat(depPathAbs); err == nil {
		t.Errorf("dependency's Cleanup() method didn't remove dependency path.")
	}

	if _, err := os.Stat(filePathAbs); err == nil {
		t.Errorf("dependency's Setup() method didn't remove dep's file %s", filename)
	}

}
