package shuffler_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/tests"
	"bitbucket.org/yasnikoff/shuffler/fs"
	"bitbucket.org/yasnikoff/shuffler/repo/repo_fixtures"

	"reflect"
	"bitbucket.org/yasnikoff/shuffler/tests/project_fixtures"
	"strings"
	"sort"
)

type SetupFx struct {
	tests.DirFx
	repoFx    *repo_fixtures.RepoFixture
	user      *project_fixtures.ProjectFixture
	targetDir *tests.DirFx
	caseName  string
}

func NewSetupFx(t *testing.T) *SetupFx {
	f := &SetupFx{DirFx:*tests.NewDirFx(t)}
	f.repoFx = repo_fixtures.NewFromDirFx(f.NewDirFx("repo"))
	f.user = f.AddProject("user")
	f.targetDir = f.NewDirFx("target")
	return f
}

func (f *SetupFx)AddProject(subPath string, files ...string) *project_fixtures.ProjectFixture {
	prFx := project_fixtures.NewFromDirFx(f.NewDirFx(subPath))
	for _, file := range files {
		prFx.NewRandFileFx(file)
	}
	prFx.Project.SetRepo(f.repoFx.Repo)
	f.repoFx.AddToIndex(prFx)

	return prFx
}

func (f *SetupFx)Expect(files ...string) {

	if err := f.user.Project.SetupPath(f.targetDir.Path(), nil); err != nil {
		f.Fatal(err)
	}

	if f.caseName == "" {
		f.caseName = "unspecified"
	}

	result, err := fs.FilesFromPath(f.targetDir.Path())
	if err != nil {
		f.Fatal(err)
	}
	sort.Strings(result)

	expected := files
	if expected == nil {
		expected = []string{}
	}
	sort.Strings(expected)

	if !reflect.DeepEqual(result, expected) {
		msg := "case " + strings.ToUpper(f.caseName) + ": project setup created these files at the target directory:\n     %q\nwant %q"
		f.Errorf(msg, result, expected)
	}
}

func TestEmptySetup(t *testing.T) {
	f := NewSetupFx(t)
	defer f.Cleanup()

	f.caseName = "project with no dependencies"

	f.Expect()
}

func TestWithNoDepsSetup(t *testing.T) {
	f := NewSetupFx(t)
	defer f.Cleanup()

	f.caseName = "project with no dependencies"

	f.user.NewRandFile("userfile")

	f.Expect("userfile")
}

func TestWithOneDepSetup(t *testing.T) {
	f := NewSetupFx(t)
	defer f.Cleanup()

	f.caseName = "project with one dependency"

	f.user.NewRandFile("userfile")

	dep := f.AddProject("dep", "depfile")
	f.user.AddDep("dep", dep, nil)

	f.Expect("userfile", "dep/depfile")

}


func TestWithDepsSetup(t *testing.T) {
	f := NewSetupFx(t)
	defer f.Cleanup()

	f.caseName = "project with 2-nd order dependency"

	f.user.NewRandFile("userfile")

	dep := f.AddProject("dep", "depfile")
	f.user.AddDep("dep", dep, nil)

	depDep := f.AddProject("depDep", "depdepfile")
	dep.AddDep("dep", depDep, nil)

	f.Expect("userfile", "dep/depfile", "dep/dep/depdepfile")
}

func TestWithDepsWithPrivateUnderscoreFilesSetup(t *testing.T) {
	f := NewSetupFx(t)
	defer f.Cleanup()

	f.caseName = "project with 2-nd order dependency with _private file"

	f.user.NewRandFile("userfile")

	dep := f.AddProject("dep", "depfile")
	f.user.AddDep("dep", dep, nil)

	depDep := f.AddProject("depDep", "depdepfile", "_private")
	dep.AddDep("dep", depDep, nil)

	f.Expect("userfile", "dep/depfile", "dep/dep/depdepfile")
}

func TestWithDepsWithPrivateDotFilesSetup(t *testing.T) {
	f := NewSetupFx(t)
	defer f.Cleanup()

	f.caseName = "project with 2-nd order dependency with .private file"

	f.user.NewRandFile("userfile")

	dep := f.AddProject("dep", "depfile")
	f.user.AddDep("dep", dep, nil)

	depDep := f.AddProject("depDep", "depdepfile", ".private")
	dep.AddDep("dep", depDep, nil)

	f.Expect("userfile", "dep/depfile", "dep/dep/depdepfile")
}


func TestWithDepsWithPrivateDepPathSetup(t *testing.T) {
	f := NewSetupFx(t)
	defer f.Cleanup()

	f.caseName = "project with dependency that has _private path"

	f.user.NewRandFile("userfile")

	dep := f.AddProject("dep", "depfile")
	f.user.AddDep("_dep", dep, nil)

	f.Expect("userfile")
}
