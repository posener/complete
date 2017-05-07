package match

// Matcher matches itself to a string
// it is used for comparing a given argument to the last typed
// word, and see if it is a possible auto complete option.
type Matcher interface {
	String() string
	Match(prefix string) bool
}
