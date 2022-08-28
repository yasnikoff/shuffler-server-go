package tests

import (
	"path/filepath"
	"os"
	"bitbucket.org/yasnikoff/shuffler/tests/util/randname"
	"testing"
	"io/ioutil"
)

var TestsRoot = filepath.Join(os.TempDir(), "shuffler tests")


type Fixture struct {
	T *testing.T
}

func (f *Fixture)Fatal(err error) {
	f.T.Fatal(err)
}

func (f *Fixture)Fatalf(msg string, args ...interface{}) {
	f.T.Fatalf(msg, args...)
}

func (f *Fixture)Error(err error) {
	f.T.Error(err)
}

func (f *Fixture)Errorf(msg string, args ...interface{}) {
	f.T.Errorf(msg, args...)
}

func (f *Fixture)Printf(msg string, args ...interface{}) {
	f.T.Logf(msg, args...)
}

func (f *Fixture)Println(args ...interface{}) {
	f.T.Log(args...)
	f.T.Logf("\n")
}

type DirFx struct {
	Fixture
	subPath         string
	parent          *DirFx
	childDirFxMap  map[string]*DirFx
	childFileFxMap map[string]*FileFx
	// TODO: checks for not having identical keys in both maps
}

func NewDirFx(t *testing.T) *DirFx {
	return NewDirFxFromPath(t, filepath.Join(TestsRoot, randname.NewSubPath(TestsRoot, nil)))
}

func NewDirFxFromPath(t *testing.T, path string) *DirFx {
	td := &DirFx{
		subPath:path,
		Fixture:Fixture{T:t},
		childDirFxMap:make(map[string]*DirFx, 0),
		childFileFxMap:make(map[string]*FileFx, 0),
	}
	td.Setup()
	return td
}

func AbsPath(relPath string) string {
	return filepath.Join(TestsRoot, relPath)
}

func (td *DirFx)Setup(){
	if err:= os.MkdirAll(td.Path(), os.ModeDir);err!=nil{
		td.Fatal(err)
	}
}

func (td *DirFx)Cleanup() {
	if err:=os.RemoveAll(td.Path());err!=nil{
		td.Fatal(err)
	}
}

func (td *DirFx)Parent()*DirFx{
	return td.parent
}

func (td *DirFx)Path() string {
	var parentPath string
	if td.parent!=nil{
		parentPath=td.parent.Path()
	}
	return filepath.Join(parentPath, td.subPath)
}
func (td *DirFx)SubPath() string {
	return td.subPath
}

func (td *DirFx)Abs(path string) string {
	return filepath.Join(td.Path(), path)
}
func (td *DirFx)AbsSlash(path string) string {
	return filepath.ToSlash(td.Abs(path))
}

func (td *DirFx)DirFxMap() map[string]*DirFx {
	return td.childDirFxMap
}
func (td *DirFx)FileFxMap() map[string]*FileFx {
	return td.childFileFxMap
}

func (td *DirFx)NewFileFx(path string, content []byte) *FileFx {
	f := &FileFx{parent:td, subPath:path, content:content, Fixture:td.Fixture}
	f.write()
	td.childFileFxMap[f.subPath]=f
	return f
}

func (td *DirFx)newSubPath() string {
	return randname.NewSubPath(td.Path(), nil)
}

func (td *DirFx)NewFile(path string, content []byte) string {
	return td.NewFileFx(path, content).Path()
}

func (td *DirFx)NewTextFileFx(path string, content string) *FileFx {
	return td.NewFileFx(path, []byte(content))
}
func (td *DirFx)NewTextFile(path string, content string) string {
	return td.NewTextFileFx(path, content).Path()
}

func (td *DirFx)NewRandFileFx(path string) *FileFx {
	content := []byte(randname.NewName(nil))
	return td.NewFileFx(path, content)
}
func (td *DirFx)NewRandFile(path string) string {
	return td.NewRandFileFx(path).Path()
}

func (td *DirFx)NewDirFx(path string) *DirFx {
	d:=&DirFx{Fixture:td.Fixture, parent:td, subPath:path, childDirFxMap:make(map[string]*DirFx), childFileFxMap:make(map[string]*FileFx)}
	d.Setup()
	td.childDirFxMap[d.subPath]= d
	return d
}
func (td *DirFx)NewDir(path string) string {
	return td.NewDirFx(path).Path()
}

func (td *DirFx)NewRandDirFx() *DirFx {
	return td.NewDirFx(td.newSubPath())
}
func (td *DirFx)NewRandDir() string {
	return td.NewDirFx(td.newSubPath()).Path()
}

//-----------------------------------------------------------------
type FileFx struct {
	Fixture
	subPath string
	parent  *DirFx
	content []byte
}

func (f *FileFx)Path() string {
	var parentPath string
	if f.parent != nil {
		parentPath = f.parent.Path()
	}
	return filepath.Join(parentPath, f.subPath)
}

func (f *FileFx)SubPath()string{
	return f.subPath
}

func (f *FileFx)Parent()*DirFx{
	return f.parent
}

func (f *FileFx)Content() []byte {
	return f.content
}

func (f *FileFx)ContentString() string {
	return string(f.content)
}

func (f *FileFx)write() {
	if err := os.MkdirAll(filepath.Dir(f.Path()), os.ModeDir); err != nil {
		f.Fatal(err)
	}
	if err := ioutil.WriteFile(f.Path(), f.content, 0x777); err != nil {
		f.Fatalf("Error in writing file fixture:\n%s", f.Path(), err)
	}
}
//-----------------------------------------------------------------
