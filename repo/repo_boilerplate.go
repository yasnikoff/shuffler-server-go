/*

THIS FILE IS GENERATED AUTOMATICALLY
DO NOT EDIT THE CONTENTS

*/

package repo

import (
	"fmt"
	"bitbucket.org/yasnikoff/shuffler"
)

var DocType *shuffler.DocType
var afterSave func(*Repo) error

func init() {
	DocType = shuffler.RegisterType("Repo", newDoc)
}

func newDoc() shuffler.Doc {
	return &Repo{}
}
/*

// SaveTo saves Repo to the path.
func (rp *Repo) SaveTo(path string) error {
	if err := shuffler.SaveDoc(rp, path); err != nil {
		return err
	}
	if afterSave != nil {
		return afterSave(rp)
	}
	return nil
}
// Save saves Repo to the its path.
func (rp *Repo) Save() error {
	return rp.SaveTo(rp.Path())
}
*/

// Reload reads document data from disk
func (rp *Repo) Reload() error {
	d, err := LoadRepo(rp.Path())
	if err != nil {
		return err
	}
	pr:= d

	*rp = *pr
	return nil
}
//
//func New()*Repo{
//	return doc.DocTypes["Repo"].New().(*Repo)
//}

func (rp *Repo) Type() *shuffler.DocType {
	return DocType
}

func (rp *Repo)String() string {
	return fmt.Sprintf("%s {%s}", rp.Type().Name(), rp.ID())
}

// Load returns Repo saved in the filesystem.
func LoadRepo(path string) (*Repo, error) {

	rp, err := shuffler.LoadDoc(path)
	if err != nil {
		return nil, err
	}
	result,ok:= rp.(*Repo)
	if !ok{
		return nil, &shuffler.WrongDocTypeError{ExpectedType:"Repo"}
	}

	return result, nil
}

