package shell

import (
	"os"
	"fmt"
	"os/exec"
	"bitbucket.org/yasnikoff/shuffler/actions"
)

const ActionType = "shell"

func init() {
	actions.Register(ActionType, func() actions.Action {
		return &Action{}
	})
}

// Shell is a wrapper around os/exec.Cmd that implements Action
type Action struct {
	ExecPath string
	Args     []string
	Dir      string
}

func (sh *Action)Run() ([]byte, error) {
	// TODO: implement timeout
	var err error
	if sh.Dir == "" {
		sh.Dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	info, err := os.Stat(sh.Dir)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("Can not call action with %q as working directory. It is not a directory.", sh.Dir)
	}
	cmd := exec.Command(sh.ExecPath, sh.Args...)
	cmd.Dir = sh.Dir
	return cmd.CombinedOutput()
}

func (sh *Action)Type() string {
	return ActionType
}

