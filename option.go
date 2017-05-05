package complete

import (
	"path/filepath"
	"strings"
)

type Option interface {
	String() string
	Matches(prefix string) bool
}

type Arg string

func (a Arg) String() string {
	return string(a)
}

func (a Arg) Matches(prefix string) bool {
	return strings.HasPrefix(string(a), prefix)
}

type ArgFileName string

func (a ArgFileName) String() string {
	return string(a)
}

func (a ArgFileName) Matches(prefix string) bool {
	full, err := filepath.Abs(string(a))
	if err != nil {
		logger("failed getting abs path of %s: %s", a, err)
	}
	prefixFull, err := filepath.Abs(prefix)
	if err != nil {
		logger("failed getting abs path of %s: %s", prefix, err)
	}

	// if the file has the prefix as prefix,
	// but we don't want to show too many files, so, if it is in a deeper directory - omit it.
	return strings.HasPrefix(full, prefixFull) && (full == prefixFull || !strings.Contains(full[len(prefixFull)+1:], "/"))
}
