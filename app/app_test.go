package app_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/app/app_fixtures"
)


func TestProjectInit(t *testing.T) {
	f := app_fixtures.New(t)
	defer f.Cleanup()

	projectPath := "new project"
	pr, err := f.App.NewProject(f.Abs(projectPath))
	if err != nil {
		t.Fatal(err)
	}
	if pr == nil {
		t.Fatal("App.NewProject() ended without error but returned nil project")
	}
}
