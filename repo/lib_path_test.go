package repo

import (
	"testing"
	"github.com/pkg/errors"
	"path/filepath"
	"fmt"
	"os"
)

func Test_LibPath(t *testing.T) {
	tests := []struct {
		Case    string

		// repo path (will be prefixed by absolute path
		// in order to get absolute path is cross-platform way)
		Repo    string
		Lib     string
		LibPath string
		AbsPath string
	}{
		{
			Case:".",
			Repo: "/abs/path/to/root",
			Lib: ".",
			LibPath: "p1",
			AbsPath:"/abs/path/to/root/p1",
		},
		{
			Case:"lib",
			Repo: "/abs/path/to/root",
			Lib: "lib",
			LibPath: "../p1",
			AbsPath:"/abs/path/to/root/p1",
		},
	}
	absPathPrefix := os.TempDir()
	for _, tc := range tests {
		r := New()
		r.SetPath(filepath.Join(absPathPrefix,tc.Repo))
		r.LibData = tc.Lib
		tc.AbsPath = filepath.Join(absPathPrefix, tc.AbsPath)
		tc := tc // capture range variable
		t.Run(fmt.Sprintf("%q-FromLibPath", tc.Case), func(t *testing.T) {
			absPath, err := r.FromLibPath(tc.LibPath)

			if err != nil {
				t.Fatal(errors.Wrapf(err, "case %q: can't create absPath from %q",tc.Case, tc.LibPath))
			}
			expected:=filepath.FromSlash(tc.AbsPath)
			if expected != absPath {
				t.Errorf("case %q: wrong abs path from lib path %q. Expected %q, got %q",
					tc.Case, tc.LibPath, expected, absPath)
			}
		})
		t.Run(fmt.Sprintf("%q-ToLibPath", tc.Case), func(t *testing.T) {
			libPath, err := r.ToLibPath(tc.AbsPath)

			if err != nil {
				t.Fatal(errors.Wrapf(err, "case %q: can't create libPath from %q",tc.Case, tc.AbsPath))
			}
			expected:=filepath.ToSlash(tc.LibPath)
			if expected != libPath {
				t.Errorf("case %q: wrong lib path from abs path %q. Expected %q, got %q",
					tc.Case, tc.AbsPath, expected,  libPath)
			}
		})


	}

}

