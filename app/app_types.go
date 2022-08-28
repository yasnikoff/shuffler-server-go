package app


import "bitbucket.org/yasnikoff/shuffler/parameters"

func (app *app)Types()parameters.TypesRegistry{
	return app.typeReg
}
