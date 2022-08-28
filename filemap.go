package shuffler

import (
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"path/filepath"
	"os"
	"regexp"
	"strings"
	"github.com/pkg/errors"
	"fmt"
)

const filesPerDepEstimate int8 = 4 // estimate of how many files to expect from a dependency.

var DefaultExcludedGlobs = []string{""}
var DefaultExcludedNames = []string{`^[.]shuffler`, `^[.]`, `^[_]`, `\w*\.blend[0-9]+`} // TODO: move .blend[0-9] to more appropriate place
var DefaultForceIncludes = []string{`__init__.py`} // TODO: find better way to deal with '__init__.py' files

// FileMap returns map of source and destination paths needed for duplicating Project.
// Destination paths are keys and source paths are values. Only Project files are listed, i.e. files that do not come from dependencies.
// Because of that the map is not sufficient for setting up a directory.
// Destination paths are relative, so that the map can be applied to any directory.
// Source paths are absolute which isn't necessary for this function because all the files should come from the Project's directory
// but they are made so for being consistent with FileMapWithDeps, where absolute source paths a necessary because they come from different Projects.
func (pr *Project)FileMap(params *parameters.Collection, basePath string) (map[string]string, error) {

	result := make(map[string]string, 0)
	errMsg := fmt.Sprintf("can't create filemap for path %q and parameters %s", basePath, params)

	excludeNames, err := getRegExps(append(DefaultExcludedNames, pr.ExcludedNames...))
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}

	excludedGlobs := append(DefaultExcludedGlobs, pr.ExcludedGlobs...)


	for path := range pr.Deps.Items() {
		excludedGlobs = append(excludedGlobs, path)
	}

	if params == nil {
		params = pr.DefaultParams()
	}

	err = filepath.Walk(pr.Path(), func(path string, info os.FileInfo, e error) error {
		relPath, err := filepath.Rel(pr.Path(), path)
		if relPath == "." {
			// do not skip project dir and do nothing with it, just go ahead
			return nil
		}
		if err != nil {
			return err
		}

		var matched bool
		for _, re := range excludeNames {
			matched = re.MatchString(info.Name());
			if matched {
				break
			}
		}
		if !matched {
			for _, pattern := range excludedGlobs {
				matched, err = filepath.Match(pattern, relPath)
				if err != nil {
					return err
				}
				if matched {
					break
				}

			}
		}

		if matched {
			// path is to be discarded
			// check for force includes
			for _, pattern := range DefaultForceIncludes {
				forceInclude, err := regexp.MatchString(pattern, relPath)
				if err != nil {
					return err
				}
				matched = (!forceInclude)&&matched
			}

		}

		if matched && info.IsDir() {
			return filepath.SkipDir
		}
		if matched {
			return nil
		}

		if info.IsDir() {
			return nil
		}
		dstPath := filepath.ToSlash(filepath.Join(basePath, relPath))
		srcPath := filepath.ToSlash(filepath.Join(pr.Path(), relPath))
		// TODO: after implementing parameters processing srcPath adn dstPath might differ
		result[dstPath] = srcPath

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// FileMap returns map of source and destination paths needed to set up a directory.
// Destination paths are keys and source paths are values. Project files and files that come from dependencies are listed,
// so that the map is sufficient for a directory setup.
func (pr *Project)FileMapWithDeps(params *parameters.Collection, basePath string) (map[string]string, error) {

	result := make(map[string]string, int(filesPerDepEstimate) * len(pr.Deps.items))

	addToResult := func(fileMap map[string]string, baseDstPath string) {
		for dst, src := range fileMap {
			result[filepath.ToSlash(filepath.Join(baseDstPath, dst))] = filepath.ToSlash(src)
		}
	}

	for subPath, dep := range pr.Deps.items {
		if strings.HasPrefix(subPath, "_") || strings.HasPrefix(subPath, ".") {
			// TODO: more configurable rules for deps' paths.
			continue
		}
		desc := dep.Descriptor()
		depPr, err := pr.repo.GetByID(desc.ID())
		if err != nil {
			return nil, err
		}
		fm, err := depPr.FileMapWithDeps(desc.Parameters, "")
		if err != nil {
			return nil, err
		}
		addToResult(fm, subPath)

	}

	prFm, err := pr.FileMap(params, basePath)
	if err != nil {
		return nil, err
	}
	addToResult(prFm, basePath)
	return result, nil
}


//-----------------------------------------------------------------------------


// getRegExps takes a slice of regexp patterns and compiles them.
func getRegExps(patterns []string) ([]*regexp.Regexp, error) {
	regExps := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, errors.Wrapf(err, "can't compile RegExp pattern %q", pattern)
		}
		regExps[i] = re
	}
	return regExps, nil
}
