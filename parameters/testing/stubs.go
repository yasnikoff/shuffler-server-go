package parameters_testing

import (
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"fmt"
	"reflect"
	"encoding/json"
)


// testType is parameters.Type implementation for tests
type StubType struct {
	Name string
}

func (t *StubType)NewValue(opt parameters.TypeOptions) (parameters.Value, error) {
	stOpt, ok:=opt.(*StubOptions)
	if !ok{
		return nil, fmt.Errorf("Wrong otions type for StubType: %T", opt)
	}
	return &StubValue{t:t, Val:stOpt.DefaultValue}, nil
}

func (t *StubType)NewOptions()parameters.TypeOptions{
	return &StubOptions{t:t}
}

func (t *StubType)String()string{
	return t.Name
}

// StubOptions is parameters.TypeOptions implementation for tests.
type StubOptions struct{
	t *StubType
	DefaultValue interface{}
}

func (o *StubOptions)Type()parameters.Type{
	return  o.t
}

func (o *StubOptions)Error()error{
	return  nil
}

func (o *StubOptions)String()string{
	return  "{}"
}


// testVal is parameters.Value implementation for tests.
type StubValue struct {
	t   *StubType
	Val interface{}
}

func (v *StubValue)Equals(val parameters.Value) bool {
	var tv *StubValue
	tv, ok := val.(*StubValue)
	if !ok {
		return false
	}
	return reflect.DeepEqual(tv.Val, v.Val)

}

func (v *StubValue)Type() parameters.Type {
	return v.t
}

func (v *StubValue)MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Val)
}

func (v *StubValue)UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.Val)
}

