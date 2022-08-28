package app

import (
	"bitbucket.org/yasnikoff/shuffler/tests/project_fixtures"
)

// TODO: tests

type DepData struct {
	subPath string
	project_fixtures.ProjectFixture
}
/*

type DepFixture struct{
	AppFixture
	Deps []DepData
}

func (f *DepFixture)Setup()error{
	if err:=f.AppFixture.Setup();err!=nil{
		return err
	}
	for _, d:=range f.Deps{
		for _, fd :=range d.fi
	}
}

func TestDep_Remove(t *testing.T) {


	f:=&DepFixture{Deps:[]DepData{
		{"dep1", project_fixtures.New()}
	}}
	if err:=f.Setup();err!=nil{
		t.Fatal(err)
	}
	defer f.Cleanup()
}
*/

