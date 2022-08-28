package typed_data_test

import (
	"testing"
	td "bitbucket.org/yasnikoff/shuffler/typed_data"
	tests "bitbucket.org/yasnikoff/shuffler/typed_data/typed_data_tests"
)

func TestRegistry(t *testing.T) {
	tests.TestRegistry(t,
		func() td.Registry {
			return td.NewDefaultRegistry("Test Registry")
		},
		func() td.Type {
			return td.NewDataType("Test Data Type", func()td.Data{
				return &tests.WorkingData{}
			})
		},
	)
}

//func TestDataType(t *testing.T) {
//
//	var name:="test"
//	data:=struct{
//		t:
//	}
//	d:=td.NewDataType(name, )
//}
