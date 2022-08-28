package shuffler_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/tests/project_fixtures"
	"bitbucket.org/yasnikoff/shuffler"
)

func TestAfterSave(t *testing.T) {
	t.Parallel()
	var called bool
	shuffler.SetAfterSave(func(pr *shuffler.Project)error{
		called=true
		return nil
	})
	f:=project_fixtures.New(t)
	defer f.Cleanup()

	if err:=f.Project.Save();err!=nil{
		t.Fatal(err)
	}

	if !called{
		t.Errorf("Project.afterSave() was not called when saving project")
	}
}
