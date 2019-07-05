// Package match contains matchers that decide if to apply completion.
//
// This package is deprecated.
package match

import "strings"

// Match matches two strings
// it is used for comparing a term to the last typed
// word, the prefix, and see if it is a possible auto complete option.
//
// Deprecated.
type Match func(term, prefix string) bool

// Prefix is a simple Matcher, if the word is it's prefix, there is a match
// Match returns true if a has the prefix as prefix
//
// Deprecated.
func Prefix(long, prefix string) bool {
	return strings.HasPrefix(long, prefix)
}

// File returns true if prefix can match the file
//
// Deprecated.
func File(file, prefix string) bool {
	// special case for current directory completion
	if file == "./" && (prefix == "." || prefix == "") {
		return true
	}
	if prefix == "." && strings.HasPrefix(file, ".") {
		return true
	}

	file = strings.TrimPrefix(file, "./")
	prefix = strings.TrimPrefix(prefix, "./")

	return strings.HasPrefix(file, prefix)
}
