package accessors

import (
	"encoding/json"
)

type String struct{
	value string
}

func (t String)Get()string{
	return t.value
}

func (t *String)Set(value string){
	t.value = value
}
func (t String)String()string{
	return t.value
}
func (t String)MarshalJSON()([]byte,error){
	return json.Marshal(t.value)
}
func (t *String)UnmarshalJSON(data []byte)error{
	return json.Unmarshal(data, &(t.value))
}

