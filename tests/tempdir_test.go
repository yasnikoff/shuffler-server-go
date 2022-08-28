package tests

import ("testing"
	"bitbucket.org/yasnikoff/shuffler/fs"
	"io/ioutil"
	"reflect"
	"path/filepath"
)

type Pather interface{
	Path() string
}

type Child interface{
	Pather
	Parent()*DirFx
	SubPath()string
}

func checkChild(t *testing.T, ch Child){
	parentPath:=ch.Parent().Path()
	childPath:=ch.Path()
	subPath:=ch.SubPath()

	if childPath!=filepath.Join(parentPath, subPath){
		t.Errorf("child.Path() != child.Parent().Path() + child.SubPath().\nPath: %q\nParent path: %q\nsubPath:%q",subPath,parentPath,subPath )
	}
}

func TestTempDir(t *testing.T) {
	f :=NewDirFx(t)
	defer f.Cleanup()

	if !fs.PathExists(f.Path()){
		t.Errorf("NewTempDir didn't create a directory at %q.", f.Path())
	}

	f.Cleanup()

	if fs.PathExists(f.Path()){
		t.Errorf("tempDir.Cleanup() didn't delete tempDir's directory %q", f.Path())
	}
}

func TestNewFile(t *testing.T) {
	f :=NewDirFx(t)
	defer f.Cleanup()

	content:=[]byte{0x11,0x10}
	path:="file"
	fileFx:= f.NewFileFx(path, content)
	loaded, err:=ioutil.ReadFile(fileFx.Path())
	if err!=nil{
		t.Fatal(err)
	}
	if !reflect.DeepEqual(content, loaded){
		t.Errorf("bytes read from the file creataed by DirFx.NewFile() is not equal to the saved bytes.\nread:%x\nwant:%x", loaded,content)
	}

}

func TestNewFileFx(t *testing.T) {
	f:=NewDirFx(t)
	defer f.Cleanup()

	content:=[]byte{0x01,0x11,0x10,0x00}
	subPath:="file"
	fileFx:=f.NewFileFx(subPath, content)
	loaded,err:=ioutil.ReadFile(fileFx.Path())
	if err!=nil{
		t.Errorf("can't fileFx's file %q\nerror:\n%s", fileFx.Path(), err)
	}
	if !reflect.DeepEqual(loaded,content){
		t.Errorf("bytes read from the file creataed by DirFx.NewFileFx() is not equal to the saved bytes.\nread:%x\nwant:%x", loaded,content)
	}
	if !reflect.DeepEqual(f.FileFxMap(), map[string]*FileFx{subPath:fileFx}){
		t.Errorf("FileFx created with DirFx.NewFileFx() is not in DirFx's children slice.")
	}
	checkChild(t, fileFx)

}

func TestNewDirFx(t *testing.T) {
	f :=NewDirFx(t)
	defer f.Cleanup()

	subPath :="dir"
	f2:=f.NewDirFx(subPath)

	if !fs.PathExists(f2.Path()){
		t.Errorf("DirFx.NewDirFx() didn't created directory at %s", f2.Path())
	}
	if !reflect.DeepEqual(f.DirFxMap(), map[string]*DirFx{subPath:f2}){
		t.Errorf("DirFx created with DirFx.NewDirFx() is not in DirFx's children slice.")
	}
	checkChild(t, f2)
}
