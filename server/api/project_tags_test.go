package api

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler"
	"bitbucket.org/yasnikoff/shuffler/tagging"
	"fmt"
)

type TagsTest struct {
	*apiTestFx
	Project *shuffler.Project
	Tags    *tagging.Collection
}

func ( t *TagsTest)subPath(subPath ...string) string {
	args := []string{"projects", t.Project.ID().String(), "tags"}
	if len(subPath) > 0 {
		args = append(args, subPath...)
	}
	return toPath(args...)
}

func NewTagsTest(t *testing.T) *TagsTest {
	r := &TagsTest{
		apiTestFx:NewRunner(t, "test"),
	}
	r.Project = r.AppFixture.NewProject()
	r.Tags = &r.Project.Tags
	return r
}

func TestTagsGet(t *testing.T) {

	r := NewTagsTest(t)
	defer r.Cleanup()

	r.CaseName = "get two tags"
	tags := []string{"tag1", "tag2"}
	r.Tags.Add(tags...)

	if err := r.Project.Save(); err != nil {
		t.Fatal(err)
	}

	r.Expected = tags
	r.Get(r.subPath())
	r.Check()
}

func TestTagsAdd(t *testing.T) {
	r := NewTagsTest(t)
	defer r.Cleanup()

	caseName := func(initial, added string) string {
		return fmt.Sprintf("add %s to %s",added, initial )
	}
	tests := []struct {
		Case        string
		InitialTags []string
		NewTags     []string
	}{
		{caseName("Empty", "Empty"), []string{}, []string{}},
		{caseName("Empty", "one tag"), []string{}, []string{"tag"}},
		{caseName("Empty", "two tags"), []string{}, []string{"tag1", "tag2"}},
		{caseName("one tag", "two tags"), []string{"tag initial"}, []string{"tag1", "tag2"}},
	}

	r.writeLogs = false

	for _, test := range tests {
		r.CaseName = test.Case
		r.Project.Tags.Clear()
		r.Project.Tags.Add(test.InitialTags...)
		if err := r.Project.Save(); err != nil {
			t.Error(err)
			continue
		}

		r.Expected = append(test.InitialTags, test.NewTags...)

		r.Post(r.subPath("add"), map[string]interface{}{"data":test.NewTags})
		r.Check()
	}

}
func TestTagsRemove(t *testing.T) {
	r := NewTagsTest(t)
	defer r.Cleanup()

	caseName := func(initial, removed string) string {
		return fmt.Sprintf("remove %s from %s",removed, initial )
	}
	tests := []struct {
		Case    string
		Initial []string
		Remove  []string
		Expect	[]string
	}{
		{caseName("Empty", "Empty"), []string{}, []string{}, []string{}},
		{caseName("Empty", "one tag"), []string{}, []string{"tag"}, []string{}},
		{caseName("Empty", "two tags"), []string{}, []string{"tag1", "tag2"}, []string{}},
		{caseName("one tag", "two tags not in initial"), []string{"tag initial"}, []string{"tag1", "tag2"}, []string{"tag initial"}},
		{caseName("two tags", "two tags"), []string{"tag1", "tag2"}, []string{"tag1", "tag2"}, []string{}},
	}

	r.writeLogs = false

	for _, test := range tests {
		r.CaseName = test.Case
		r.Project.Tags.Clear()
		r.Project.Tags.Add(test.Initial...)
		if err := r.Project.Save(); err != nil {
			t.Error(err)
			continue
		}

		r.Expected = test.Expect

		r.Post(r.subPath("remove"), map[string]interface{}{"data":test.Remove})
		r.Check()
	}

}
