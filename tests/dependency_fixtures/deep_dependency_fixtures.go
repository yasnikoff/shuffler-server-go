package dependency_fixtures

import (
	"bitbucket.org/yasnikoff/shuffler/tests/project_fixtures"
	"testing"
)

// DeepDepFx manages test case when project dependency has dependency too.
type DeepDepFx struct {
	DepFixture
	// dependency of dependency
	DepDep     *project_fixtures.ProjectFixture

	DepPath    string
	DepDepPath string
}

// NewDeepDepFx setups and returns new DeepDepFx. Project (DeepDepFx.User.Project) is set up to have dependency
// (DeepDepFx.Dep.Project) at the path "depPath".
// Dependency Project (DeepDepFx.Dep.Project) has its own dependency (DeepDepFx.DepDep.Project)
// which is mounted to Dep at depDepPath.
func NewDeepDepFx(t *testing.T, depPath string, depDepPath string) *DeepDepFx{
	f := &DeepDepFx{DepFixture:*New(t)}

	f.DepPath = depPath
	f.DepDepPath = depDepPath

	f.DepDep = f.NewProjectFx("depdep")

	f.User.NewRandFile("project_file")
	f.Dep.NewRandFile("dep_file")
	f.DepDep.NewRandFile("dep_dep_file")

	f.Dep.AddDep(depDepPath, f.DepDep, nil)
	f.User.AddDep(depPath, f.Dep, nil)

	f.UpdateIndex()
	return f
}
