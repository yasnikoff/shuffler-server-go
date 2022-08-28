package matchers

import (
	"testing"
	"fmt"
	"strings"
)

type testMatcher struct {
	// matches contains all and only the strings that this matcher will match.
	matches []string
}

func (t *testMatcher)Match(path string) bool {
	for _, m := range t.matches {
		if m == path {
			return true
		}
	}
	return false
}


// All

func TestAll_Match(t *testing.T) {
	m := All
	tests := []string{
		"path",
		"/path",
		"path/",
	}
	for _, path := range tests {
		if !m.Match(path) {
			t.Errorf("All must mutch all paths. It didn't match %q", path)
		}
	}
}


// Compound

type CmpT struct {
	Case  string
	paths [][]string
}

func (ct *CmpT)ExpectedPositives() []string {
	result := []string{}
	for _, paths := range ct.paths {
		result = append(result, paths...)
	}
	return result
}

func (ct *CmpT)Matchers() []Matcher {
	result := make([]Matcher, len(ct.paths))
	for i, pathList := range ct.paths {
		result[i] = &testMatcher{matches:pathList}
	}
	return result
}

// Msg returns test case message for path that didn't match
func (ct *CmpT)Msg(path string) string {
	return fmt.Sprintf("\nCase %s: path %q was not matched.\n", strings.ToUpper(ct.Case), path)
}

func (ct *CmpT)Run(m Matcher, t *testing.T) {
	for _, p := range ct.ExpectedPositives() {
		if !m.Match(p) {
			t.Errorf(ct.Msg(p))
		}
	}
}

var compoundTests []CmpT = []CmpT{
	{
		Case: "no matchers",
		paths:[][]string{},
	},
	{
		Case: "one matcher",
		paths:[][]string{{"path1"}},
	},
	{
		Case: "two matchers",
		paths:[][]string{{"path1"}, {"path2"}},
	},
	{
		Case: "three matchers",
		paths:[][]string{{"path1"}, {"path2"}, {"path3"}},
	},
}

func TestCompound(t *testing.T) {

	for _, test := range compoundTests {
		m := NewCompound(test.Matchers()...)
		test.Run(m, t)

	}
}

func TestCompound_Append(t *testing.T) {

	for _, test := range compoundTests {
		m := NewCompound()
		m.Append(test.Matchers()...)
		test.Run(m, t)
	}
}

func TestNotMatcher(t *testing.T) {
	tests := map[string]struct {
		M        Matcher
		Expected map[string]bool
	}{
		"test matcher":{
			M:&testMatcher{matches:[]string{"one", "two"}},
			Expected:map[string]bool{"one":false, "two":false, "three":true},
		},
		"All matcher":{
			M:All,
			Expected:map[string]bool{"one":false, "two":false, "three":false},
		},
	}

	for caseName, test := range tests {

		n := Not(test.M)
		for path, expected := range test.Expected {
			if n.Match(path) != expected {
				t.Errorf("'Not' matcher returned wrong result for case %s on %q string. Got:%v, Want:%v",strings.ToUpper(caseName),path, !expected, expected)
			}
		}
	}
}
