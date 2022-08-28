package shuffler

import (
	"reflect"
)

type DepData map[string]Dependency

func (dc *Dependencies)Data() DepData {
	return dc.items
}
func (dc *Dependencies)SetData(data interface{}) error {
	d, ok := data.(DepData)
	if !ok {
		return &SetDataWrongArgType{
			Expected:"DepData",
			Actual:reflect.TypeOf(data).Name()}
	}
	return dc.SetItems(d)
}

