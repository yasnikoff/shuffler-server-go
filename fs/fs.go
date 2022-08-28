// Package fs provides file system utilities for the Shuffler project
package fs

import (
	"os"
	"path/filepath"
	"io"
	"fmt"
	"bitbucket.org/yasnikoff/goutil/errutil/errorlist"
	"sync"
	"bitbucket.org/yasnikoff/shuffler/fs/matchers"

)

// Normalize returns absolute path from given path.
func NormalizePath(path string) (string, error) {
	if path == "" {
		// don't normalize empty path because it's a zero value for strings and
		// may be here because of some bug. Force users to be explicit.
		return "", fmt.Errorf("path is empty")
	}

	if path == "." {
		path, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return path, nil
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.FromSlash(filepath.Clean(path)), nil
}

// GetInfo encapsulates os.FileInfo and os.PathError for filesystem path.
// This struct can be obtained once for a path and then used by multiple
// functions without reading the path's info again.
type PathInfo struct {
	Path string
	Stat os.FileInfo
	Err  error
}

// GetInfo returns pathIfo for given path.
func GetInfo(path string) *PathInfo {
	stat, err := os.Stat(path)
	return &PathInfo{path, stat, err}
}

// IsDir returns true if PathInfo represents directory
func IsDir(info *PathInfo) bool {
	return info.Err == nil && info.Stat.IsDir()
}

// Exists returns true if PathInfo represents existing and readable path.
func Exists(info *PathInfo) bool {
	return !os.IsNotExist(info.Err)
}

// Copy copies src file to dst.
func Copy(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only src, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()

	if err := os.MkdirAll(filepath.Dir(dst), os.ModeDir); err != nil {
		return err
	}

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func FilesFromPath(path string) ([]string, error) {
	return MatchFilesFromPath(path, matchers.All)
}

func MatchFilesFromPath(path string, matcher matchers.Matcher)([]string,error){
	result := []string{}
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err!=nil{
			return err
		}
		subPath, err := filepath.Rel(path, p)
		if err != nil {
			return err
		}
		subPath = filepath.ToSlash(subPath)
		if info.IsDir() {
			return nil
		}

		if !matcher.Match(subPath){
			return nil
		}
		result = append(result, subPath)

		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CopyFileMap(dstDir string, fileMap map[string]string) error {
	var errList *errorlist.ErrorList
	errors := make(chan error)

	errDone := make(chan struct{})
	go func() {
		for err := range errors {
			errList.AddNotNil(err)
		}
		errDone <- struct{}{}
	}()

	jobs := sync.WaitGroup{}
	jobs.Add(len(fileMap))
	for dst, src := range fileMap {
		go func(dst, src string) {
			errors <- Copy(filepath.Join(dstDir, dst), src)
			jobs.Done()
		}(dst, src)
	}

	jobs.Wait()
	close(errors)
	<-errDone
	return errList
}
