package typed_data

import (
	"fmt"
)

// RegistryError

const registryErrorMsg string = "registry %q error"

type RegistryError struct {
	Reg Registry
}

func (err *RegistryError)Error() string {
	return fmt.Sprintf(registryErrorMsg, err.Reg)
}


// ErrNoTypeWithName

const errNoTypeWithNameMsg string = "no type with the name %q"

type ErrNoTypeWithName struct {
	RegistryError
	Name string
}

func (err ErrNoTypeWithName)Error() string {
	return fmt.Sprintf("%s: %s", err.RegistryError.Error(), fmt.Sprintf(errNoTypeWithNameMsg, err.Name))
}


// TypeNotRegisteredError

const typeNotRegisteredErrorMsg string = "type %q is not registered"

type TypeNotRegisteredError struct {
	RegistryError
	Type Type
}

func (err TypeNotRegisteredError)Error() string {
	return fmt.Sprintf("%s: %s", err.RegistryError.Error(), fmt.Sprintf(typeNotRegisteredErrorMsg, err.Type))
}

const typeAlreadyRegisteredErrorMsg string = "type %q is already registered"

type TypeAlreadyRegisteredError struct {
	RegistryError
	Type Type
}

func (err TypeAlreadyRegisteredError)Error() string {
	return fmt.Sprintf("%s: %s", err.RegistryError.Error(), fmt.Sprintf(typeAlreadyRegisteredErrorMsg, err.Type))
}
