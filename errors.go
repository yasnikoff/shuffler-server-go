package shuffler

import (
	"fmt"
)

type WrongDocTypeError struct {
	Doc          interface{}
	ExpectedType string
}

func (err *WrongDocTypeError)Error() string {
	return fmt.Sprintf("wrong doc type of %s; expected type: %s, got %T", err.Doc, err.ExpectedType)
}

func IsWrongDocType(err error) bool {
	switch err.(type){
	case *WrongDocTypeError:
		return true
	}
	return false
}
