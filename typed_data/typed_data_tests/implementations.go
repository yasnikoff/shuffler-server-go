package typed_data_tests

import (
	"github.com/pkg/errors"

	td "bitbucket.org/yasnikoff/shuffler/typed_data"
	"encoding/json"
)

var testDataType td.Type = td.NewDataType(
	"TestDataType",
	func() td.Data {
		return &WorkingData{}
	})

type StateData struct {
	String string
	Int    int
}

// Dysfunctional Data implementation
type BrokenData struct{
	state int
}

func(d *BrokenData)State()interface{}{
	return &BrokenStateData{d.state+1}
}

func(d *BrokenData)SetState(s interface{})error{
	return nil
}
func(d *BrokenData)Type()td.Type{
	return testDataType
}

// Dysfunctional json.Marshaler - json.Unmarshaler
type BrokenStateData struct{
	state int
}
func( d *BrokenStateData)MarshalJSON()([]byte, error){
	return json.Marshal(d.state+1)
}
func (d *BrokenStateData)UnmarshalJSON(data []byte)error{
	var state int
	if err:=json.Unmarshal(data, &state);err!=nil{
		return err
	}
	d.state = state +1
	return nil
}

type WorkingData struct {
	stringField string
	intField    int
}

func (td *WorkingData)State() interface{} {
	return &StateData{
		td.stringField,
		td.intField,
	}
}

func (td *WorkingData)SetState(data interface{}) error {
	d, ok := data.(*StateData)
	if !ok {
		return errors.Errorf("data %s is not *StateData", data)
	}
	td.intField = d.Int
	td.stringField = d.String

	return nil
}

func (td *WorkingData)Type() td.Type {
	return testDataType
}

// Registry that does not register types.
type brokenRegistry struct{}

func (r *brokenRegistry)String() string {
	return "Broken Registry"
}
func (r *brokenRegistry)Register(name string, dt td.Type) error {
	return nil
}

func (r *brokenRegistry)Get(name string) (td.Type, error) {
	return nil, nil
}

func (r *brokenRegistry)GetTypeName(t td.Type) (string, bool) {
	return "", false
}


// Simple Registry implementation.
type WorkingRegistry struct {
	td.Registry
	dataTypeName string
	dataType td.Type // for simplicity only one item can be stored
}

func (r *WorkingRegistry)String() string {
	return "Working Registry"
}

func (r *WorkingRegistry)Register(name string, dt td.Type) error {
	if r.dataType != nil {
		return errors.New("This Registry implementations only supports calling Register() once")
	}
	r.dataType = dt
	r.dataTypeName = name
	return nil
}
func (r *WorkingRegistry)Get(name string) (td.Type, error) {
	if name == r.dataTypeName{
		return r.dataType, nil
	}
	return nil, errors.Errorf("no type with the name %q in %s", name, r)
}


