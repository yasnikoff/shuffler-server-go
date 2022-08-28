/*

THIS FILE IS GENERATED AUTOMATICALLY
DO NOT EDIT THE CONTENTS

*/

package shuffler

import (
	"fmt"
)

var DocTypeProject *DocType
var afterSave func(*Project) error

func init() {
	DocTypeProject = RegisterType("Project", newProject)
}


// SaveTo saves Project to the path.
func (pr *Project) SaveTo(path string) error {
	if err := SaveDoc(pr, path); err != nil {
		return err
	}
	if afterSave != nil {
		return afterSave(pr)
	}
	return nil
}
// Save saves Project to the its path.
func (pr *Project) Save() error {
	return pr.SaveTo(pr.Path())
}

// Reload reads document data from disk
func (pr *Project) Reload() error {
	d, err := LoadDoc(pr.Path())
	if err != nil {
		return err
	}
	pr, ok := (d).(*Project)
	if !ok {
		return fmt.Errorf("can't convert doc from %q to type Project", pr.Path())
	}
	//TODO: !? fix this in the template
	*pr = *pr
	return nil
}
//
//func New()*Project{
//	return doc.DocTypes["Project"].New().(*Project)
//}
/*

func (pr *Project) Type() *DocType {
	return DocTypeProject
}
*/

func (pr *Project)String() string {
	return fmt.Sprintf("%s {%s}", pr.Type().Name(), pr.ID())
}

// Load returns Project saved in the filesystem.
func LoadProject(path string) (*Project, error) {

	pr, err := LoadDoc(path)
	if err != nil {
		return nil, err
	}
	result, ok := pr.(*Project)
	if !ok {
		return nil, &WrongDocTypeError{Doc:result, ExpectedType:"Project"}
	}
	return result, nil
}

