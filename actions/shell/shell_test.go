package shell

import (
	"path/filepath"
	"os/exec"
	"testing"

	"bitbucket.org/yasnikoff/shuffler/tests"
	"bitbucket.org/yasnikoff/shuffler/actions"
)

func TestShellActionWorkingDir(t *testing.T) {
	f := tests.NewDirFx(t)
	defer f.Cleanup()

	cmdName := "bash"
	c := actions.NewCollection()
	a := &Action{ExecPath:cmdName}
	c.Add(a)

	existingDir := f.Path()
	notExistingDir := filepath.Join(existingDir, "no such dir")

	var err error

	a.Dir = notExistingDir
	if _, err = c.RunAction(0); err == nil {
		t.Errorf("Action called in not existing working dir must return error.")
	}

	a.Dir = f.NewRandFile("file.txt")
	if _, err := c.RunAction(0); err == nil {
		t.Errorf("Action called with file as working dir must return error.")
	}

	if _, err := exec.LookPath(cmdName); err != nil {
		t.Skipf("Executable %q is not found in the system's %%PATH%%. Skipping rest of the test.", cmdName)
	}
	a.Dir = existingDir
	if _, err = c.RunAction(0); err != nil {
		t.Errorf("No-op Action called in existing working dir must not return error. Returned rerror: %s", err.Error())
	}

}

func TestShell(t *testing.T) {

	actions.TestAction(&Action{}, t)
	actions.TestAction(&Action{Dir:"some dir"}, t)
	actions.TestAction(&Action{Args:[]string{"Arg1", "Arg2"}}, t)

}
