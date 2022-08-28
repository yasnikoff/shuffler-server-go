package typed_data

import (
	"github.com/pkg/errors"
	"encoding/json"
	"fmt"
)

var ErrNoSuchTypeRegistered error = errors.New("type is not registered in the registry")

// Instances of the Data interface can provide data (through the State() method call)
// that represents their state.
// They are also can restore their state from that data.
// The data is expected to be serializable by json package.
type Data interface {
	// Data that is serializable and that can be used later by SetState to fully restore
	// object's state.
	State() interface{}
	// Instance is expected to be able to fully restore its state based on the data.
	SetState(state interface{}) error
	// Type of the Data
	Type() Type
}

type Registry interface {
	String() string
	// Register DataType under the given name.
	Register(name string, t Type) error
	// Return DataType registered under the given name.
	Get(name string) (Type, error)
	// Return the name of the given Type
	GetTypeName(Type) (name string, ok bool)
}

type Type interface {
	// Return new default Data of this Type
	New() Data
}


// Default implementation of the Registry interface.
func NewDefaultRegistry(name string) Registry {
	return &defaultTypeRegistry{
		items:make(map[string]Type),
		names:make(map[Type]string),
	}
}

// Simple implementation of Type interface.
func NewDataType(name string, constructor func() Data) Type {
	return &dataType{name:name, constructor:constructor}
}


// Simple Type implementation
type dataType struct {
	// The name of the data type.
	name        string
	// Returns new instance of data ready to call SetState() on.
	constructor func() Data
}

func (dt *dataType)Name() string {
	return dt.name
}

func (dt *dataType)New() Data {
	return dt.constructor()
}


// Default implementation of TypeRegistry
type defaultTypeRegistry struct {
	items map[string]Type
	names map[Type]string
	Name  string
}

func (r *defaultTypeRegistry)String() string {
	return r.Name
}

func (r *defaultTypeRegistry)Register(name string, dt Type) error {
	if _, ok := r.items[name]; ok {
		//return errors.Wrapf(ErrTypeAlreadyRegistered, "DataType name: %q", name)
		return TypeAlreadyRegisteredError{RegistryError{r}, dt}
	}
	r.items[name] = dt
	r.names[dt] = name
	return nil
}

func (r *defaultTypeRegistry)Get(name string) (Type, error) {
	t, ok := r.items[name]
	if !ok {
		return nil, errors.Wrapf(ErrNoSuchTypeRegistered, "type: %q, registry: %s", t, r)
	}
	return t, nil
}

func (r *defaultTypeRegistry)GetTypeName(dt Type) (string, bool) {
	name, ok := r.names[dt]
	return name, ok
}

type Marshaler interface {
	Unmarshal([]byte) (Data, error)
	Marshal(Data) ([]byte, error)
}

type defaultJSONMarshaler struct {
	r Registry
}

// Return Marshaler implementation for JSON format.
func NewJSONMarshaler(r Registry) Marshaler {
	return &defaultJSONMarshaler{r:r}
}

func (m *defaultJSONMarshaler)String() string {
	return fmt.Sprintf("json Marshaller for Registry %q", m.r)
}

func (m *defaultJSONMarshaler)Marshal(d Data) ([]byte, error) {

	errMsg := "can't marshal data %s"

	wrapErr := func(err error) error {
		msg := fmt.Sprintf(errMsg, d)
		return errors.Wrapf(err, msg)
	}

	typeName, ok := m.r.GetTypeName(d.Type())
	if !ok {
		return nil, wrapErr(TypeNotRegisteredError{
			RegistryError:RegistryError{Reg:m.r},
			Type:d.Type()},
		)
	}

	var state = struct {
		t string `json:"type"`
		d interface{} `json:"data"`
	}{
		t: typeName,
		d:d.State(),
	}

	return json.Marshal(&state)
}

func (m *defaultJSONMarshaler)Unmarshal(data []byte) (Data, error) {
	state := struct {
		t  string `json:"type"`
		rm json.RawMessage `json:"data"`
	}{}

	const errMsg string = "can't unmarshal data %s"

	err := json.Unmarshal(data, &state)
	if err != nil {
		return nil, errors.Wrapf(err, errMsg, data)
	}

	t, err := m.r.Get(state.t)
	if err != nil {
		return nil, errors.Wrapf(err, errMsg, data)
	}

	d := t.New()
	if err := json.Unmarshal(state.rm, d); err != nil {
		return nil, errors.Wrapf(err, errMsg, data)
	}

	return d, nil
}
