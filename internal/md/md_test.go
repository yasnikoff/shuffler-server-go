package md_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/internal/md"
	"os"
	"path/filepath"
	"bitbucket.org/yasnikoff/goutil/fsutil/patherrors"
	"bitbucket.org/yasnikoff/shuffler/tests"
)

func TestDirInfo_emptyMdDir(t *testing.T) {
	t.Parallel()
	f := tests.NewDirFx(t)
	defer f.Cleanup()

	if err := os.MkdirAll(md.DirPath(f.Path()), os.ModeDir); err != nil {
		t.Fatal(err)
	}
	path := f.Path()
	dirUnderTest := "dir with empty metadata folder"

	if ok := md.Initialized(path); !ok {
		t.Errorf("For &s got md.Initialized(d)==false. Want true.", dirUnderTest)
	}

	direrr := md.DirError(md.GetInfo(path))

	if direrr != nil {
		t.Errorf("For &s got md.DirError(d)=%v. Want nil.", dirUnderTest, direrr)
	}
}

func TestDirInfo_notExist(t *testing.T) {
	t.Parallel()
	f := tests.NewDirFx(t)
	defer f.Cleanup()
	path := filepath.Join(f.Path(), "noSuchPath")
	dirUnderTest := "non-existent folder"

	direrr := md.DirError(md.GetInfo(path))

	if !patherrors.Is(direrr, patherrors.DoesNotExist) {
		t.Errorf("For &s got md.DirError(d) is not 'NotExists'.", dirUnderTest)
	}

}

