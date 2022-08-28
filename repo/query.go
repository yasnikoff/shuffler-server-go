package repo

import (
	"sync"
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"bitbucket.org/yasnikoff/shuffler/tagging"
	"bitbucket.org/yasnikoff/shuffler"
)



// Projects iterates over all projects in the Repo and returns slice of IDs of those that matched.
// May be expensive for large repositories.
func (r *Repo)Projects(match shuffler.ProjectMatcher) []id.ID {
	// TODO: bound concurrency
	// TODO: channel va lock
	wg := &sync.WaitGroup{}
	lock := sync.Mutex{}
	result := make([]id.ID, 0)
	for ID := range r.Index.Items() {
		wg.Add(1)
		go func(ID id.ID) {
			defer wg.Done()
			d, err := r.GetByID(ID)
			if err != nil {
				return
			}
			if match(d) {
				lock.Lock()
				result = append(result, d.ID())
				lock.Unlock()
			}
		}(ID)
	}
	wg.Wait()
	return result
}

func (r *Repo)SaveQuery(name, query string)error{
	r.savedQueriesOnce.Do(func(){
		if r.SavedQueriesData==nil{
			r.SavedQueriesData=make(map[string]string)
		}
	})
	r.SavedQueriesData[name]=query
	return r.Save()
}

func (r *Repo)SavedQueries()map[string]string{
	return r.SavedQueriesData
}


//----Tags---------------------------------------------------------------------

type TagQuery struct {
	tagging.Matcher
}

func (q *TagQuery)matchProject(pr *shuffler.Project)bool{
	return q.Match(&pr.Tags)
}

func (q *TagQuery)Execute(r shuffler.Repo)([]id.ID, error)  {
	r=r.(*Repo)
	return r.Projects(q.matchProject), nil
}

