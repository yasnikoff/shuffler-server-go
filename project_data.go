package shuffler

import (
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"bitbucket.org/yasnikoff/shuffler/tagging"
	"bitbucket.org/yasnikoff/shuffler/actions"
)

type ProjectData struct {
	Header
	Name          string `json:"name, omitempty"`
	Origin        id.ID `json:"origin, omitempty"`
	Sig           *parameters.Signature `json:"signature, omitempty"`
	Deps          *Dependencies `json:"dependencies,omitempty"`
	Tags          *tagging.Collection `json:"tags, omitempty"`
	Collections   []id.ID        `json:"collections,omitempty"`
	ExcludedNames []string `json:"excluded names,omitempty"`
	ExcludedGlobs []string `json:"excluded globs,omitempty"`
	Actions       *actions.Collection        `json:"actions, omitempty"`
}

func (pr *Project)Data() interface{} {
	return &ProjectData{
		Header:pr.Header,
		Name:pr.Name(),
		Origin:pr.Origin,
		Sig:pr.Signature(),
		Deps:pr.Deps,
		Tags:&pr.Tags,
		Collections:pr.collections,
		ExcludedNames:pr.ExcludedNames,
		ExcludedGlobs:pr.ExcludedGlobs,
		Actions:pr.Actions,
	}
}

