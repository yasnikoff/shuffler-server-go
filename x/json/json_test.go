package json

import (
	"encoding/json"
	"testing"
)

type Child struct{
	name string
	value int
}

func (ch *Child)MarshalJSON()([]byte, error){
	return json.Marshal(map[string]interface{}{"name":ch.name, "value":ch.value})
}

type Parent struct{
	Name string
	child *Child
}

func(p *Parent)MarshalJSON()([]byte, error){
	return json.Marshal(ParentData{Name:p.Name, Child:p.child})
}

type ParentData struct{
	Name string `json:"name"`
	Child *Child `json:"child"`
}

func testMarshal(t *testing.T){
	child := &Child{name:"child 1", value:12}
	parent:=&Parent{child:child}

	data, err:=json.Marshal(parent)
	if err!=nil{
		t.Errorf("Error: %s", err)
	}
	t.Logf("%s", data)
}

