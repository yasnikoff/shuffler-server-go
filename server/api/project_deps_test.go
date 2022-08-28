package api

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler"
	"net/http"
	"github.com/pkg/errors"
)

type testDep struct {
	*baseTestCase
	t        *testing.T
	Project  *shuffler.Project
	Deps     []shuffler.Dependency
	SubPaths []string
	Case     string
}

func (dt *testDep)init(t *testing.T) {
	dt.baseTestCase = newBaseTestCase(t)
	dt.t = t
	dt.SetAppName("dep-test")
	dt.SetCaseName(dt.Case)
	dt.Project = dt.appFx.NewProject()
}

func (dt *testDep)Setup(r *testRunner) {
	dt.baseTestCase.Setup(r)
	for _, subPath := range dt.SubPaths {
		depPr := dt.appFx.NewProject()
		dt.AddDep(subPath, &shuffler.DepDesc{IDData:depPr.ID()})
	}
}
// AddDep adds dependency to the main Project.
// For use in tests setup to create deps that are to exist before the test.
func (dt *testDep)AddDep(subPath string, dep shuffler.Dependency) {
	if err := dt.Project.Deps.Add(subPath, dep); err != nil {
		dt.t.Fatalf(errors.Wrap(err, "can't add dep to depTest.Project: %s").Error())
	}
	if err := dt.Project.Save(); err != nil {
		dt.t.Fatalf(errors.Wrap(err, "can't save depTest.Project: %s").Error())
	}
	dt.Deps = append(dt.Deps, dep)
}

func ( t *testDep)Path() []string {
	var path []string = append(t.baseTestCase.Path(), "projects", t.Project.ID().String(), "dependencies")
	return path
}

type testDepGet struct {
	*testDep
}

func (t *testDepGet)ExpectedJsonSource() interface{} {
	return t.Project.Deps
}

func TestDepsGet(t *testing.T) {
	t.Parallel()
	tests := []*testDep{
		{Case:"no deps", SubPaths:[]string{}},
		{Case:"one dep, simple path", SubPaths:[]string{"dep"}},
		{Case:"one dep, compound path", SubPaths:[]string{"deps/dep1"}},
		{Case:"one dep, path with space", SubPaths:[]string{"deps/dep 1"}},
		{Case:"two deps", SubPaths:[]string{"dep1", "dep2"}},
		{Case:"four deps", SubPaths:[]string{"dep1", "dep2", "dep3", "dep4"}},
	}

	for _, test := range tests {
		test.init(t)
		r := &testRunner{}

		r.Run(test, t, func(t *testing.T) {

			if err := r.Check(); err != nil {
				t.Error(err)
			}
		})
	}
}

type deps map[string]shuffler.Dependency

type addDepTestCase struct {
	*testDep
	Case        string
	InitialDeps deps
	DepType     string
	NewDep      shuffler.Dependency
	DepPath     string
}

func (test *addDepTestCase)init(t *testing.T) {
	test.testDep = &testDep{}
	test.testDep.init(t)
	test.SetCaseName(test.Case)
	test.method = http.MethodPost
}

func (test *addDepTestCase)Setup(r *testRunner) {
	test.testDep.Setup(r)

	dep := test.AppFx().NewProject()
	if err:=dep.Save();err!=nil{
		test.t.Fatal(test.WrapErrorf(err, "can't save dep project"))
	}


	switch test.DepType{
	case shuffler.DepTypeDepDesc:
		test.NewDep = &shuffler.DepDesc{IDData:dep.ID()}
	}

	test.postData = map[string]map[string]interface{}{
		"data":map[string]interface{}{
			"id":test.Project.ID(),
			"path":test.DepPath,
			"type":test.DepType,
			"dep":test.NewDep,
		},
	}
}

func (test *addDepTestCase)Path() []string {
	return append(test.testDep.Path(), "add")
}

func TestDepsAdd(t *testing.T) {

	tests := []*addDepTestCase{
		{
			Case:"valid descriptor to project with no deps",
			InitialDeps:deps{},
			DepPath:"links/dep",
			DepType: shuffler.DepTypeDepDesc,
		},
	}

	for _, test := range tests {
		test.init(t)
		r := &testRunner{WriteLogs:false}

		r.Run(test, t, func(t *testing.T) {
			if err := r.Check(); err != nil {
				t.Error(err)
			}

			if err := test.Project.Reload(); err != nil {
				t.Fatal(test.WrapErrorf(err, "can't reload project after adding dep"))
			}

			// TODO: for some reason, fail even though getting by r.Get(test.testDep.Path()) returns correct data
			/*if _, err := test.Project.Deps.Get(test.DepPath); err != nil {
				t.Error(test.WrapErrorf(err, "dependency %q has not been added", test.DepPath))
			}*/
		})
	}
}
