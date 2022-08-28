package repo_test

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/repo/repo_fixtures"
	"path/filepath"
	"bitbucket.org/yasnikoff/shuffler/tagging"
	"encoding/json"
	"reflect"
	"bitbucket.org/yasnikoff/shuffler"
	"bitbucket.org/yasnikoff/shuffler/repo"
	"bitbucket.org/yasnikoff/shuffler/types/id"
)

func TestTags_JSON(t *testing.T) {
	t.Parallel()
	tc := &tagging.Collection{}
	Tags := []string{"tag1", "tag2"}
	tc.Add(Tags...)
	data, err := json.Marshal(tc)
	if err != nil {
		t.Fatal(err)
	}
	tc2 := &tagging.Collection{}
	if err := json.Unmarshal(data, tc2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(tc.Tags(), tc2.Tags()) {
		t.Errorf("unmarshaled tags collection has Tags:\n%s\noriginal collection has Tags:\n%s", tc2.Tags(), tc.Tags())
	}

}

type TagFixture struct {
	repo_fixtures.RepoFixture
	TagsByPath map[string][]string
	t          *testing.T
}

func NewTagFixture(t *testing.T, TagsByPath map[string][]string) *TagFixture {
	if TagsByPath==nil{
		TagsByPath=make(map[string][]string, 0)
	}
	f := &TagFixture{
		RepoFixture:*repo_fixtures.NewDefault(t),
		TagsByPath:TagsByPath,
		t:t,
	}
	for path, tags := range f.TagsByPath {
		f.Add(path, tags...)
	}

	// update index after adding tags
	if err := f.Repo.UpdateIndex(); err != nil {
		t.Fatal(err)
	}
	return f
}

func (f *TagFixture)Add(path string, tags ...string) {
	pr := f.NewProjectWithPath(path)
	pr.Tags.Add(tags...)
	if err:=pr.Save();err!=nil{
		f.t.Fatal(err)
	}
}

type tagTest struct {
	query         string
	expectedPaths []string
}

func runTagTest(test tagTest, f *TagFixture, t *testing.T) {
	returnedIDs := f.Repo.TagsIndex.Get(test.query)
	if len(returnedIDs) != len(test.expectedPaths) {
		t.Errorf("Tags index Get() returned wrong number of IDs. Got %d, Want %d",
			len(returnedIDs), len(test.expectedPaths))
	}
	IDS:
	for _, ID := range returnedIDs {

		pr, err := f.Repo.GetByID(ID)
		if err != nil {
			t.Error(err)
			continue IDS
		}
		path, err := filepath.Rel(f.Path(), pr.Path())
		if err != nil {
			t.Error(err)
		}

		for _, expPath := range test.expectedPaths {
			if path == expPath {
				continue IDS
			}
		}
		t.Errorf("Tags index Get() returned path %q that is not in expected paths for query %q", path, test.query)
	}
}

func TestTags_Intersection(t *testing.T) {
	t.Parallel()
	tagsByPath := map[string][]string{
		"Project1": []string{"tag1", "tag2"},
		"Project2": []string{"tag2", "tag3"},
	}

	f := NewTagFixture(t, tagsByPath)

	defer f.Cleanup()

	tests := []tagTest{
		{"tag1", []string{"Project1"}},
		{"tag2", []string{"Project1", "Project2"}},
		{"tag3", []string{"Project2"}},
	}

	for _, test := range tests {
		runTagTest(test, f, t)
	}
}

func TestTags_UpdateIndexTwice(t *testing.T) {
	t.Parallel()
	tagsByPath := map[string][]string{
		"Project1": []string{"tag1", "tag2"},
		"Project2": []string{"tag2", "tag3"},
	}

	f := NewTagFixture(t, tagsByPath)
	defer f.Cleanup()

	// updating index the second time should not change it (for example, insert ID twice)
	if err := f.Repo.UpdateIndex(); err != nil {
		t.Fatal(err)
	}

	tests := []tagTest{
		{"tag1", []string{"Project1"}},
		{"tag2", []string{"Project1", "Project2"}},
		{"tag3", []string{"Project2"}},
	}

	for _, test := range tests {
		runTagTest(test, f, t)
	}
}

func TestTags_UpdateIndex(t *testing.T) {
	t.Parallel()

	f := NewTagFixture(t, map[string][]string{
		"Project1": []string{"tag"},
	})

	defer f.Cleanup()

	//f.Repo.Reload()
	f.Add("Project2", "tag")

	f.Repo.UpdateIndexForPath(f.Abs("Project2"))
	tests := []tagTest{
		{"tag", []string{"Project1", "Project2"}},
	}

	for _, test := range tests {
		runTagTest(test, f, t)
	}
}

func TestTags_UpdateIndexReload(t *testing.T) {
	t.Parallel()

	f := NewTagFixture(t, map[string][]string{
		"Project1": []string{"tag"},
	})

	defer f.Cleanup()

	f.Repo.Reload()
	f.Add("Project2", "tag")
	f.Repo.UpdateIndexForPath(f.Abs("Project2"))
	tests := []tagTest{
		{"tag", []string{"Project1", "Project2"}},
	}

	for _, test := range tests {
		runTagTest(test, f, t)
	}
}

func TestTagQuery(t *testing.T) {
	t.Parallel()

	var q shuffler.ProjectQuery
	q = &repo.TagQuery{}

	tag := "tag"
	query := "tag"
	q.Parse(query)

	f := repo_fixtures.NewDefault(t)
	defer f.Cleanup()

	pr := f.NewProjectWithPath("project")

	var result []id.ID
	var expectedResult []id.ID
	var err error

	result, err = q.Execute(f.Repo)
	expectedResult = []id.ID{}
	if err != nil {
		t.Fatal(err)
	}

	if len(result) > 0 {
		t.Errorf("Tag query on repo with no tags returned \n%s\nwant empty result", result, expectedResult)
	}

	pr.Tags.Add(tag)

	if err := pr.Save(); err != nil {
		t.Fatal(err)
	}

	if err := f.Repo.UpdateIndex(); err != nil {
		t.Fatal(err)
	}

	result, err = q.Execute(f.Repo)
	expectedResult = []id.ID{pr.ID()}
	if err != nil {
		t.Fatal(err)
	}

	if !IDSlicesEqual(result, expectedResult) {
		t.Errorf("Tag query returned \n%s\nwant\n%s", result, expectedResult)
	}

}

func IDSlicesEqual(s1, s2 []id.ID) bool {

	return reflect.DeepEqual(s1, s2)
	/*if len(s1) != len(s2) {
		return false
	}
	PR1:
	for _, pr1 := range s1 {
		for _, pr2 := range s2 {
			if pr1 == pr2 {
				continue PR1
			}
		}
		return false
	}
	return true*/
}
