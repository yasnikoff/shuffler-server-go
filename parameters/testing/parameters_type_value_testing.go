package parameters_testing

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/parameters"
)

type TypeTest struct{
	*testing.T
	Type *parameters.Type
}

func (t *TypeTest)Run(){
	t.Parallel()
}


