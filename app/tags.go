package app

import "bitbucket.org/yasnikoff/shuffler/types/id"

func (app *app)AddTags(projectID id.ID, tags ...string) error {
	pr, err := app.GetProject(projectID)
	if err != nil {
		return err
	}
	pr.Tags.Add(tags...)

	if err:=pr.Save();err!=nil{
		return err
	}
	if err:=app.Repo().UpdateIndexForPath(pr.Path());err!=nil{
		return err
	}

	return nil
}
