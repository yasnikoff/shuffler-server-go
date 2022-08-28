package shuffler_test

import (
	"bitbucket.org/yasnikoff/shuffler/internal/md"
	"bitbucket.org/yasnikoff/shuffler/tests/project_fixtures"
	"path/filepath"
	"testing"
	"os"
	"io/ioutil"
	"strings"
	"bitbucket.org/yasnikoff/shuffler/tests"
	"bitbucket.org/yasnikoff/shuffler"
)

const (
	testDataDir = "testdata"
	invalidPath = "|||||\\ /://::#&&&????///"
	noSuchPath = "testdata/no_shuch_path"
)

// TODO: don't use Project as test document?

func TestLoad_nonExistentPath(t *testing.T) {

	funcUnderTest := "Load() called on non-existent path "

	pr, err := shuffler.LoadProject(noSuchPath)

	//if !patherrors.Is(err, patherrors.DoesNotExist) {
	//	t.Errorf(funcUnderTest+"must return a patherror.DoesNotExist Got %s (%s)", reflect.TypeOf(err), err)
	//}
	if err == nil {
		t.Errorf(funcUnderTest + "must return an error. Got nil.")
	}

	if pr != nil {
		t.Errorf(funcUnderTest + "must return nil project. Got %v", pr)
	}
}

func TestLoad_notDir(t *testing.T) {
	path := filepath.Join(testDataDir, "notDir.txt")

	pr, err := shuffler.LoadProject(path)

	funcUnderTest := "Load() called on a path that is not a directory "

	//if !patherrors.Is(err, patherrors.IsNotDir) {
	//	t.Errorf(funcUnderTest+"must return NotDirError. Got %v (%s)", reflect.TypeOf(err), err)
	//}

	if err == nil {
		t.Errorf(funcUnderTest + "must return an error. Got nil")
	}

	if pr != nil {
		t.Errorf(funcUnderTest + "must return nil project. Got %v", pr)
	}
}

func TestLoad_uninitialized(t *testing.T) {

	f := tests.NewDirFx(t)
	defer f.Cleanup()

	path := f.Path()

	funcUnderTest := "Load() called on uninitialized path "

	pr, err := shuffler.LoadProject(path)

	if err != md.ErrPathNotInitialized {
		t.Errorf(funcUnderTest + "must return an ErrPathNotInitialized. Got %v", err)
	}

	if pr != nil {
		t.Errorf(funcUnderTest + "must return nil project. Got %v", pr)
	}

}

func TestLoad_corrupted_mdDirIsNotDir(t *testing.T) {

	path := filepath.Join(testDataDir, "corrupted", "mdDirIsNotDir")

	funcUnderTest := "Load() called on a path with metadata dir being not a dir "

	pr, err := shuffler.LoadProject(path)

	if err == nil {
		t.Errorf(funcUnderTest + "must return an error.Got nil")
	}
	if pr != nil {
		t.Errorf(funcUnderTest + "must return nil project. Got %v", pr)
	}

}

func TestLoad_corrupted_invalidJson(t *testing.T) {
	path := filepath.Join(testDataDir, "corrupted", "invalidJson")

	funcUnderTest := "Load() called on a path with invalid json data in metadata.json file "
	pr, err := shuffler.LoadProject(path)

	if err == nil {
		t.Errorf(funcUnderTest + "must return an error.Got nil")
	}

	if pr != nil {
		t.Errorf(funcUnderTest + "must return nil project. Got %v", pr)
	}
}

func TestLoad_invalidPath(t *testing.T) {

	funcUnderTest := "Load() called on invalid path "
	pr, err := shuffler.LoadProject(invalidPath)

	if pr != nil {
		t.Errorf(funcUnderTest + "must return nil project. Got %v", pr)
	}

	if err == nil {
		t.Errorf(funcUnderTest + "must return an error.\nCalled on: \"%s\" Got %v", invalidPath, err)
	}
}

func TestLoad_normal(t *testing.T) {
	f := project_fixtures.New(t)
	defer f.Cleanup()

	funcUnderTest := "Load(path) called on valid initialized path "

	pr, err := shuffler.LoadProject(f.Path())
	if err != nil {
		t.Errorf(funcUnderTest + "must return nil error. Got: %s", err)
	}

	if pr == nil {
		t.Errorf(funcUnderTest + "must return project. Got nil.")
	}

	if pr.Path() != f.Path() {
		t.Errorf(funcUnderTest + "must return doc with doc.Path() equal to the path.")
		t.Logf("Project path:\n%q", f.Path())
		t.Logf("Loaded project path:\n%q", pr.Path())
	}
}


// New() tests

func TestNew_empty(t *testing.T) {

	f := tests.NewDirFx(t)
	defer f.Cleanup()

	funcUnderTest := "New() called on empty dir "

	pr := shuffler.NewProject()
	err := pr.SaveTo(f.Path())
	if err != nil {
		t.Errorf(funcUnderTest + "must return nil error. Got: %s", err)
	}

	if pr == nil {
		t.Errorf(funcUnderTest + "must return pointer to Project. Got nil.")
	}
}

func TestSave_alreadyInitialized(t *testing.T) {

	f := project_fixtures.New(t)
	defer f.Cleanup()

	funcUnderTest := "Save() called with path with another doc "

	pr2 := shuffler.NewProject()
	err := pr2.SaveTo(f.Path())

	if err == nil {
		t.Errorf(funcUnderTest + "must return error. Got: nil. Tried to save doc with id %s to path with doc htat have id %s", pr2.ID(), f.Project.ID())
	}
}

func TestNew_normal(t *testing.T) {

	f := tests.NewDirFx(t)
	defer f.Cleanup()

	path := f.Path()

	funcUnderTest := "New() called on existing, not initialized directory "

	pr := shuffler.NewProject()
	err := pr.SaveTo(path)
	if err != nil {
		t.Fatalf(funcUnderTest + "must initialie it. Error: %s", err)
	}

	if err != nil {
		t.Fatalf(funcUnderTest + "must initialie it. Error: %s", err)
	}

	expectedPath, err := filepath.Abs(path)
	realPath := pr.Path()
	if realPath != expectedPath {
		t.Errorf(funcUnderTest + "must return project with its .path field set to the path passed to the func. Expected: %s. Got: %s", expectedPath, realPath)
	}

}

func TestSave_invalidPath(t *testing.T) {

	pr := shuffler.NewProject()
	err := pr.SaveTo(invalidPath)

	if err == nil {
		t.Errorf("Init() on invalid path must return an error. Got %v", err)
	}
}

// Project's path must be absolute even when initialized with not absolute path.
func TestProject_pathIsAbs(t *testing.T) {

	// TODO: tests for relative paths ".", "some/path", "../some/path

	f := tests.NewDirFx(t)
	defer f.Cleanup()

	path := f.Path()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Chdir(cwd)
	}()

	if err := os.Chdir(path); err != nil {
		t.Fatal(err)
	}

	pr := shuffler.NewProject()

	if err := pr.SaveTo(path); err != nil {
		t.Fatal(err)
	}

	abspath, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	prpath := pr.Path()

	if prpath != abspath {
		t.Errorf("Project path is not set to absolute path to its root directory. p.Path()==\"%s\", want %s", prpath, abspath)
	}
}

func TestSave(t *testing.T) {
	f := tests.NewDirFx(t)
	defer f.Cleanup()

	path := f.Path()
	pr := shuffler.NewProject()
	err := pr.SaveTo(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := pr.SaveTo(""); err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadFile(md.FilePath(path))
	if err != nil {
		t.Fatal(err)
	}
	json := string(data)

	if !strings.Contains(json, pr.ID().String()) {
		t.Errorf("Saved json does not contain project's ID. Json: %s, ID: %s", json, pr.ID())
	}
}
func TestLoad(t *testing.T) {
	f := project_fixtures.New(t)

	defer f.Cleanup()

	pr, err := shuffler.LoadProject(f.Path())
	if err != nil {
		t.Fatal(err)
	}
	if pr.ID() != f.Project.ID() {
		t.Errorf("Loaded project's ID differs from the saved one. Got: %v, want: %v", f.Project.ID(), pr.ID())
	}
	if pr.Name() != f.Project.Name() {
		t.Errorf("Names for loaded and saved docs are different: %q != %q ", pr.Name(), f.Project.Name())
	}
}

func TestInitialized2(t *testing.T) {
	f := tests.NewDirFx(t)
	defer f.Cleanup()
	path := f.Path()
	pr := shuffler.NewProject()
	if err := pr.SaveTo(path); err != nil {
		t.Fatal(err)
	}

	funcUnderTest := "Initialized() called on initialized path "
	result := md.Initialized(path)
	if result != true {
		t.Errorf(funcUnderTest + "must return true. Got: %v", result)
	}
}

func TestSave_overwriteSelf(t *testing.T) {
	pf := project_fixtures.New(t)
	defer pf.Cleanup()

	if err := pf.Project.SaveTo(""); err != nil {
		t.Errorf("Save() called with empty string on saved doc returned: %v. Want nil.", err)
	}

}

// TODO: if doc has empty path than trying to save to the empty path must return error


func TestLoadToNewProject(t *testing.T) {
	f := tests.NewDirFx(t)
	defer f.Cleanup()

	name := "test project"
	p := shuffler.NewProject()
	p.SetName(name)
	if err := p.SaveTo(f.Path()); err != nil {
		t.Error(err)
	}
	p2, err := shuffler.LoadProject(f.Path())
	if err != nil {
		t.Error(err)
	}
	if p2.Name() != p.Name() {
		t.Errorf("Names of saved and loaded docs are different: %q != %q ", p.Name(), p2.Name())
	}
	if p2.ID() != p.ID() {
		t.Errorf("IDs for loaded and saved docs are different: %q != %q ", p.ID(), p2.ID())
	}
}

func TestLoadHeader(t *testing.T) {
	f := project_fixtures.New(t)
	defer f.Cleanup()

	h, err := shuffler.LoadHeader(f.Path())
	if err != nil {
		t.Fatal(err)
	}

	if h.Type().Name() != f.Project.Type().Name() {
		t.Errorf("Header loaded from path with %T must have the same document type. Got %#v. Want %#v.", f.Project, h.Type(), f.Project.Type())
	}
}
