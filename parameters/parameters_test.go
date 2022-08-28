package parameters_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/parameters"
	pt "bitbucket.org/yasnikoff/shuffler/parameters/testing"
)

func TestTypesRegistry(t *testing.T) {
	pt.TestRegistry(t, parameters.NewTypesRegistry)

}

