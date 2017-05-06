package complete

import (
	"path/filepath"
	"strings"
)

// Matcher matches itself to a string
// it is used for comparing a given argument to the last typed
// word, and see if it is a possible auto complete option.
type Matcher interface {
	String() string
	Match(prefix string) bool
}

// MatchPrefix is a simple Matcher, if the word is it's prefix, there is a match
type MatchPrefix string

func (a MatchPrefix) String() string {
	return string(a)
}

func (a MatchPrefix) Match(prefix string) bool {
	return strings.HasPrefix(string(a), prefix)
}

// MatchFileName is a file name Matcher, if the last word can prefix the
// MatchFileName path, there is a possible match
type MatchFileName string

func (a MatchFileName) String() string {
	return string(a)
}

func (a MatchFileName) Match(prefix string) bool {
	full, err := filepath.Abs(string(a))
	if err != nil {
		Log("failed getting abs path of %s: %s", a, err)
	}
	prefixFull, err := filepath.Abs(prefix)
	if err != nil {
		Log("failed getting abs path of %s: %s", prefix, err)
	}

	// if the file has the prefix as prefix,
	// but we don't want to show too many files, so, if it is in a deeper directory - omit it.
	return strings.HasPrefix(full, prefixFull) && (full == prefixFull || !strings.Contains(full[len(prefixFull)+1:], "/"))
}
