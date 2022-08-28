package parameters

import (
	"fmt"
	"github.com/pkg/errors"
	"encoding/json"
	"strings"
)

var emptyReg TypesRegistry
var emptyRegError = errors.New("default registry is designed to be always empty")

func init() {
	emptyReg = &emptyRegistry{}
}

// emptyRegistry serves as a no-op stub for TypeRegistry.
type emptyRegistry struct{}

func (r *emptyRegistry)Add(name string, t Type) error {
	return emptyRegError
}

func (r *emptyRegistry)Get(name string) (Type, error) {
	return nil, emptyRegError
}

func (r *emptyRegistry)NameOf(t Type) string {
	return ""
}

func (r *emptyRegistry)Items() map[string]Type {
	return nil
}

type TypeOptions interface {
	// Type returns type that it is applied to.
	Type() Type

	// Returns nil if the options can be used to produce value.
	// Otherwise returns errors
	Error() error

	String()string
}

// paramSpec holds information about one parameter in a signature.
type paramSpec struct {
	Name    string
	Type    Type
	Options TypeOptions
}

func (ps *paramSpec)NewValue() Value {
	val, err := ps.Type.NewValue(ps.Options)
	if err != nil {
		panic(err) // errors must be checked in AddParameter
	}
	return val
}

func (ps *paramSpec)String()string{
	return fmt.Sprintf("", ps.Name, ps.Type, ps.Options )
}

//type Signature *signature

// Signature represents mapping from parameters' names to their types and value constraints.
// Use NewSignature to create instances of the type.
type Signature struct {
	items        []paramSpec
	paramsByName map[string]*paramSpec
	reg          TypesRegistry
}

// NewSignature returns new empty signature.
func NewSignature(r TypesRegistry) *Signature {
	if r == nil {
		r = emptyReg
	}
	sgn := &Signature{
		paramsByName:make(map[string]*paramSpec),
		reg:r,
	}
	return sgn
}

func (sgn *Signature)AddParameter(name string, t Type, opt TypeOptions) error {
	if opt == nil {
		opt = t.NewOptions()
	}
	if err := opt.Error(); err != nil {
		return errors.Wrapf(err, "options for parameter %q are not valid", name)
	}

	if _, ok := sgn.paramsByName[name]; ok {
		return fmt.Errorf("Parameter %q is already in the signature", name)
	}
	ps := paramSpec{Name:name, Type:t, Options:opt}
	sgn.items = append(sgn.items, ps)
	sgn.paramsByName[name] = &ps

	return nil
}

func (sgn *Signature)SetTypesRegistry(r TypesRegistry) {
	sgn.reg = r
}

// Items returns parameter specifications slice
func (sgn *Signature)Items() []paramSpec {
	return sgn.items
}

// typeData is used in Type json marshalling
type paramSpecData struct {
	Name     string `json:"name"`
	TypeName string `json:"type"`
	Options  TypeOptions `json:"options"`
}

func (sgn *Signature)String() string {
	parts := make([]string, len(sgn.items))
	for i, t := range sgn.items {
		parts[i] = fmt.Sprintf("(%s)", t)
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

func (sgn *Signature)MarshalJSON() ([]byte, error) {
	sgnData := make([]paramSpecData, len(sgn.items))
	for i, t := range sgn.items {
		sgnData[i] = paramSpecData{Name: t.Name, TypeName:sgn.reg.NameOf(t.Type), Options:t.Options}
	}
	return json.Marshal(sgnData)
}

func (sgn *Signature)UnmarshalJSON(data []byte) error {

	var paramsData []struct {
		Name    string `json:"name"`
		Type    string`json:"type"`
		Options json.RawMessage `json:"options"`
	}

	if err := json.Unmarshal(data, &paramsData); err != nil {
		return errors.Wrapf(err, "can't unmarshal json data %s", data)
	}

	sgn.items = make([]paramSpec,0, len(paramsData))

	for _, typeData := range paramsData {
		t, err := sgn.reg.Get(typeData.Type)
		if err != nil {
			return errors.Wrapf(err, "can't unmarshal json data %s", data)
		}
		opt := t.NewOptions()
		if err := json.Unmarshal(typeData.Options, opt); err != nil {
			return errors.Wrapf(err, "can't unmarshal type options data data %s", typeData.Options)

		}
		sgn.AddParameter(typeData.Name, t, opt)
	}
	return nil
}

// NewCollection returns new parameters collection of the given signature.
func (sgn *Signature)NewCollection() *Collection {
	result := NewCollection(sgn)
	for _, parSpec := range sgn.items {
		result.items[parSpec.Name] = parSpec.NewValue()
	}
	return result
}
