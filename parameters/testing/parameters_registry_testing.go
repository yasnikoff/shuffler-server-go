package parameters_testing

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"reflect"
	"encoding/json"
	"strings"
	"fmt"
)

type RegistryConstructor func() parameters.TypesRegistry // used to create registry instances for tests

// TestRegistry runs some test that every Registry must pass.
// rc is constructor function used to create registry instances for tests.
func TestRegistry(t *testing.T, rc RegistryConstructor) {
	rt := &regTest{t:t, newReg:rc}
	rt.Run()
}

// regTest runs tests for particular registry implementation.
type regTest struct {
	t      *testing.T
	newReg RegistryConstructor // is used to create new registry instance for every test
}

// Type NameTypePair represents single Name-Type pair.
type NameTypePairs map[string]parameters.Type
type TestCasesMap map[string]NameTypePairs // mapping from test case name to a slice of name-type pairs

// NormalTestCases returns TestCasesMap of normal test cases.
// "Normal" means that adding NameTypePairs from the Map is expected to run without any errors.
func NormalTestCases() TestCasesMap {

	return TestCasesMap{
		"one type":
		NameTypePairs{
			"t1": &StubType{"T1"},
		},
		"two types":
		NameTypePairs{
			"t1": &StubType{"T1"},
			"t2": &StubType{"T2"},
		},

	}
}

func (rt *regTest)PopulateReg(r parameters.TypesRegistry, nameTypePairs NameTypePairs) {
	for name, t := range nameTypePairs {
		r.Add(name, t)
	}
}

func (rt *regTest) Run() {
	rt.t.Parallel()
	rt.TestAdd_normal()
	rt.TestAdd_error()
	rt.TestGet()
	rt.testNameOf()
}

// TestAdd_normal adds specified name-type pairs to a new registry.
// All additions a expected to return nil errors.
func (rt *regTest) TestAdd_normal() {
	tests := NormalTestCases()
	for _, test := range tests {
		r := rt.newReg()
		for name, t := range test {
			err := r.Add(name, t)
			if err != nil {
				rt.t.Errorf("Add(%q,%s) returned error:%q. Want nil.", name, t, err, test)
			}
		}
	}
}


// addErrorTestData is a struct for TestAdd_error test cases.
type addErrorTestData struct {
	name string
	Type *StubType
}

func (d *addErrorTestData) String() string {
	return fmt.Sprintf("{%q, %s}", d.name, d.Type.Name)
}


// TestAdd_error adds specified types with specified names to a new registry.
// Adding last name-type pair is expected to return specified error.
func (rt *regTest)TestAdd_error() {
	t1 := &StubType{"T1"}
	t2 := &StubType{"T2"}
	tests := []struct {
		types []addErrorTestData
		err   error // expected error when adding last type
	}{
		{[]addErrorTestData{{"t1", t1}, {"t2", t1}}, parameters.TypeRegisteredError},
		{[]addErrorTestData{{"t1", t1}, {"t1", t2}}, parameters.NameRegisteredError},
	}

	for _, test := range tests {
		r := rt.newReg()
		for ind, td := range test.types {
			err := r.Add(td.name, td.Type)
			if ind + 1 == len(test.types) {
				if err != test.err {
					rt.t.Errorf("Adding last type from %s.\nExpected error returned from Add(%q, %s) is: %q, want: %q", test.types, td.name, td.Type, err, test.err)
				}
			} else {
				if err != nil {
					rt.t.Errorf("Unexpected error while calling Add((%q,%s)).\nIt is a name-type pair number %d  from %#v\n Error: :%s", td.name, td.Type, ind + 1, test.types, err)
				}
			}

		}
	}
}

// testGetData holds test cases for TestGet
type testGetData struct {
	caseName string
	types    map[string]parameters.Type
}

func ( rt *regTest)TestGet() {
	tests := NormalTestCases()
	for caseName, types := range tests {
		r := rt.newReg()
		rt.PopulateReg(r, types)

		for name, t := range types {
			expected := t
			result, err := r.Get(name)
			if err!=nil{
				rt.t.Error(err)
				continue
			}
			if !reflect.DeepEqual(expected, result) {
				rt.t.Errorf("case %q: After addidng name-type pairs %s,\nGet(%q) returned %s, want %s", strings.ToUpper(caseName), types, name, result, expected)
			}
		}
	}

}

// testNameOf tests NameOf method of parameters.TypeRegistry.
func (rt *regTest)testNameOf() {
	tests := NormalTestCases()
	for caseName, types := range tests {
		r := rt.newReg()
		rt.PopulateReg(r, types)
		for name, t := range types {
			expected := name
			result := r.NameOf(t)
			if !reflect.DeepEqual(expected, result) {
				rt.t.Errorf("case %q: After addidng name-type pairs %s,\nNameOf(%q) returned %q, want %q", strings.ToUpper(caseName), types, t, result, expected)
			}

		}
	}
}

func(rt *regTest)testJSON(){
	tests:=NormalTestCases()
	for caseName, types:=range tests{
		r:=rt.newReg()
		rt.PopulateReg(r, types)
		data, err:=json.Marshal(r)
		if err!=nil{
			rt.t.Fatalf("case %q, fatal error while marshalling to json: %s",strings.ToUpper(caseName),err)
		}
		r2:=rt.newReg()
		err=json.Unmarshal(data, r2)
		if err!=nil{
			rt.t.Fatalf("case %q, fatal error while unmarshalling from json: %s",strings.ToUpper(caseName),err)

		}
		if !reflect.DeepEqual(r, r2){
			rt.t.Errorf("case %q, json roundtrip didn't return the same registry.", strings.ToUpper(caseName))
		}
	}

}
