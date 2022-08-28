package shuffler

import "bitbucket.org/yasnikoff/shuffler/types/id"

type Tagger interface {
	Tags() []string
	AddTags(...string)
	SetTags(...string)
}

type ProjectQuery interface {
	Parse(string) error
	Execute(Repo) ([]id.ID, error)
	String() string
}

