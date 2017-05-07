package match

import "strings"

// Prefix is a simple Matcher, if the word is it's prefix, there is a match
type Prefix string

func (a Prefix) String() string {
	return string(a)
}

// Match returns true if a has the prefix as prefix
func (a Prefix) Match(prefix string) bool {
	return strings.HasPrefix(string(a), prefix)
}
