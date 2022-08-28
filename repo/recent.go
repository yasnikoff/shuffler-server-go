package repo

import (
	"bitbucket.org/yasnikoff/shuffler/types/id"
)

func (r *Repo)Recent() []id.ID {
	return r.RecentData.Elements()
}

func (r *Repo)AddToRecent(ID id.ID) error{
	r.RecentData.Add(ID)
	return r.Save()
}



