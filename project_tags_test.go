package shuffler_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/tests/project_fixtures"
	"bitbucket.org/yasnikoff/shuffler"
	"reflect"
)

func TestTags(t *testing.T) {
	f:=project_fixtures.New(t)
	defer f.Cleanup()

	pr:=f.Project
	tags:=[]string{"1", "2"}
	if len(pr.Tags.Get())>0{
		t.Errorf("New project has tags: %s, %s", pr, tags)
	}
	pr.Tags.Add(tags...)
	if !reflect.DeepEqual(pr.Tags.Get(),tags){
		t.Errorf("project Tags()==%q, want %q", pr.Tags.Get(), tags)
	}

	if err:=pr.Save();err!=nil{
		t.Fatal(err)
	}

	loadedPr,err:=shuffler.LoadProject(pr.Path())
	if err!=nil{
		t.Fatal(err)
	}
	if !reflect.DeepEqual(loadedPr.Tags.Get(), tags){
		t.Errorf("loaded project's Tags()==%q, want %q", loadedPr.Tags.Get(),tags)
	}
}
