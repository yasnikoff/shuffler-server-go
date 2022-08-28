package shuffler_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/tests"
	"reflect"
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"fmt"
	"bitbucket.org/yasnikoff/shuffler/tests/project_fixtures"
	"bitbucket.org/yasnikoff/shuffler/repo/repo_fixtures"
	"path/filepath"
	"bitbucket.org/yasnikoff/shuffler"
	"strings"
)

// Symbols "." and "_" mean exclusion for project files, but when present in the project path they must not cause exclusion of all project files.
// Project whose path contain the symbols must return exported files with FileList() as usual.
const PathWithExcludedSymbols = "path/with/.excluded/_symbols/some_inside.too"


// FileMapTest is used to test Project.FileMap(*parameters.Collection) function.
type FileMapTest interface {
	// Result returns actual result got from method under test
	Result() (map[string]string, error)
	// Expected returns expected result of Project.FileMap()
	Expected() map[string]string
	// Parameters() returns *parameters.Collection to be passed to Project.FileListP(p) call.
	Parameters() *parameters.Collection
	// Setup prepares user project fixture for the test. Typically adding dependencies and creating project files.
	Setup(user *FileMapTestFx) error
	// Msg() returns error description when the test fails.
	Msg(result, expected map[string]string) string
}

type FileMapTestFx struct {
	tests.DirFx // main fixture dir, sub dir of the root. Here all the tests take place. Sub path contains special symbols like "." and "_"
	root       *tests.DirFx
	userFx     *project_fixtures.ProjectFixture
	repoFx     *repo_fixtures.RepoFixture
	projectFxs map[string]*project_fixtures.ProjectFixture
}

func NewFileMapTestFx(t *testing.T) *FileMapTestFx {
	root := tests.NewDirFx(t)
	f := &FileMapTestFx{
		DirFx:*root.NewDirFx(PathWithExcludedSymbols),
		root:root,

	}
	f.repoFx = repo_fixtures.NewFromDirFx(f.NewDirFx("repo"))
	return f
}

func (f *FileMapTestFx)GetProjectFx(subPath string) (*project_fixtures.ProjectFixture, bool) {
	pFx, ok := f.projectFxs[subPath]
	return pFx, ok
}

func (f *FileMapTestFx)AddProject(subPath string, files ...string) *project_fixtures.ProjectFixture {
	prFx, ok := f.GetProjectFx(subPath)
	if ok {
		f.Fatalf("Path %q already contains project", subPath)
	}
	prFx = project_fixtures.NewFromDirFx(f.NewDirFx(subPath))
	for _, file := range files {
		prFx.NewRandFileFx(file)
	}
	f.repoFx.AddToIndex(prFx)
	prFx.Project.SetRepo(f.repoFx.Repo)

	return prFx
}

func (f *FileMapTestFx)Setup(baseDir string) {
	f.DirFx = *tests.NewDirFxFromPath(f.T, baseDir)
}

func (f *FileMapTestFx)Cleanup() {
	f.root.Cleanup()
}

func (f *FileMapTestFx)Reset() {
	f.Cleanup()
	f.Setup(f.Path())
}

func (f *FileMapTestFx)Result(params *parameters.Collection) map[string]string {
	pr := f.userFx.Project
	result, err := pr.FileMap(params, "")

	if err != nil {
		f.T.Fatal(err)
	}
	return result
}

func (f *FileMapTestFx)Run(t FileMapTest) {
	if err := t.Setup(f); err != nil {
		f.Fatal(err)
	}

	result, err := t.Result()
	if err != nil {
		f.Fatal(err)
	}
	expected := t.Expected()

	/*
		f.simplifyFileMap(result)
		f.simplifyFileMap(expected)
		f.Printf("case: %s\n", t.Msg(result, expected))
		f.Printf("   result:\n%v\nexpected:%v", result, expected)
		*/

	if !reflect.DeepEqual(result, expected) {
		f.simplifyFileMap(result)
		f.simplifyFileMap(expected)

		f.T.Errorf(t.Msg(result, expected))
	}

	f.Reset()

}

func (f *FileMapTestFx)simplifyFileMap(filemap map[string]string) {
	// substitute {abs root path} for root path prefix for readability sake
	absPlaceHolder := "{root abs path}/"
	absPathPrefix := filepath.ToSlash(f.Path())
	for dst, src := range filemap {
		if !strings.HasPrefix(src, absPathPrefix) {
			f.Fatalf("scr path %q is not an absolute subpath of fixture's path", src)
		}
		newSrc := strings.TrimPrefix(src, absPathPrefix)
		newSrc = strings.TrimPrefix(newSrc, "/")
		filemap[dst] = absPlaceHolder + newSrc
	}
}
// simple test, for project without dependencies.
type simpleFileMapTest struct {
	user     *shuffler.Project
	files    []string
	expected map[string]string
	msg      string
}

func (t *simpleFileMapTest)Result() (map[string]string, error) {
	return t.user.FileMap(t.Parameters(), "")

}

func (t *simpleFileMapTest)Expected() map[string]string {
	for dst, src := range t.expected {
		if !filepath.IsAbs(src) {
			t.expected[dst] = filepath.ToSlash(filepath.Join(t.user.Path(), src))
		}
	}
	return t.expected
}

func (t *simpleFileMapTest)Parameters() *parameters.Collection {
	return nil
}

func (t *simpleFileMapTest)Setup(f *FileMapTestFx) error {
	f.userFx = f.AddProject("user", t.files...)
	t.user = f.userFx.Project
	return nil
}

func (t *simpleFileMapTest)Msg(result, expected map[string]string) string {
	if t.msg == "" {
		t.msg = "didn't return expected file list"
	}
	fmtStr := fmt.Sprintf("project.FileList() %s.\nfiles:\t%%q\nGot:\t%%q\nWant:\t%%q", t.msg)
	return fmt.Sprintf(fmtStr, t.files, result, expected)
}

func TestFileLists_Without_Deps(t *testing.T) {
	f := NewFileMapTestFx(t)
	defer f.Cleanup()

	tests := []simpleFileMapTest{
		{
			msg:"must not return any files when called on empty project",
			files:[]string{},
			expected:map[string]string{},
		},
		{
			msg:"must not return files from project that start with \"_\"",
			files:[]string{"_excluded.txt"},
			expected:map[string]string{},
		},
		{
			msg:"must return '__init__.py' files) ",
			files:[]string{"__init__.py"},
			expected:map[string]string{"__init__.py":"__init__.py"},
		},
		{
			msg:"must not return files from project that start with \"_\"",
			files:[]string{"file1.txt", "file2.txt", "_excluded.txt"},
			expected:map[string]string{"file1.txt":"file1.txt", "file2.txt":"file2.txt"},
		},
		{
			msg:"must not return files from project that start with \"_\"",
			files:[]string{"dir/file1.txt", ".excluded/file.txt", "_excluded/file.txt"},
			expected:map[string]string{"dir/file1.txt":"dir/file1.txt"},
		},
		{
			msg:"must not return files from project that start with \".\"",
			files:[]string{".excluded.txt"},
			expected:map[string]string{},

		},
		{
			msg:"must not return files from project that start with \".\"",
			files:[]string{"file1.txt", ".excluded.txt"},
			expected:map[string]string{"file1.txt":"file1.txt"},

		},
		{
			msg:"must not return files from project that start with \".\"",
			files:[]string{"dir/.excluded.txt"},
			expected:map[string]string{},

		},
		{
			msg:"must not return files from project with extensions .blend1, .blend2, ... ",
			files:[]string{"file.blend1", "file.blend2", "file.blend9", "file.blend10", "file.blend01"},
			expected:map[string]string{},

		},
		{
			msg:"must return files from project with extension .blend",
			files:[]string{"file.blend"},
			expected:map[string]string{"file.blend":"file.blend"},

		},
		{
			msg:"must not return files from project that start with \".\"",
			files:[]string{"file1.txt", "dir/.excluded.txt"},
			expected:map[string]string{"file1.txt":"file1.txt"},

		},
		{
			msg:"didn't rturn expected file list",
			files:[]string{"dir/file1.txt",

				".excluded.txt",
				"_excluded.txt",
				"dir1/.excluded.txt",
				"dir1/_excluded.txt",
				"dir1/_excluded/file.txt",
				"dir1/.excluded/file.txt"},

			expected:map[string]string{"dir/file1.txt":"dir/file1.txt"},
		},

	}
	for _, test := range tests {
		f.Run(&test)
	}

}

type depData struct {
	src      string                 // path of dependency project, relative to fixture's root. Project will be created here during setup and path will be converted to abs.
	dst      string                 // relative dependency path
	params   *parameters.Collection // parameters of the dependency
	files    []string
	deps     []*depData             // dependencies of the dependency
	expected map[string]string      // dep's expected file map. Src paths are relative to the project. Will be converted to abs in the call to Expected().
}

func (d *depData)setup(f *FileMapTestFx, parentFx *project_fixtures.ProjectFixture) {

	depFx, ok := f.GetProjectFx(d.src)
	if !ok {
		depFx = f.AddProject(d.src, d.files...)
	} else {
		if len(d.files) > 0 {
			f.Fatalf("Attempt to setup the same project twice at %q", d.src)
		}
	}

	parentFx.AddDep(d.dst, depFx, d.params)

	for _, dep := range d.deps {
		dep.setup(f, depFx)
	}
	d.src = filepath.ToSlash(depFx.Path())
}

func (d *depData)addToFileMap(fm map[string]string, dstBaseSubPath string) {
	newDst := filepath.ToSlash(filepath.Join(dstBaseSubPath, d.dst))
	for _, dep := range d.deps {
		dep.addToFileMap(fm, newDst)
	}
	for dst, src := range d.expected {
		fm[filepath.ToSlash(filepath.Join(newDst, dst))] = filepath.ToSlash(filepath.Join(d.src, src))
	}
}

// FileMap of a Project with dependencies
type FileMapWithDepsTest struct {
	userFx *project_fixtures.ProjectFixture
	msg    string
	files  []string
	deps   []*depData
}

func (f *FileMapWithDepsTest)Setup(fm *FileMapTestFx) error {
	fm.userFx = fm.AddProject("user", f.files...)

	for _, dep := range f.deps {
		dep.setup(fm, fm.userFx)
	}
	f.userFx = fm.userFx
	return nil
}

func (f *FileMapWithDepsTest)Msg(result, expected map[string]string) string {
	msg := f.msg
	if msg == "" {
		msg = "project with dependencies"
	}
	msg = "FileMapWithDeps method of a " + strings.ToUpper(msg) + " returned wrong result.\nGot: %v\nWant:%v"
	return fmt.Sprintf(msg, result, expected)
}
func (f *FileMapWithDepsTest)Result() (map[string]string, error) {
	return f.userFx.Project.FileMapWithDeps(f.Parameters(), "")
}
func (f *FileMapWithDepsTest)Expected() map[string]string {
	result := make(map[string]string, 0)
	for _, dep := range f.deps {
		dep.addToFileMap(result, "")
	}
	for _, file := range f.files {
		result[file] = f.userFx.AbsSlash(file)
	}
	return result
}
func (f *FileMapWithDepsTest)Parameters() *parameters.Collection {
	return f.userFx.Project.DefaultParams()
}

func TestFileMapWithDeps(t *testing.T) {

	f := NewFileMapTestFx(t)
	defer f.Cleanup()

	tests := []*FileMapWithDepsTest{
		&FileMapWithDepsTest{
			msg:"empty project",
			files:[]string{},
			deps:[]*depData{},
		},
		&FileMapWithDepsTest{
			msg:"project with no dependencies",
			files:[]string{"file"},
			deps:[]*depData{},
		},
		&FileMapWithDepsTest{
			msg:"project with one dependency",
			files:[]string{"file"},
			deps:[]*depData{
				&depData{
					src:"depSrc",
					dst:"dep",
					files:[]string{"depfile"},
					expected:map[string]string{"depfile":"depfile"},
					deps:[]*depData{},
				},
			},
		},
		&FileMapWithDepsTest{
			msg:"project with private dependency",
			files:[]string{"file"},
			deps:[]*depData{
				&depData{
					src:"depSrc",
					dst:"_dep",
					files:[]string{"depfile"},
					expected:map[string]string{},
					deps:[]*depData{},
				},
			},
		},
		&FileMapWithDepsTest{
			msg:"project with 2-nd order dependency",
			files:[]string{"file"},
			deps:[]*depData{
				&depData{
					src:"depSrc",
					dst:"dep",
					files:[]string{"depfile"},
					expected:map[string]string{"depfile":"depfile"},
					deps:[]*depData{
						&depData{
							src:"depdepSrc",
							dst:"dep",
							files:[]string{"depdepfile"},
							expected:map[string]string{"depdepfile":"depdepfile"},
						},
					},
				},
			},
		},

	}

	for _, test := range tests {
		f.Run(test)
	}

}
