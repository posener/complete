package match

import "fmt"

// Matcher matches itself to a string
// it is used for comparing a given argument to the last typed
// word, and see if it is a possible auto complete option.
type Matcher interface {
	fmt.Stringer
	Match(prefix string) bool
}
