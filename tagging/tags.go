package tagging

import (
	"sort"
	"encoding/json"
	"strings"
	"math"
)

type Collection struct {
	tags []string
}

func (c *Collection)Add(names ...string) {
	tagsMap := make(map[string]struct{})
	for _, tag := range c.tags {
		tagsMap[tag] = struct{}{}
	}
	for _, tag := range names {
		tagsMap[tag] = struct{}{}
	}
	c.tags = make([]string, 0)
	for tag := range tagsMap {
		if tag == "" {
			continue
		}
		c.tags = append(c.tags, tag)
	}
	sort.Strings(c.tags)
}
func (c *Collection)Remove(names ...string) {
	result := make([]string, 0, int(math.Max(0, float64(len(c.tags) - len(names)))))
	nameLen := len(names)
	for _, tag := range c.tags {
		i := sort.SearchStrings(names, tag)
		if i >= nameLen || tag != names[i] {
			// tag not in names
			result = append(result, tag)
		}
	}
	c.tags = result
}

func (c *Collection)Set(names ...string) {
	c.tags = []string{}
	c.Add(names...)
}

func (c *Collection)Get() []string {
	return c.tags
}

func (c *Collection)Tags() []string {
	return c.Get()
}

func (c *Collection)String() string {
	return strings.Join(c.Get(), ";")
}

func (c *Collection)Clear() {
	c.tags = []string{}
}

func Parse(s string) ([]string, error) {
	return strings.Split(s, ";"), nil
}

func (c *Collection)Parse(s string) error {
	t, err := Parse(s)
	if err != nil {
		return err
	}
	c.Add(t...)
	return nil
}

func (c *Collection)MarshalJSON() ([]byte, error) {
	return json.Marshal(c.tags)
}

func (c *Collection)UnmarshalJSON(data []byte) error {
	var tags = []string{}
	if err := json.Unmarshal(data, &tags); err != nil {
		return err
	}
	if len(tags) > 0 {
		c.tags = make([]string, len(tags))
		c.Add(tags...)
	}
	return nil
}
