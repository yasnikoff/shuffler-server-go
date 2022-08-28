package shuffler

import (
	"fmt"
	"encoding/json"
)

var docTypes = make(map[string]*DocType)
type EmptyDocConstructor func() Doc

type DocType struct {
	name        string
	constructor EmptyDocConstructor
}

func (dt *DocType)Name() string {
	return dt.name
}

func (dt *DocType)New() Doc {
	d := dt.constructor()
	h:=d.GetHeader()
	h.setType(dt)
	h.init(d, dt)
	return d
}


func (dt *DocType)MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.Name())
}

func (dt *DocType)UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	t, ok := docTypes[name]
	if !ok {
		return fmt.Errorf("Unknown document type: %s", name)
	}
	*dt = *t
	return nil
}

func RegisterType(name string, constructor EmptyDocConstructor) *DocType {

	if _, ok := docTypes[name]; ok {
		panic(fmt.Errorf("document type %v already registered", name))
	}
	dt := &DocType{name, constructor}
	docTypes[name] = dt
	return dt
}
