package blender

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/actions"
	"bitbucket.org/yasnikoff/shuffler/actions/shell"
	"os/exec"
	"regexp"
	"fmt"
	"path/filepath"
	"bufio"
	"bytes"
)

var blendInPath bool

const blendNotInPathMsg = "Executable \"blender\" is not in the system's %%PATH%%. Skipping."

func init() {
	_, err := exec.LookPath("blender")
	if err == nil {
		blendInPath = true
	}
}

func TestAction(t *testing.T) {
	actions.TestAction(&Action{}, t)
	actions.TestAction(&Action{BlendFile:"main.blend", Py:"script.py", PyArgs:[]string{"1", "2"}}, t)
}

type BlendTest struct {
	T   *testing.T

	// blender.Action to test
	A   *Action

	// blender CombinedOutput returned from Run()
	Out []byte

	// error returned from Run()
	Err error
}

func (bt *BlendTest)RunAction() {
	if !blendInPath {
		bt.T.Skipf(blendNotInPathMsg)
	}

	bt.Out, bt.Err = bt.A.Run()
	if bt.Err != nil {
		bt.T.Fatal(bt.Err)
	}
	if bt.Match(fmt.Sprintf("%s': No such file or directory", bt.A.BlendFile)) {
		bt.T.Fatal("Blender file %q not found.", bt.A.BlendFile)
	}
	if bt.Match(fmt.Sprintf(`%s" could not be opened: No such file or directory`, bt.A.Py)) {
		bt.T.Fatal("Script file %q is not found", bt.A.Py)
	}
}

func (bt *BlendTest)Match(pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.Match(bt.Out)
}

func (bt *BlendTest)FindLine(line string)bool{
	scanner:=bufio.NewScanner(bytes.NewReader(bt.Out))
	for scanner.Scan(){
		if scanner.Text()==line{
			return true
		}
	}
	return false
}

func (bt *BlendTest)PrintOutput() {
	bt.T.Log(string(bt.Out))
}

func TestScriptArgs(t *testing.T) {

	t.Parallel()

	a := &Action{
		Action:shell.Action{
			Dir:"testdata",
		},
		//Blend:"main.blend",
		Py:"test.py",
		PyArgs:[]string{
			"Arg1",
			"Arg2",
		},
	}

	test := &BlendTest{T:t, A:a}
	test.RunAction()

	//test.PrintOutput()

	//verify that sys.argv[0] is set to absolute path to the script
	expectedPath, err:=filepath.Abs(a.Dir)
	if err!=nil{
		t.Fatal(err)
	}
	expectedPath=filepath.Join(expectedPath, a.Py)
	expectedPath=`TEST SCRIPT: ` + expectedPath
	if !test.FindLine(expectedPath){
		t.Errorf("Test script didnt print this expected line with the path to the script:\n%s", expectedPath)
	}

	// verify that only script args are left in the sys.argv array
	expectedArgs := `TEST ARGS: ['Arg1', 'Arg2']`
	if !test.FindLine(expectedArgs) {
		t.Errorf("Test script didn't print expected line with the script args:\n%s\n", expectedArgs)
	}

}

// TODO: test with and without .blend file
