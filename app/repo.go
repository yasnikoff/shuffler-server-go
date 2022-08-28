package app

import (
	"bitbucket.org/yasnikoff/shuffler/repo"
	"time"
	"fmt"
	"bitbucket.org/yasnikoff/shuffler"
)

func (app *app)Repo() shuffler.Repo {
	if app.repo == nil {
		r, err := repo.Discover(".")
		if err != nil {
			app.Fatal(fmt.Errorf("Can't discover repo"))
		}
		app.SetRepo(r)
	}

	// TODO: remove when implementing not-in-memory repo index
	r, ok := app.repo.(*repo.Repo)
	if ok {

		if !r.IsUpToDate(2 * time.Second) {
			// TODO: benchmarks
			err := r.Reload()
			if err != nil {
				app.Logger().Printf("can't load repo from %q:\n%s", r.Path(), err)
				return app.repo
			}
		}
	}
	return app.repo
}

func (app *app)SetRepo(r shuffler.Repo) {
	//r.SetLogger(app.Logger())
	app.repo = r
}

func (app *app)SaveQuery(name, query string) error {
	r := app.Repo()
	// TODO; sanity checks
	return r.SaveQuery(name, query)
}
