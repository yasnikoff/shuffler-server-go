package md
import (
	"bitbucket.org/yasnikoff/shuffler/fs"
)

// Info holds filesystem-related info about path and its metadata
type Info struct{
	Root pathInfo
	Dir pathInfo
	File pathInfo
}

// pathInfo holds info about filesystem path
type pathInfo struct{
	Path string
	Exists bool
	IsDir bool
}

func getPathInfo(path string)pathInfo{
	var inf = fs.GetInfo(path)
	return pathInfo{inf.Path, fs.Exists(inf), fs.IsDir(inf)}
}

// GetInfo returns filesystem-related info about path and its metadata
func GetInfo(path string)*Info{
	return &Info{
		Root:getPathInfo(path),
		Dir:getPathInfo(DirPath(path)),
		File:getPathInfo(FilePath(path)),
	}
}
