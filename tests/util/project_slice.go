package util

import (
	"bitbucket.org/yasnikoff/shuffler"
)

// ProjectSlicesEqualByID return true if both slices contain projects that are equal by their ID .
func ProjectSlicesEqualByID(s1, s2 []*shuffler.Project) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i].ID() != s2[i].ID() {
			return false
		}
		/*if !reflect.DeepEqual(s1[i],s2[i]){
			return false
		}*/
	}
	return true
}
