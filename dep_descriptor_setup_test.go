package shuffler_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/tests/dependency_fixtures"
	"bitbucket.org/yasnikoff/shuffler/fs"
	"path/filepath"
)

type DepSetupFixture struct {
	dependency_fixtures.DeepDepFx
}

func NewDepSetupFx(t *testing.T) *DepSetupFixture {
	f := &DepSetupFixture{DeepDepFx:*dependency_fixtures.NewDeepDepFx(t, "dep", "depdep")}
	return f
}

// Check checks existence of dependency files (if depsExist==false, files must not exist).
// Also checks if the project files exist (they must always be present).
func (f *DepSetupFixture)CheckFilesExist(depsExist bool, condition string) {

	var requirement string = "must exist."
	if !depsExist {
		requirement = "must not exist."
	}
	//
	for _, fileFx := range f.User.FileFxMap() {
		file := fileFx.Path()
		if !fs.PathExists(file) {
			f.Errorf("%s %s %q %s", condition, "project file", file, "must exist.")
		}
	}
	for _, fileFx := range f.Dep.FileFxMap() {
		file := f.User.Abs(filepath.Join(f.DepPath, fileFx.SubPath()))
		if fs.PathExists(file) != depsExist {
			f.Errorf("%s %s %q %s", condition, "dependency file", file, requirement)
		}
	}
	for _, fileFx := range f.DepDep.FileFxMap() {
		file := f.User.Abs(filepath.Join(f.DepPath, f.DepDepPath, fileFx.SubPath()))
		if fs.PathExists(file) != depsExist {
			f.Errorf("%s %s %q %s", condition, "2-nd order dependency file", file, requirement)
		}
	}

}



func TestDepUpdate(t *testing.T) {
	//t.Skip()
	f := NewDepSetupFx(t)
	//defer f.Cleanup()

	f.CheckFilesExist(false, "Before Update() method of project dependencies")
	f.User.UpdateAllDeps()
	f.CheckFilesExist(true, "After Update() method of project dependencies")
	f.User.CleanUpAllDeps()
	f.CheckFilesExist(false, "After Cleanup() method of project dependencies")
	f.User.UpdateAllDeps()
	f.CheckFilesExist(true, "After Update() method of project dependencies")

}
