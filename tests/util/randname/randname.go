// Package randname provides functions for generating random names for directories and files.
// It is intended to be used in test fixtures that need random names for temporary files and directories.
// Seeding is omitted intentionally, so generated names are always the same.
// This helps in detecting tests that fail to cleanup after themselves.
package randname

import (
	"math/rand"
	"bitbucket.org/yasnikoff/shuffler/fs"
	"path/filepath"
	"fmt"
	"os"
)
var DefaultNameOptions NameOptions
var DefaultSubPathOptions SubPathOptions

func init(){
	DefaultNameOptions  = NameOptions{
		Len:16,
		Runes:[]rune(`0123456789abcdef`),
	}
	DefaultSubPathOptions  = SubPathOptions{
		NameOptions:DefaultNameOptions,
		Tries:32,
	}

}

// NameOptions provides options for random name creation.
// Len - length of the name.
// Runes - runes to use in the name. There is no guarantee that all of the runes will be used.
type NameOptions struct {
	Len uint8
	Runes []rune
}



// NewNameOptions returns new NameOptions struct initialized with values from DefaultNameOptions.
func NewNameOptions()*NameOptions{
	return &NameOptions{
		Len:DefaultNameOptions.Len,
		Runes:DefaultNameOptions.Runes,
	}
}

// SubPathOptions provides options for random sub path creation.
// Tries - number of times to try to generate sub path name before giving up.
type SubPathOptions struct{
	NameOptions
	Tries uint8
}


// NewSubPathOptions returns new NameOptions struct initialized with values from DefaultSubPathOptions.
func NewSubPathOptions()*SubPathOptions{
	return &SubPathOptions{
		NameOptions:*NewNameOptions(),
		Tries:DefaultSubPathOptions.Tries,
	}
}


// New returns random name with specified options.
func NewName(opt *NameOptions) string {
	if opt==nil{
		opt=&DefaultNameOptions
	}
	result := make([]rune, 0, opt.Len)
	runesCount:=len(opt.Runes)
	for i:= uint8(0); i < opt.Len; i++ {
		result = append(result, opt.Runes[rand.Intn(runesCount)])
	}
	return string(result)
}


// NewSubpath returns random name that can be used as subdirectory or file name inside the path.
// Panics if the path does not exists or after trying opt.Tries times
// and checking that all the generated sub paths already exist.
func NewSubPath(path string, opt *SubPathOptions) string {
	if opt==nil{
		opt=&DefaultSubPathOptions
	}

	var name string
	if !fs.PathExists(path) {
		if err:=os.MkdirAll(path, os.ModeDir);err!=nil{
			panic(fmt.Errorf("Can't create dir %q", path))
		}
	}
	for i := uint8(0); i < opt.Tries; i++ {
		name = NewName(&(opt.NameOptions))
		if !fs.PathExists(filepath.Join(path, name)) {
			return name
		}
	}
	panic(fmt.Errorf("Can't create new temp dir. Tried pick a random name %d times, all dirs exist.", opt.Tries))
}
