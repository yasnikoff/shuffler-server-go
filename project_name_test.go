package shuffler

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/tests"
	"encoding/json"
)

func TestNameSet(t *testing.T) {
	pr := NewProject()
	name := "test project"
	pr.SetName(name)
	if pr.Name() != name {
		t.Errorf("Project.SetName(%s) didn't set the name. Project.Name()==%q, want %[1]q", name, pr.Name())
	}
}

func TestProjectNameJSONRoundtrip(t *testing.T) {

	pr:=NewProject()
	name:="test name"
	pr.SetName(name)
	data,err:=json.Marshal(pr)
	if err!=nil{
		t.Fatal(err)
	}
	pr2:=NewProject()
	if err:=json.Unmarshal(data, pr2);err!=nil{
		t.Fatal(err)
	}
	if pr2.Name()!=name{
		t.Errorf("Unmarshalled project name is %q, want it be equal to initial value %q.", pr2.Name(), name)
	}
}

func TestNameSaveLoad(t *testing.T) {
	f := tests.NewDirFx(t)
	defer f.Cleanup()

	pr := NewProject()
	name := "test project"

	pr.SetName(name)
	if pr.Name()!=name{
		t.Fatalf("Project.SetName(%s) didn't set the name. Project.Name()==%q, want %[1]q", name, pr.Name())

	}
	pr.SaveTo(f.Path())

	pr, err:=LoadProject(pr.Path())
	if err!=nil{
		t.Fatal(err)
	}
	if pr.Name() != name {
		t.Errorf("Loaded project has different name from the saved one. Project.Name()==%q, want %[1]q", name, pr.Name())
	}
}
func TestNameNew(t *testing.T) {

	pr:=NewProject()
	if pr.Name()!=""{
		t.Errorf("New not save project without name set must have empty name. pr.Name()==%q.",pr.Name())
	}

}


