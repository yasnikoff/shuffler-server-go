package app_test

import (
	"testing"
	"reflect"
	"bitbucket.org/yasnikoff/shuffler/app/app_fixtures"

	"bitbucket.org/yasnikoff/shuffler"
	"bitbucket.org/yasnikoff/shuffler/tests/project_fixtures"
	"bitbucket.org/yasnikoff/shuffler/fs"
	"github.com/termie/go-shutil"
	"bitbucket.org/yasnikoff/shuffler/tests"
	"github.com/pkg/errors"
	"fmt"
	"os"
)

func ProjectsAreEqual(p1, p2 *shuffler.Project) error {
	if p1 == nil {
		return fmt.Errorf("first project is nil")
	}
	if p2 == nil {
		return fmt.Errorf("second project is nil")
	}

	if p2.ID() != p1.ID() {
		return fmt.Errorf("ids are different:\n%s\n%s", p1.ID(), p2.ID())
	}

	// project.Name() returns dir's name if NameData is empty,
	// so the same project can have different names based on the directory
	if p2.NameData != p1.NameData {
		return fmt.Errorf("names are differnet:\n%s\n%s", p1.Name(), p2.Name())
	}
	// TODO: decide about project.Name() behaviour.
	// Is convenience of defaulting to dir's name worth the complication?

	if !reflect.DeepEqual(p1.Tags.Get(), p2.Tags.Get()) {
		return fmt.Errorf("Tags are different.\n%q\n%q", p1.Tags.Get(), p2.Tags.Get())

	}

	if !reflect.DeepEqual(p1.Deps.Items(), p2.Deps.Items()) {
		return fmt.Errorf("deps are different:\n%s\n%s", p1.Deps.Items(), p2.Deps.Items())

	}
	dstFiles, err := fs.FilesFromPath(p1.Path())
	if err != nil {
		return err
	}
	srcFiles, err := fs.FilesFromPath(p2.Path())
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(dstFiles, srcFiles) {
		return fmt.Errorf("files are different:\n%s\n%s", srcFiles, dstFiles)
	}

	return nil
}

type ProjectDuplicateFixture struct {
	app_fixtures.AppFixture
	SrcPrFx *project_fixtures.ProjectFixture
	SrcPr   *shuffler.Project
	DstPr   *shuffler.Project

	OldName string
	NewName string

	NewPath string
}

func NewProjectDuplicateFixture(t *testing.T) *ProjectDuplicateFixture {
	//t.Parallel()
	f := &ProjectDuplicateFixture{
		AppFixture:*app_fixtures.New(t),
		OldName:"origin",
		NewName:"duplicate",
	}
	f.SrcPrFx = project_fixtures.NewFromDirFx(f.NewDirFx(f.OldName))
	f.SrcPr = f.SrcPrFx.Project
	f.SrcPr.SetName(f.OldName)
	if err := f.Repo.AddToIndex(f.SrcPr); err != nil {
		f.Fatal(err)
	}
	return f
}

func (f *ProjectDuplicateFixture)Duplicate() {
	f.NewPath = f.Abs(f.NewName)
	var err error
	f.DstPr, err = f.App.DuplicateProject(f.SrcPr.ID().String(), f.NewName, f.NewPath)
	if err != nil {
		f.T.Fatal(err)
	}
}

func (f *ProjectDuplicateFixture)Files() (srcFiles, dstFiles []string) {
	var err error
	srcFiles, err = fs.FilesFromPath(f.SrcPr.Path())
	if err != nil {
		f.Fatal(err)
	}
	dstFiles, err = fs.FilesFromPath(f.DstPr.Path())
	if err != nil {
		f.Fatal(err)
	}
	return

}

func (f *ProjectDuplicateFixture)SaveSrc() {
	if err := f.SrcPr.Save(); err != nil {
		f.T.Fatal(err)
	}
}

func (f *ProjectDuplicateFixture)CheckProjects(CaseDescr string) {
	if f.DstPr == nil {
		f.T.Fatalf("%s: app.DuplicateProject() returned nil project", CaseDescr)
	}

	if f.DstPr.Name() != f.NewName {
		f.T.Errorf("%s: app.DuplicateProject() returend project with wrong name. Got: %q, want: %q", CaseDescr, f.DstPr.Name(), f.NewName)
	}

	if f.DstPr.Path() != f.NewPath {
		f.T.Errorf("%s: app.DuplicateProject() returned project with wrong path. Got: %q, want: %q", CaseDescr, f.DstPr.Path(), f.NewPath)
	}

	if f.DstPr.ID() == f.SrcPr.ID() {
		f.T.Errorf("%s: app.DuplacateProject() returned project with the same id as the original one.", CaseDescr)
	}

	if f.DstPr.Origin != f.SrcPr.ID() {
		f.Errorf("%s: app.DuplicateProject() returned project with its Origin property NOT set to source project. Origin: %v, want: %s", CaseDescr, f.SrcPr.Origin, f.DstPr.ID())
	}

	if !reflect.DeepEqual(f.SrcPr.Tags.Get(), f.DstPr.Tags.Get()) {
		f.T.Errorf("%s: app.DuplacateProject() returned project with tags different from the original one. Got:%q, want:%q", CaseDescr, f.DstPr.Tags.Get(), f.SrcPr.Tags.Get())

	}

	if !reflect.DeepEqual(f.SrcPr.Deps.Items(), f.DstPr.Deps.Items()) {
		f.T.Errorf("%s: app.DuplacateProject() returned project with dependencies different from the original one. Got:%q, want:%q", CaseDescr, f.DstPr.Deps.Items(), f.SrcPr.Deps.Items())

	}
	dstFiles, srcFiles := f.Files()
	if !reflect.DeepEqual(dstFiles, srcFiles) {
		f.T.Errorf("%s: app.DuplacateProject() didn't copy project files properly. dst files:%v, src files:%v", CaseDescr, dstFiles, srcFiles)
	}
}

func TestProjectDuplicateEmptyProject(t *testing.T) {
	f := NewProjectDuplicateFixture(t)
	defer f.Cleanup()

	f.Duplicate()
	f.CheckProjects("Empty project")

}

func TestProjectDuplicateProjectWithTags(t *testing.T) {
	f := NewProjectDuplicateFixture(t)
	defer f.Cleanup()

	tags := []string{"one", "two", "three"}
	f.SrcPr.Tags.Add(tags...)

	f.SaveSrc()
	f.Duplicate()
	f.CheckProjects("Project with tags")
}

func TestProjectDuplicateProjectWithDeps(t *testing.T) {
	f := NewProjectDuplicateFixture(t)
	defer f.Cleanup()

	depPr := f.NewProject()
	f.App.AddDep(f.SrcPr.ID().String(), "dep", depPr.ID().String())

	f.SaveSrc()
	f.Duplicate()
	f.CheckProjects("Project with deps")

}

func TestProjectDuplicateProjectWithFiles(t *testing.T) {

	tests := map[string][]string{
		"one file":{"file1"},
		"two files":{"file1", "file2"},
		"file in metadata dir":{"./.shuffler/file1"},
	}

	for caseName, files := range tests {

		f := NewProjectDuplicateFixture(t)
		defer func(f *ProjectDuplicateFixture) {
			f.Cleanup()
		}(f)

		for _, file := range files {
			f.SrcPrFx.NewRandFile(file)
		}

		f.SaveSrc()
		f.Duplicate()
		f.CheckProjects(caseName)
	}

}

type ProjectMoveFixture struct {
	*app_fixtures.AppFixture
	newPath    string
	OldProject *shuffler.Project
	NewProject *shuffler.Project
	Files      []string
	copyFx     *tests.DirFx
	WriteLogs  bool
}

type ProjectMoveOptions struct {
	Case string
	OldPath string
	NewPath string
	Files   []string
}

func NewProjectMoveFixture(t *testing.T, opt *ProjectMoveOptions) *ProjectMoveFixture {

	if opt.OldPath == "" {
		opt.OldPath = "old path"
	}

	if opt.NewPath == "" {
		opt.NewPath = "new path"
	}

	a := app_fixtures.New(t)
	pFx := a.NewProjectFx(opt.OldPath)

	f := &ProjectMoveFixture{
		AppFixture:a,
		OldProject:pFx.Project,
		newPath:opt.NewPath,
		copyFx:tests.NewDirFx(t),
		Files:opt.Files,
		WriteLogs:false,
	}

	if err := f.Repo.UpdateIndex(); err != nil {
		f.T.Fatal(err)
	}

	for _, file := range f.Files {
		pFx.NewRandFile(file)
	}

	return f
}

func (f *ProjectMoveFixture)Cleanup() {
	f.AppFixture.Cleanup()
	f.copyFx.Cleanup()
}

func (f *ProjectMoveFixture)TargetPath() string {
	return f.AppFixture.Abs(f.newPath)
}

func (f *ProjectMoveFixture)Check() {
	if err := ProjectsAreEqual(f.NewProject, f.OldProject); err != nil {
		f.T.Errorf(errors.Wrap(err, "Project is not the same after moving").Error())
	}
}

func (f *ProjectMoveFixture)Move() {
	var err error

	// make project copy
	f.Logf("recoding project state, makeing project's copy")

	// CopyTree returns error if the path exists
	if err := os.Remove(f.copyFx.Path()); err != nil {
		f.T.Fatalf("can't remove copyFx's dir", err)
	}
	if err := shutil.CopyTree(f.OldProject.Path(), f.copyFx.Path(), nil); err != nil {
		f.T.Fatal(err)
	}

	f.Logf("loading project copy")
	f.OldProject, err = shuffler.LoadProject(f.copyFx.Path())
	if err != nil {
		f.T.Fatal(err)
	}

	f.Logf("moving project")
	if _, err := f.AppFixture.App.MoveProject(f.OldProject.ID(), f.TargetPath()); err != nil {
		f.Fatal(err)
	}

	f.Logf("loading new project")
	f.NewProject, err = shuffler.LoadProject(f.Abs(f.newPath))
	if err != nil {
		f.T.Fatal(err)
	}
}

func (f *ProjectMoveFixture)Logf(format string, args ...interface{}) {
	if !f.WriteLogs {
		return
	}
	f.T.Logf(format, args...)
}

func TestProjectMoveNormal(t *testing.T) {
	tests := []*ProjectMoveOptions{
		{},
		{Files:[]string{"file.txt"}},
	}
	for _, test := range tests {
		f := NewProjectMoveFixture(t, test)
		defer func(f *ProjectMoveFixture) {
			f.Cleanup()
		}(f)
		f.Move()
		f.Check()
	}
}

func TestProjectMoveToNotEmptyDir(t *testing.T){

	a:=app_fixtures.New(t)
	defer a.Cleanup()

	p:=a.NewProjectFx("project")
	if err:=a.Repo.AddToIndex(p.Project);err!=nil{
		t.Fatal(err)
	}

	d:=a.NewDirFx("target")
	d.NewRandFile("file")

	_, err:=a.App.MoveProject(p.Project.ID(), d.Path())
	if err==nil{
		t.Errorf("app.MoveProject() moved project to not empty dir. Expect it fail with an error.")
	}

}
