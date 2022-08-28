package shuffler

import (
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"bitbucket.org/yasnikoff/shuffler/tagging"
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"bitbucket.org/yasnikoff/shuffler/actions"
)


// Main package's structure.
// Holds metadata about a directory with user data.
type Project struct {
	Header
	repoClient
	NameData      string `json:"name, omitempty"`
	Origin        id.ID `json:"origin, omitempty"`
	SigData       *parameters.Signature `json:"signature, omitempty"`
	Deps          *Dependencies `json:"dependencies,omitempty"`
	Tags          tagging.Collection `json:"tags"`
	collections   []id.ID        `json:"collections,omitempty"`
	ExcludedNames []string `json:"excluded names,omitempty"`
	ExcludedGlobs []string `json:"excluded globs,omitempty"`
	Actions       *actions.Collection        `json:"actions, omitempty"`
}

func NewProject() *Project {
	return DocTypeProject.New().(*Project)
}

func newProject() Doc {
	pr := &Project{
		SigData:parameters.NewSignature(nil),
		Tags:tagging.Collection{},
		Actions:actions.NewCollection(),
	}
	pr.Deps = newDeps(pr)
	return pr
}
