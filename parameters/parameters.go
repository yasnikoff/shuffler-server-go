package parameters

import (
	"encoding/json"
	"github.com/pkg/errors"
	"fmt"
)

var NameRegisteredError = errors.New("name is already in the registry")
var TypeRegisteredError = errors.New("type is already in the registry")

type Value interface {
	json.Marshaler
	json.Unmarshaler
	Equals(Value) bool
	Type() Type
}

// Type is used to create new Values.
type Type interface {
	// NewValue returns new default value for the type with given options.
	NewValue(opt TypeOptions) (Value, error)

	// NewOptions returns options for the type initialised to default value(s).
	NewOptions()TypeOptions

	// Name returns the type's string representation
	String()string

}

type TypesRegistry interface {
	// Add registers Type t under the given name.
	// If the name or the Type is already in the registry, en error is returned.
	// Registry may check the Type for some criteria and return en error.
	Add(name string, t Type) error

	// Get returns Type with given name or nil if nothing is found.
	Get(name string) (Type, error)

	// NameOf returns name of the type in the registry
	NameOf(t Type) string

	// Remove deletes Type from the registry and returns it.
	// If there is no Type registered under given name, nil is returned.
	//Remove(name string) Type

	// Items returns name to Type map of all types in the registry.
	Items() map[string]Type
}



// NewTypesRegistry creates new types registry.
func NewTypesRegistry() TypesRegistry {
	return &typesRegistry{
		types:make(map[string]Type),
		names:make(map[Type]string),
	}
}

type typesRegistry struct {
	types map[string]Type
	names map[Type]string
}

func (r *typesRegistry)Add(name string, t Type) error {
	if _, ok := r.names[t]; ok {
		return TypeRegisteredError
	}
	if _, ok := r.types[name]; ok {
		return NameRegisteredError
	}
	r.names[t] = name
	r.types[name] = t
	return nil
}

func (r *typesRegistry)Get(name string) (Type, error) {
	t, ok := r.types[name]
	if !ok {

		return nil, fmt.Errorf("no type with the name %q in the registry", name)
	}
	return t, nil
}

func (r *typesRegistry)NameOf(t Type) string {
	n := r.names[t]
	return n
}
/*

func (r *typesRegistry)Remove(name string) Type {
	return nil
}
*/

func (r *typesRegistry)Items() map[string]Type {
	return r.types
}

// valueData is used for unmarshaling json to Value
type valueData struct {
	typeName string
	data     json.RawMessage
}

// Collection holds values of parameters and a reference to Signature.
//type Collection *collection

type Collection struct {
	Sig   *Signature
	items map[string]Value
}

func (c *Collection)MarshalJSON() ([]byte, error) {
	itemsData := make(map[string]*valueData, len(c.items))
	for name, value := range c.items {
		itemsData[name] = &valueData{typeName:c.Sig.reg.NameOf(value.Type())}
	}
	return json.Marshal(itemsData)
}

func NewCollection(s *Signature) *Collection {
	if s == nil {
		s = NewSignature(nil)
	}
	return &Collection{
		Sig: s,
		items:make(map[string]Value),
	}
}

