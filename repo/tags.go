package repo

import (
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"sync"
	"strings"
	"encoding/json"
)

type tagsIndex struct {
	once sync.Once
	tags map[string][]id.ID
}

func (ti *tagsIndex)init() {
	if ti.tags == nil {
		ti.tags = make(map[string][]id.ID)
	}
}

// Get returns slice of ID that satisfy the query.
// For now, query is just a tag name, case insensitive.
func (ti *tagsIndex)Get(query string) []id.ID {
	query = strings.ToLower(query)
	return ti.tags[query]
}

func (ti *tagsIndex)Add(ID id.ID, tags ...string) {
	ti.once.Do(ti.init)
	TAGS:
	for _, tag := range tags {
		tagL := strings.ToLower(tag)
		for _, existingID := range ti.tags[tagL] {
			if ID == existingID {
				continue TAGS
			}
		}
		ti.tags[tagL] = append(ti.tags[tagL], ID)
	}
}

func (ti *tagsIndex)MarshalJSON() ([]byte, error) {
	return json.Marshal(ti.tags)
}

func (ti *tagsIndex)UnmarshalJSON(data []byte) error {

	return json.Unmarshal(data, &ti.tags)
}
