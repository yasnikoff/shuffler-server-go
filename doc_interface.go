package shuffler

import (
	// "bitbucket.org/yasnikoff/shuffler/types/id"
)

type DocData interface{

}

type Doc interface {
	//ID() id.ID
	//SetID(id.ID)

	Path()string
	SetPath(string) error

	GetHeader()*Header
	//
	//Type() DocType
	//SetType(DocType)

	//SavedTime() time.Time
	//SetSavedTime(time.Time)

	//Data() DocData
	//SetState()DocData)error

	Save()error
	SaveTo(string)error

	//json.Marshaler
	//json.Unmarshaler
}
