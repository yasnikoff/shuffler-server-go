package app_fixtures

import (
	"bitbucket.org/yasnikoff/shuffler/app"
	"bitbucket.org/yasnikoff/shuffler/repo/repo_fixtures"
	"testing"
)

type AppFixture struct {
	*repo_fixtures.RepoFixture
	App app.App
}

func New(t *testing.T) *AppFixture {
	f := &AppFixture{RepoFixture:repo_fixtures.New(t)}
	f.App = app.New()
	f.App.SetRepo(f.Repo)
	return f
}

