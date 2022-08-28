package shuffler

import (
	"fmt"
)

type Data interface {
	Data() interface{}
	SetData(interface{}) error
}

const setDataWrongArgTypeMsg string = "Wrong argument type for SetData method. Expected: %s got %s"

type SetDataWrongArgType struct {
	Expected string
	Actual   string
}

func (err *SetDataWrongArgType)Error() string {
	return fmt.Sprintf(setDataWrongArgTypeMsg, err.Expected, err.Actual)
}

