package match

import (
	"strings"
)

// File is a file name Matcher, if the last word can prefix the
// File path, there is a possible match
type File string

func (a File) String() string {
	return string(a)
}

// Match returns true if prefix's abs path prefixes a's abs path
func (a File) Match(prefix string) bool {

	// special case for current directory completion
	if a == "./" && (prefix == "." || prefix == "") {
		return true
	}

	cmp := strings.TrimPrefix(string(a), "./")
	prefix = strings.TrimPrefix(prefix, "./")
	return strings.HasPrefix(cmp, prefix)
}
