package actions

import (
	"testing"
)

func init() {
	Register(mockTypeName, func() Action {
		return newMockAction()
	})
}

func TestMockAction(t *testing.T) {
	test := &actionTest{A:newMockAction("cmd"), T:t}
	test.Run()
}

type collectionTest struct {
	T *testing.T
}
