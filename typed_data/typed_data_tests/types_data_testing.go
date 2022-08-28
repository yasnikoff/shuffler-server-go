// Provides functions for testing Data and Registry implementations.
package typed_data_tests

import (
	"testing"
	td "bitbucket.org/yasnikoff/shuffler/typed_data"
	"reflect"
	"github.com/pkg/errors"
	"encoding/json"
)


// Registry tests

var ErrCantGetTypeAfterRegister = errors.New("can't Get type after Register")

// Run tests that every Regestry implementation must pass.
// regConstructor() - function that will be called to get new Registry instance for every test.
// dtConstructor() - function that will be called to get new DataType instance for every test.
func TestRegistry(t *testing.T, regConstructor func() td.Registry, dtConstructor func() td.Type) {
	t.Run("Register Type", func(t *testing.T) {
		if err := RegisterCheck(regConstructor(), dtConstructor()); err != nil {
			t.Errorf(err.Error())
		}
	})
}

// Return error indicating that Register() - Get() functions don't work as expected.
func RegisterCheck(r td.Registry, dt td.Type) error {
	const errMsg string = "error registering type %s in registry %s"

	name:= "Test Type"
	r.Register(name, dt)

	dt2, err:=r.Get(name)
	if err!=nil{
		return errors.Wrapf(err,
			errMsg, dt, r)
	}
	if !reflect.DeepEqual(dt2, dt) {
		return errors.Wrapf(ErrCantGetTypeAfterRegister,
			errMsg, dt, r)
	}
	return nil
}


// Data tests

func TestData(t *testing.T, constructor func()td.Data){
	t.Run("Data()->SetData()", func(t *testing.T){
		if err:=SetStateCheck(constructor());err!=nil{
			t.Errorf(err.Error())
		}
	})

	t.Run("JSON Marshal->Unmarshal", func(t *testing.T){
		if err:=JSONCheck(constructor());err!=nil{
			t.Errorf(err.Error())
		}
	})
}


func SetStateCheck(d td.Data) error {
	stateData := d.State()
	newData := d.Type().New()
	if err := newData.SetState(stateData); err != nil {
		return errors.Wrap(err, "can't SetState(data) that was obtained by Data()")
	}
	if !reflect.DeepEqual(newData, d) {
		return errors.New("after d1.SetState(d.Data()) d1!=d")
	}
	return nil
}

func JSONCheck(d td.Data) error {
	stateData := d.State()
	bData, err := json.Marshal(stateData)
	if err != nil {
		return errors.Wrapf(err, "can't marshal state data of %s to JSON", d)
	}
	newStateData := d.Type().New().State()
	if err := json.Unmarshal(bData, newStateData); err != nil {
		return errors.Wrapf(err, "can't unmarshal JSON %s ", bData)
	}
	if !reflect.DeepEqual(stateData, newStateData) {
		return errors.Wrapf(err,
			"unmarshaled state data %q does not equal original state data %q",
			newStateData, stateData)
	}
	return nil
}
