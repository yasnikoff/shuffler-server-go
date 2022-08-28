// package blender provides shuffler Action for running blender python scripts.
package blender

import (
	"fmt"
	"bitbucket.org/yasnikoff/shuffler/actions"
	"bitbucket.org/yasnikoff/shuffler/actions/shell"
	"path/filepath"
)

const typeName = "blender"

// blender process will return this exit code if python exception is raised
const pythonExitCode = 2

var ScriptNotSpecifiedErr = fmt.Errorf("Script file is not specified.")
var BlendNotSpecifiedErr = fmt.Errorf(".blend file is not specified.")

func init() {
	actions.Register(typeName, New)
}

func New() actions.Action {
	return &Action{}
}

type Action struct {
	shell.Action
	BlendFile string `json:"blendfile"`
	Py        string  `json:"script"`
	PyArgs    []string `json:"args"`
}

func (a *Action)Run() ([]byte, error) {
	if a.Py == "" {
		return nil, ScriptNotSpecifiedErr
	}
	/*if a.Blend == "" {
		return nil, BlendNotSpecifiedErr
	}*/
	// blender process will have another working dir than this process.
	// To prevent bugs, paths to .blend and .py files are converted to absolute paths.
	absRoot, err := filepath.Abs(a.Dir)
	if err != nil {
		return nil, err
	}

	args := []string{
		"-b",
		"--python-expr", runnerText,
		"--python-exit-code", fmt.Sprintf("%d", pythonExitCode),
		"--",
		filepath.Join(absRoot, a.Py),

	}
	if a.BlendFile != "" {
		args = append([]string{filepath.Join(absRoot, a.BlendFile)}, args...)
	}
	a.Args = append(args, a.PyArgs...)
	a.ExecPath = "blender"
	return a.Action.Run()

}

func (a *Action)fullPath(path string) string {
	return filepath.Join(a.Dir, path)
}

func (a *Action)Type() string {
	return typeName
}
