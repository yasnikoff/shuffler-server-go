package project_fixtures

import (
	"bitbucket.org/yasnikoff/shuffler"
	"bitbucket.org/yasnikoff/shuffler/tests"
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"testing"
)

type ProjectFixture struct {
	tests.DirFx
	Project *shuffler.Project
}

func New(t *testing.T) *ProjectFixture {
	root := tests.NewDirFx(t)
	return &ProjectFixture{DirFx:*root, Project:NewProject(t, root.Path())}
}

func NewFromDirFx(f *tests.DirFx) *ProjectFixture {
	return &ProjectFixture{DirFx:*f, Project:NewProject(f.T, f.Path())}
}
func NewProject(t *testing.T, path string) *shuffler.Project {
	pr := shuffler.NewProject()
	if err := pr.SaveTo(path); err != nil {
		t.Fatal(err)
	}
	return pr
}

func (f *ProjectFixture)Save() {
	if err := f.Project.Save(); err != nil {
		f.T.Fatal(err)
	}
}

func (f *ProjectFixture)AddDep(path string, dp *ProjectFixture, params *parameters.Collection) *shuffler.DepDesc {
	if params == nil {
		params = dp.Project.DefaultParams()
	}

	desc := &shuffler.DepDesc{IDData:dp.Project.ID(), Parameters:params}
	if err := f.Project.Deps.Add(path, desc); err != nil {
		f.Fatal(err)
	}
	if err := f.Project.Save(); err != nil {
		f.Fatal(err)
	}
	return desc
}

func (f *ProjectFixture)UpdateAllDeps() {
	if err := f.Project.Deps.UpdateAll(); err != nil {
		f.Fatal(err)
	}
}
func (f *ProjectFixture)CleanUpAllDeps() {

	if err := f.Project.Deps.CleanUpAll(); err != nil {
		f.Fatal(err)
	}
}
