package shuffler

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/fs/matchers"
)

type matcherTest struct {
	Name      string
	M         matchers.Matcher
	positives map[string]string
	negatives map[string]string
}

func (mt *matcherTest)Run(t *testing.T) {
	matcherName := mt.Name
	for caseName, path := range mt.positives {
		if !mt.M.Match(path) {
			t.Errorf("%q MUST match the path %q. (%s)", matcherName, path, caseName)
		}
	}
	for caseName, path := range mt.negatives {
		if mt.M.Match(path) {
			t.Errorf("%s MUST NOT match the path %q (%s).", matcherName, path, caseName)
		}
	}
}

func TestShufflerDirMatcher(t *testing.T) {
	test := matcherTest{
		positives:map[string]string{
			"shuffler metadata dir name":".shuffler",
			"shuffler metadata dir name with leading ./":"./.shuffler",
			`shuffler metadata dir name with leading .\\`:".\\.shuffler",
		},

		negatives:map[string]string{
			"shuffler without leading dot":"shuffler",
			"shuffler with leading underscore":"_shuffler",
			"shuffler with leading letter":"ashuffler",
			"shuffler with other letters after":"shufflerr",
			"shuffler with other words":"shuffler dir",
		},
		M:ShufflerDirMatcher,
		Name:"ShufflerDirMatcher",

	}

	test.Run(t)
}

func TestMdDirContentMatcher(t *testing.T) {
	test := matcherTest{
		positives:map[string]string{
			"shuffler metadata dir name":".shuffler",
			"file under metadata dir":".shuffler/file",
			"file under metadata dir, windows separator":".shuffler\\file",
			"file deep under metadata dir":".shuffler/subdir/file",
			"file deep under metadata dir, mixed separators":".shuffler/subdir\\file",
			"file deep under metadata dir, windows sepaator":".shuffler\\subdir\\file",
		},

		negatives:map[string]string{
			"shuffler without leading dot":"shuffler/file",
			"shuffler with leading underscore":"_shuffler/file",
			"shuffler with leading letter":"ashuffler/file",
			"shuffler with other letters after":"shufflerr/file",
			"shuffler with other words":"shuffler dir/file",
		},
		M:MdDirContentMatcher,
		Name:"MdDirContentMatcher",

	}

	test.Run(t)

}

func TestMdFileMatcher(t *testing.T) {
	test := matcherTest{
		positives:map[string]string{
			"forward slash separator":".shuffler/metadata.json",
			"back slash separator":".shuffler\\metadata.json",
		},

		negatives:map[string]string{
			"shuffler without leading dot":"shuffler/file",
			"shuffler with leading underscore":"_shuffler/file",
			"shuffler with leading letter":"ashuffler/file",
			"shuffler with other letters after":"shufflerr/file",
			"shuffler with other words":"shuffler dir/file",
			"file under metadata dir":".shuffler/file",
			"file under metadata dir, windows separator":".shuffler\\file",
			"file deep under metadata dir":".shuffler/subdir/file",
			"file deep under metadata dir, mixed separators":".shuffler/subdir\\file",
			"file deep under metadata dir, windows sepaator":".shuffler\\subdir\\file",

		},
		M:MdFileMatcher,
		Name:"MdFileMatcher",

	}

	test.Run(t)

}
