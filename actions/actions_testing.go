package actions

import (
	"testing"
	"encoding/json"
	"reflect"
)


const mockTypeName = "mockAction"

// Type ActionTest is designed to be used in development of Action implementations.
// It runs several tests that every implementation must pass.
type actionTest struct {
	A Action
	T *testing.T
}

// Run method runs all the tests tests on t.A instance.
func (test *actionTest)Run() {
	test.TestJSONRoundtrip()
}

// TestJSONRoundtrip tests json encoding-decoding.
func (test *actionTest)TestJSONRoundtrip() {
	t := test.T
	data, err := json.Marshal(test.A)
	if err != nil {
		t.Fatal(err)
	}

	b := types[test.A.Type()]()
	if err := json.Unmarshal(data, b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(test.A, b) {
		t.Errorf("json marshal-unmarshal roundtrip returned different Action.\nWant: %#v\n Got:%#v", test.A, b)
	}
}

type mockAction struct {
	commands  []string
	err       error
	output []byte
	callCount int
}

func newMockAction(cmd ...string) *mockAction {
	return &mockAction{commands:cmd}
}

func (a *mockAction)Type() string {
	return mockTypeName
}
func (a *mockAction)Run() ([]byte, error) {
	a.callCount++
	if a.err == error(nil) {
		return a.output, nil
	} else {
		return a.output, a.err
	}
}

func (a *mockAction)MarshalJSON() ([]byte, error) {
	return json.Marshal(a.commands)
}

func (a *mockAction)UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &a.commands)
}

// TestAction runs some tests on Action "a" that are common for all actions.
func TestAction(a Action, t *testing.T){
	test:=&actionTest{A:a,T:t}
	test.Run()
}
