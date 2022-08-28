package tagging

import (
	"testing"
	"reflect"
	"encoding/json"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

func TestAddTags(t *testing.T) {
	c := new(Collection)
	tags := []string{"tag2", "tag1"}
	expectedTags := []string{"tag1", "tag2"} //  sorted
	c.Add(tags...)
	if !reflect.DeepEqual(expectedTags, c.Tags()) {
		t.Errorf("Project.Tags()==%s,\nWant %s\n", c.Tags(), expectedTags)
	}
}

func TestAddDuplicateTag(t *testing.T) {
	//pr := shuffler.NewProject()()
	c := new(Collection)
	dupTag := "duplicate tag"
	tags := []string{"tag1", "tag3", dupTag, dupTag, "tag2", ""}
	expectedTags := []string{dupTag, "tag1", "tag2", "tag3"} // sorted, no blanks
	c.Add(tags...)
	if !reflect.DeepEqual(c.Tags(), expectedTags) {
		t.Errorf("After adding tags with duplicates, project tags are not as expected.\ngot\t\t%q\nwant\t%q", c.Tags(), expectedTags)
	}
}

func TestTagsTwice(t *testing.T) {
	tags1 := []string{"tag1", "tag twice"}
	tags2 := []string{"tag2", "tag twice"}
	expected := []string{"tag twice", "tag1", "tag2"}
	tc := &Collection{}
	tc.Add(tags1...)
	tc.Add(tags2...)
	if !reflect.DeepEqual(tc.Tags(), expected) {
		t.Errorf("After adding tags two times with the same tag Tags()==%q. Want %q", tc.Tags(), expected)
	}

}

func TestRemoveTags(t *testing.T) {

	tests := []struct {
		Case     string
		Initial  []string
		Remove   []string
		Expected []string
	}{
		{Case:"remove the only tag",
			Initial: []string{"tag"},
			Remove:[]string{"tag"},
			Expected:[]string{},
		},
		{Case:"remove first tag from two",
			Initial: []string{"tag1", "tag2"},
			Remove:[]string{"tag1"},
			Expected:[]string{"tag2"},
		},
		{Case:"remove second tag from two",
			Initial: []string{"tag1", "tag2"},
			Remove:[]string{"tag2"},
			Expected:[]string{"tag1"},
		},
		{Case:"remove first tag from three",
			Initial: []string{"tag1", "tag2", "tag3"},
			Remove:[]string{"tag1"},
			Expected:[]string{"tag2", "tag3"},
		},
		{Case:"remove second tag from three",
			Initial: []string{"tag1", "tag2", "tag3"},
			Remove:[]string{"tag2"},
			Expected:[]string{"tag1", "tag3"},
		},
		{Case:"remove first two tags from three",
			Initial: []string{"tag1", "tag2", "tag3"},
			Remove:[]string{"tag1", "tag2"},
			Expected:[]string{"tag3"},
		},
		{Case:"remove last two tags from three",
			Initial: []string{"tag1", "tag2", "tag3"},
			Remove:[]string{"tag2", "tag3"},
			Expected:[]string{"tag1"},
		},
		{Case:"remove all from three",
			Initial: []string{"tag1", "tag2", "tag3"},
			Remove:[]string{"tag1", "tag2", "tag3"},
			Expected:[]string{},
		},
		{Case:"remove first and last from three",
			Initial: []string{"tag1", "tag2", "tag3"},
			Remove:[]string{"tag1", "tag3"},
			Expected:[]string{"tag2"},
		},
		{Case:"remove three tags from two",
			Initial: []string{"tag1", "tag3"},
			Remove:[]string{"tag1", "tag2", "tag3"},
			Expected:[]string{},
		},
		{Case:"remove tags not in collection",
			Initial: []string{"tag1", "tag3"},
			Remove:[]string{"1", "2", "3"},
			Expected:[]string{"tag1", "tag3"},
		},
	}
	for _, test := range tests {
		c := new(Collection)
		c.Add(test.Initial...)
		if !reflect.DeepEqual(c.tags, test.Initial) {
			t.Fatal("tagging.Collection.Remove() test: failed to add initial tags to collection")
		}
		c.Remove(test.Remove...)
		if !reflect.DeepEqual(test.Expected, c.Tags()) {
			t.Errorf("tagging.Collection.Remove(), case %s: after removing %s from %s got %s,Want %s", strings.ToUpper(test.Case), test.Remove, test.Initial, c.Tags(), test.Expected)
		}
	}

}

func TestSet(t *testing.T) {
	names := []string{"1", "2", "3", "4" }
	names1 := names[:2]
	names2 := names[2:]
	tc := &Collection{}
	tc.Set(names1...)
	if !reflect.DeepEqual(names1, tc.Tags()) {
		t.Errorf("After setting tags %q collection's tags are %q, want %[1]q", names1, tc.Tags())
	}
	tc.Set(names2...)
	if !reflect.DeepEqual(names2, tc.Tags()) {
		t.Errorf("After setting tags %q collection's tags are %q, want %[1]q", names2, tc.Tags())
	}

}
/*

func TestTagsSaveLoad(t *testing.T) {
	f := project_fixtures.New(t)
	defer f.Cleanup()

	tags := []string{"tag1", "tag2"}
	f.Project.Tags.Add(tags...)

	if err := f.Project.Save(); err != nil {
		t.Fatal(err)
	}
	pr, err := shuffler.LoadProject(f.Project.Path())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(pr.Tags(), tags) {
		t.Errorf("Loaded project's tags: %q, want %q", pr.Tags(), tags)
	}
}
*/

type TagParseTest struct {
	S string   // Tag String
	T []string //
}

func TestParseRoundtrip(t *testing.T) {
	tests := []TagParseTest{
		{
			S:"tag1;tag2",
			T:[]string{"tag1", "tag2"},
		},
		{
			S:"",
			T:[]string{},
		},
		{
			S:"",
			T:[]string{},
		},
	}
	for _, test := range tests {
		c := new(Collection)
		c.Add(test.T...)
		if c.String() != test.S {
			t.Errorf("Tags collection String()==%q, want %q", c.String(), test.S)
		}
		err := c.Parse(test.S)
		if err != nil {
			t.Error(err)
			continue
		}
		if len(c.Tags()) != len(test.T) {
			t.Errorf("Tags collection Parse() returned %d tags, want %d", len(c.Tags()), len(test.T))
		}
		if !reflect.DeepEqual(c.Tags(), test.T) {
			t.Errorf("Tags collection Parse() returned %q, want %q", c.Tags(), test.T)
		}
	}

}

func TestParse(t *testing.T) {
	tests := []TagParseTest{
		{"", []string{}},
		{";;;;;", []string{}},
		{"tag;;;;", []string{"tag"}},
		{";2;;;1;", []string{"1", "2"}},
		{";2;;;2;", []string{"2"}},
		{";TAG;", []string{"TAG"}},

	}
	for _, test := range tests {
		c := new(Collection)
		err := c.Parse(test.S)
		if err != nil {
			t.Error(err)
			continue

		}
		if !reflect.DeepEqual(c.Tags(), test.T) {
			t.Errorf("Tags collection Parse() returned %q for %q string, want %q", c.Tags(), test.S, test.T)
		}
	}
}

func TestString(t *testing.T) {
	tests := []TagParseTest{
		{T:[]string{"t1", "t2"}, S:"t1;t2"},
		{T:[]string{}, S:""},
		{T:[]string{"3", "2", "1"}, S:"1;2;3"},
	}
	for _, test := range tests {
		c := new(Collection)
		c.Add(test.T...)
		s := c.String()
		if s != test.S {
			t.Errorf("Tags collection String() returned %q for %s, want %q", c.String(), test.T, test.S)
		}
	}

}

func TestTagQueryString(t *testing.T) {
	s := "tag;1;2"
	q := Matcher{}
	if err := q.Parse(s); err != nil {
		t.Fatal(err)
	}
	if q.String() != s {
		t.Errorf("after parsing string %q tagging.Query.String()==%q, want %q", s, q, s)
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		Case   string
		Input  string
		Output []string
	}{
		{
			Case: "duplicate tag",
			Input:`["duplicate","duplicate", "tag 2"]`,
			Output:[]string{"duplicate", "tag 2"},
		},
		{
			Case: "unsorted tags",
			Input:`["z","a"]`,
			Output:[]string{"a", "z"},
		},
		{
			Case: "unsorted tags with duplicates",

			Input:`["z","a"]`,
			Output:[]string{"a", "z"},
		},
	}
	var data []byte
	for _, test := range tests {
		data = []byte(test.Input)
		var tags Collection
		if err := json.Unmarshal(data, &tags); err != nil {
			t.Error(errors.Wrapf(err, "%s: can't unmarshal data %s to tags collection", strings.ToUpper(test.Case), data))
			continue
		}
		result := tags.Tags()
		expected := test.Output
		if !reflect.DeepEqual(result, expected) {
			t.Error(fmt.Errorf("%s: \nInput:\t\t%s\nGot:\t\t%s\nWant:\t\t%s",
				strings.ToUpper(test.Case), test.Input, result, expected))
			continue
		}
	}
}
