package accessors

import ("time"
	"encoding/json"
)

type Time struct{
	value time.Time
}

func (t Time)Get()time.Time{
	return t.value
}

func (t *Time)Set(value time.Time){
	t.value = value
}
func (t Time)String()string{
	return t.value.String()
}
func (t Time)MarshalJSON()([]byte,error){
	return json.Marshal(t.value)
}
func (t *Time)UnmarshalJSON(data []byte)error{
	return json.Unmarshal(data, &(t.value))
}


