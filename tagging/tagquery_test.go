package tagging

import (
	"testing"
)

type TagMatcherTest struct {
	QS       string // query string
	TS       string // tags string
	Expected bool
}

func (test *TagMatcherTest)Run(t *testing.T) {
	m := Matcher{}
	if err := m.Parse(test.QS); err != nil {
		t.Errorf("error parsing string %q: %s", test.QS, err)
	}
	tc := &Collection{}
	if err := tc.Parse(test.TS); err != nil {
		t.Errorf("error parsing string %q: %s", test.TS, err)
	}

	if m.Match(tc) != test.Expected {
		t.Errorf("query %q for tags %q returned %v, want %v", test.QS, test.TS, !test.Expected, test.Expected)
	}
}

func TestTagMatcher(t *testing.T) {
	t.Parallel()
	tests := []TagMatcherTest{
		{"", "", true},
		{"", "tag", true},
		{"1", "1", true},
		{"a", "a", true},
		{"a", "b", false},
		{"1", "1;2", true},
		{"tag", "", false},
		{"1;2", "1;2", true},
		{"1;2", "1;2;3", true},
		{"1;2;3", "1;2", false},
		{"1;2;3", "1;2", false},

		{"1;2;3 || 1", "1;2", true},
		{"1;2;3 || 1;2", "1;2", true},

		{"1;2;3 || 1;2", "", false},
		{"1;2;3 || 1;2 || ", "", true},
		{"1;2;3 || 1;2 || ;", "", true},
		{"1;2;3 || 1;2 ||", "", true},
		{"1;2;3 || 1;2 ||;", "", true},

	}
	for _, test := range tests {
		test.Run(t)

	}
}
