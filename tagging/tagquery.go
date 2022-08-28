package tagging

import (
	"strings"
)

type Matcher struct {
	raw       string
	tagGroups []Collection // every group is an "AND" group of tags. groups are "OR" joined.
}

func matchGroup(tags []string, group []string) bool {

	GroupTags:
	for _, gt := range group {
		for _, t := range tags {
			if t == gt {
				continue GroupTags
			}
		}
		return false
	}
	return true
}

func (tm *Matcher)Match(tags *Collection) bool {

	for _, group := range tm.tagGroups {
		if matchGroup(tags.Get(), group.Get()) {
			return true
		}
	}
	return false
}

func (tm *Matcher)Parse(s string) error {
	result := []Collection{}
	for _, groupStr := range strings.Split(s, "||") {
		groupStr = strings.TrimSpace(groupStr)
		tagGroup := *new(Collection)
		if err := tagGroup.Parse(groupStr); err != nil {
			return err
		}
		result = append(result, tagGroup)
	}
	tm.tagGroups = result
	tm.raw = s
	return nil
}

func (tm *Matcher)String() string {
	return tm.raw
}
