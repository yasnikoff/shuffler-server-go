package actions

import (
	"testing"
	"reflect"
	"bitbucket.org/yasnikoff/shuffler/tests"
	"bytes"
	"fmt"
	"encoding/json"
)
// testJSONFieldIsPresent checks that json data has field with fieldName.
func (ct *collectionTest)testJSONFieldIsPresent(fieldName string, data []byte) {
	t := ct.T
	field := &bytes.Buffer{}
	field.WriteString(fmt.Sprintf("\"%s\":", fieldName))
	if bytes.Index(data, field.Bytes()) == -1 {
		t.Errorf("Format of json encoding of actions.Collection changed. Field %q is not present.", fieldName)
	}
}

func TestCollection_MarshalJSON(t *testing.T) {
	c := NewCollection()

	c.Add(newMockAction("1", "1.1"))
	c.Add(newMockAction("2"))

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}

	test := &collectionTest{T:t}
	test.testJSONFieldIsPresent("type", data)
	test.testJSONFieldIsPresent("data", data)

	c2 := NewCollection()
	if err := json.Unmarshal(data, c2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(c.actions, c2.actions) {
		t.Errorf("json marshal-unmarshal roundtrip returned different actions.Collection.\nWant: %v\nGot: %v", c.actions, c2.actions)
	}
}

func TestRunActionIndex(t *testing.T) {
	c := NewCollection()

	c.Add(newMockAction())
	c.Add(newMockAction())

	testReturnError := func(c *Collection, ind int, mustReturnErr bool) {
		var err error
		_, err = c.RunAction(ind)
		if err == nil && mustReturnErr {
			t.Errorf("RunAction() called with out-of-range index didn't return an error.\nCollection len: %d. Index: %d", len(c.actions), ind)
		}
		if err != nil && !mustReturnErr {
			t.Errorf("RunAction() called with out-of-range index returned an error.\nCollection len: %d. Index: %d", len(c.actions), ind)

		}
	}

	testReturnError(c, -2, true)
	testReturnError(c, -1, true)
	testReturnError(c, 0, false)
	testReturnError(c, 1, false)
	testReturnError(c, 2, true)
	testReturnError(c, 3, true)

}

func TestRunAction(t *testing.T) {
	f := tests.NewDirFx(t)
	defer f.Cleanup()

	a := newMockAction()
	a.output = []byte("action output")

	c := NewCollection()
	c.Add(a)
	out, err := c.RunAction(0)

	if err != nil {
		t.Fatalf("RunAction unexpectedly returned error %q", err.Error())
	}
	if !reflect.DeepEqual(out, a.output) {
		t.Errorf("RunAction retuned wrong action output.\nGot: %v\nWant:%v", out, a.output)
	}

	if a.callCount != 1 {
		t.Errorf("RunAction called Action.Run %d times. Want 1.", a.callCount)
	}
}
