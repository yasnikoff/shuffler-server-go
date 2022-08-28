package matchers

import (
	"path/filepath"
	"regexp"
)

// All is matcher that matches all paths
var All all

type all struct{}

func (a all)Match(path string)bool{
	return true
}

type Matcher interface{
	Match(path string)bool
}

type Compound interface{
	Matcher
	Append(m ...Matcher)
}
type compound struct{
	matchers []Matcher
}
func NewCompound(m ...Matcher) Compound {
	c:=&compound{}
	c.Append(m...)
	return c
}

func (c *compound)Match(path string)bool{
	for _, m :=range c.matchers{
		if m.Match(path){
			return true
		}
	}
	return false
}

func (c *compound)Append(m ...Matcher){
	c.matchers=append(c.matchers, m...)
}

type Glob struct{
	Globs []string
}

func (m *Glob)Match(path string)bool{
	for _, pattern:=range m.Globs{

		if ok, err:=filepath.Match(pattern, path); ok && err==nil{
			return true
		}
	}
	return false
}

type regexpMatcher struct{
	regexes []*regexp.Regexp
}

func NewRegex(patterns ...string) Matcher {
	m:=&regexpMatcher{regexes:make([]*regexp.Regexp, len(patterns))}

	for i, pattern:=range patterns{
		m.regexes[i]=regexp.MustCompile(pattern)
	}

	return m
}

func (r *regexpMatcher)Match(path string)bool{
	for _, re:=range r.regexes{
		if re.Match([]byte(path)){
			return true
		}
	}
	return false
}

func Not(m Matcher)Matcher{
	return &not{Matcher:m}
}

type not struct{
	Matcher
}
func (n *not)Match(path string)bool{
	return !n.Matcher.Match(path)
}
