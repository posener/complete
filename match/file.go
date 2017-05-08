package match

import (
	"path/filepath"
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
	full, err := filepath.Abs(string(a))
	if err != nil {
		return false
	}
	prefixFull, err := filepath.Abs(prefix)
	if err != nil {
		return false
	}

	// if the file has the prefix as prefix,
	// but we don't want to show too many files, so, if it is in a deeper directory - omit it.
	return strings.HasPrefix(full, prefixFull) && (full == prefixFull || !strings.Contains(full[len(prefixFull)+1:], "/"))
}
