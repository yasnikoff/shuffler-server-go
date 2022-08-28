package shuffler

import (
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"bitbucket.org/yasnikoff/shuffler/parameters"
	"fmt"
)


const DepTypeDepDesc string = "descriptor"

type DepDesc struct {
	IDData     id.ID `json:"id"`
	Parameters *parameters.Collection `json:"parameters,omitempty"`
}

func (desc *DepDesc)ID() id.ID {
	return desc.IDData
}

func(desc *DepDesc)Type()string{
	return DepTypeDepDesc
}

func(desc *DepDesc)Descriptor()*DepDesc{
	return desc
}

func (desc *DepDesc)SetID(ID id.ID) {
	desc.IDData = ID
}

func (desc *DepDesc)String() string {
	return fmt.Sprintf("%s(%s)", desc.IDData, desc.Parameters)
}
