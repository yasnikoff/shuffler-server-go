package parameters_test

import (
	"testing"
	pt "bitbucket.org/yasnikoff/shuffler/parameters/testing"
	p "bitbucket.org/yasnikoff/shuffler/parameters"
	"strings"
	"github.com/pkg/errors"
	"fmt"
	"encoding/json"
	"reflect"
)

// NormalSignatures returns mapping from test case name to the signature data for that case.
// All the signature data returned is valid and all the tests must pass for all of them.
func NormalSignatures() map[string]sgnData {
	return map[string]sgnData{
		"one param":sgnData{
			{Name:"param1", TypeName:"Type1"},
		},
		"two params":sgnData{
			{Name:"param1", TypeName:"Type1"},
			{Name:"param2", TypeName:"Type2"},
		},
		"space in param name":sgnData{
			{Name:"param 1", TypeName:"Type1"},
		},
		"space in type name":sgnData{
			{Name:"param1", TypeName:"Type 1"},
		},
	}
}

// paramData holds data for parameters used in tests.
type paramData struct {
	Name     string
	TypeName string
}

type sgnData []paramData

// sgnFromData returns new signature populated with given data.
func sgnFromData(d sgnData, r p.TypesRegistry) (*p.Signature, error) {
	sgn := p.NewSignature(r)
	for _, data := range d {
		t, err := r.Get(data.TypeName)
		if err != p.NameRegisteredError {
			t = &pt.StubType{Name:data.TypeName}

			if err := r.Add(data.TypeName, t); err != nil {
				return nil, err
			}
		}
		sgn.AddParameter(data.Name, t, nil)
	}
	return sgn, nil
}


// sgnTest represents single signature test case.
type sgnTest struct {
	Case string
	Sig  p.Signature
	Reg  p.TypesRegistry
	T    *testing.T
}

func newSgnTest(caseName string, sgnData sgnData, t *testing.T) *sgnTest {

	r := p.NewTypesRegistry()

	sgn, err := sgnFromData(sgnData, r)
	if err != nil {
		panic(err)
	}
	return &sgnTest{Case:caseName, Sig:*sgn, Reg:r, T:t}
}

func (s *sgnTest)Errorf(format string, params ...interface{}) {
	s.T.Errorf("case %q: %s", strings.ToUpper(s.Case), fmt.Sprintf(format, params...))
}

func (s *sgnTest)Error(err error) {
	s.Errorf("%s", err)
}

func TestSignature_AddParameter(t *testing.T) {

	t.Parallel()

	CASES:
	for caseName, paramsData := range NormalSignatures() {
		st := newSgnTest(caseName, paramsData, t)
		items := st.Sig.Items()

		expected := len(paramsData)
		result := len(st.Sig.Items())
		if expected != result {
			st.Errorf("Signature has %d items, expect %d", result, expected)
		}

		for i, data := range paramsData {
			tp, err := st.Reg.Get(data.TypeName)
			if err != nil {
				t.Error(err)
				continue CASES

			}
			expected := *p.NewParamSpec(data.Name, tp, tp.NewOptions())
			result := items[i]
			if !reflect.DeepEqual(result, expected) {
				st.Errorf("expect parameter with name %q to be\n%q, got\n%q",
					data.Name, expected, result)
			}
		}
	}

}

func sigsEqual(s1, s2 *p.Signature) (result bool, reason error) {
	items1 := s1.Items()
	items2 := s2.Items()

	if len(items1) != len(items2) {
		return false, fmt.Errorf("param lengths are different: %d and %d", len(items1), len(items2))
	}
	for i, p1 := range items1 {
		p2 := items2[i]
		if !reflect.DeepEqual(p1, p2) {
			return false, fmt.Errorf("Params # %s are different", i)
		}
	}

	if !reflect.DeepEqual(s1.ParamsByName(), s2.ParamsByName()) {
		return false, fmt.Errorf("paramsByName mappings are different:\n%s\n%s", s1.ParamsByName(), s2.ParamsByName())
	}
	return true, nil
}

func TestJSON(t *testing.T) {

	t.Parallel()

	for caseName, paramsData := range NormalSignatures() {
		st := newSgnTest(caseName, paramsData, t)
		s1 := &st.Sig
		data, err := json.Marshal(s1)
		if err != nil {
			st.Error(errors.Wrapf(err, "can't marshal signature %s", s1))
		}
		s2 := p.NewSignature(st.Reg)
		if err := json.Unmarshal(data, s2); err != nil {
			st.Error(errors.Wrapf(err, "can't unmarshal json %q to signature", data))
		}

		if ok, err := sigsEqual(s1, s2); !ok {
			st.Error(err)
		}
	}
}
