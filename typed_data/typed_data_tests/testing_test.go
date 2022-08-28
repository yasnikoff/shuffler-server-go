// Tests for typed_data_tests itself.
// Makes sure that provided testing functions
// really capture errors in Data and Registry implementations.
package typed_data_tests

import (
	"bitbucket.org/yasnikoff/shuffler/typed_data"
	"testing"
	"github.com/pkg/errors"
)

// Test registerError() func on broken registry
// (Test for false positives)
func TestRegisterErrorOnBrokenRegistry(t *testing.T) {
	r := &brokenRegistry{}
	dt := typed_data.NewDataType(
		"TestDataType",
		func() typed_data.Data {
			return &WorkingData{}
		},
	)
	if err := RegisterCheck(r, dt); errors.Cause(err) != ErrCantGetTypeAfterRegister {
		t.Errorf("testRegister() didn't detect disfunctional Register()-Get() funcions on registry %q", r)
	}
}


// Return working registry implementation with one data type.
func DefaulTestRegistry()typed_data.Registry{
	r := &WorkingRegistry{}
	typeName:="test type"
	dt := typed_data.NewDataType(
		"TestDataType",
		func() typed_data.Data {
			return &WorkingData{}
		})
	if err:=r.Register(typeName, dt);err!=nil{
		panic(err)
	}
	return r
}

// Test registerError() func on working registry.
// (Test for  false negatives)
func TestRegisterErrorOnWorkingRegistry(t *testing.T) {
	r := &WorkingRegistry{}
	dt := typed_data.NewDataType(
		"TestDataType",
		func() typed_data.Data {
			return &WorkingData{}
		})
	if err := RegisterCheck(r, dt); err != nil {
		t.Errorf(
			"testRegister() detected error in Register() - Get() functions of working Registry %q",
			r)
	}
}

//Test State()->SetState() on broken Data.
// (Test for false positives)
func TestSetDataCheckBroken(t *testing.T){
	b:=&BrokenData{5}
	if err:= SetStateCheck(b);err==nil{
		t.Errorf("SetStateError(b) on broken Data implementation didn't return error")
	}
}
//Test State()->SetState() on working Data.
// (Test for false negatives)
func TestSetDataCheckWorking(t *testing.T) {
	d:=&WorkingData{"some string", 66}
	if err:= SetStateCheck(d);err!=nil{
		t.Errorf("SetStateError(b) on working Data implementation returns error")
	}
}

//Test JSON Marshal->Unmarshal on broken Data.
// (Test for false positives)
func TestJSONCheckBroken(t *testing.T){
	b:=&BrokenData{5}
	if err:= JSONCheck(b);err==nil{
		t.Errorf("SetStateError(b) on broken Data implementation didn't return error")
	}
}
//Test JSON Marshal->Unmarshal on working Data.
// (Test for false negatives)
func TestJSONCheckWorking(t *testing.T) {
	d:=&WorkingData{"some string", 66}
	if err:= JSONCheck(d);err!=nil{
		t.Errorf("SetStateError(b) on working Data implementation returns error")
	}
}
