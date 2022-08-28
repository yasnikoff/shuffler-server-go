package shuffler

import (
	"path/filepath"
)


func (pr *Project)Name() string {
	if pr.NameData == "" {
		if pr.path == "" {
			return ""
		} else {
			return filepath.Base(pr.Path())
		}
	}
	return pr.NameData
}

func (pr *Project)SetName(name string){
	pr.NameData = name
}
