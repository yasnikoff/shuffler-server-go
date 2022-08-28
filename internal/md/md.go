// Package md deals with saving and loading metadata object to and from the filesystem.
package md

import (
	"bitbucket.org/yasnikoff/goutil/fsutil/patherrors"
	"bitbucket.org/yasnikoff/shuffler/fs"
	"os"
	"path/filepath"
)

const (
	DirName string = ".shuffler"
	FileName string = "metadata.json"
)

var DirSubPath string = DirName
var FileSubPath string = filepath.Join( DirSubPath,FileName)


// mdDirPath returns full path to the metadata directory for the path.
func DirPath(path string) string {
	return filepath.Join(path, DirSubPath)
}

func FilePath(path string) string{
	return filepath.Join(path, FileSubPath)
}

// MdError checks given path for some common errors related to its metadata dir.
func DirErrorForPath(path string) error {
	if err := DirError(GetInfo(path)); err != nil {
		return err
	}
	return nil
}

func DirError(info *Info) *patherrors.PathError {

	path := info.Root.Path

	if !info.Root.Exists {
		err := patherrors.New(patherrors.DoesNotExist, path)
		err.Add(ErrPathNotInitialized)
		return err
	}
	if !info.Root.IsDir {
		err := patherrors.New(patherrors.IsNotDir, path)
		return err
	}

	if !info.Dir.Exists {
		err := patherrors.New(ErrPathNotInitialized, path)
		err.Add(patherrors.DoesNotExist, DirSubPath)
		return err
	}

	if !info.Dir.IsDir {
		err := patherrors.New(ErrCorruptedMetadata, path)
		err.Add(patherrors.IsNotDir, DirSubPath)
		return err
	}

	return nil
}

// Initialized returns true if path already has metadata.
func Initialized(path string) bool {
	info := GetInfo(path)
	// TODO: check for metadata file.txt
	return info.Dir.IsDir
}

// Clear deletes all metadata from the path
func Clear(path string)error{
	p, err := fs.NormalizePath(path)
	if err != nil {
		return err
	}
	return os.RemoveAll(DirPath(p))
}


