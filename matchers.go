package shuffler

import "bitbucket.org/yasnikoff/shuffler/fs/matchers"

// ShufflerDirMatcher matches the name of  shuffler metadata directory
var ShufflerDirMatcher = matchers.NewRegex(`^((./)|(.\\))?\.shuffler$`)

// MdFileMatcher matches shuffler's metadata.json file
var MdFileMatcher = matchers.NewRegex(`^\.shuffler[/,\\]metadata.json`)

// MdDirContentMatcher matches files under  shuffler metadata directory
var MdDirContentMatcher = matchers.NewRegex(`^\.shuffler[/,\\]*`)

// ProjectDuplicateFilesMatcher matches files that must be copied when duplicating project
var ProjectDuplicateFilesMatcher = matchers.Not(MdFileMatcher)

// ProjectMoveFileMatcher matches files that must be moved when project is moved
var ProjectMoveFileMatcher = matchers.All
