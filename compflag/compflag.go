// Package compflag provides a handful of standard library-compatible flags with bash complition capabilities.
//
// Usage
//
// 	import "github.com/posener/complete/v2/compflag"
//
// 	var (
// 		// Define flags...
// 		foo = compflag.String("foo", "", "")
// 	)
//
// 	func main() {
// 		compflag.Parse()
// 		// Main function.
// 	}
//
// Alternatively, the library can just be used with the standard library flag package:
//
// 	import (
// 		"flag"
// 		"github.com/posener/complete/v2/compflag"
// 	)
//
// 	var (
// 		// Define flags...
// 		foo = compflag.String("foo", "", "")
// 		bar = flag.String("bar", "", "")
// 	)
//
// 	func main() {
// 		complete.CommandLine()
// 		flag.Parse()
// 		// Main function.
// 	}
package compflag

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/posener/complete/v2"
)

// FlagSet is bash completion enabled flag.FlagSet.
type FlagSet flag.FlagSet

// Parse parses command line arguments.
func (fs *FlagSet) Parse(args []string) error {
	return (*flag.FlagSet)(fs).Parse(args)
}

func (fs *FlagSet) Visit(fn func(*flag.Flag))     { (*flag.FlagSet)(fs).Visit(fn) }
func (fs *FlagSet) VisitAll(fn func(*flag.Flag))  { (*flag.FlagSet)(fs).VisitAll(fn) }
func (fs *FlagSet) Arg(i int) string              { return (*flag.FlagSet)(fs).Arg(i) }
func (fs *FlagSet) Args() []string                { return (*flag.FlagSet)(fs).Args() }
func (fs *FlagSet) NArg() int                     { return (*flag.FlagSet)(fs).NArg() }
func (fs *FlagSet) NFlag() int                    { return (*flag.FlagSet)(fs).NFlag() }
func (fs *FlagSet) Name() string                  { return (*flag.FlagSet)(fs).Name() }
func (fs *FlagSet) PrintDefaults()                { (*flag.FlagSet)(fs).PrintDefaults() }
func (fs *FlagSet) Lookup(name string) *flag.Flag { return (*flag.FlagSet)(fs).Lookup(name) }
func (fs *FlagSet) Parsed() bool                  { return (*flag.FlagSet)(fs).Parsed() }

// Complete performs bash completion if needed.
func (fs *FlagSet) Complete() {
	complete.Complete(fs.Name(), complete.FlagSet((*flag.FlagSet)(CommandLine)))
}

var CommandLine = (*FlagSet)(flag.CommandLine)

// Parse parses command line arguments. It also performs bash completion when needed.
func Parse() {
	CommandLine.Complete()
	CommandLine.Parse(os.Args[1:])
}

func predictBool(value bool, prefix string) []string {
	// If false, typing the bool flag is expected to turn it on, so there is nothing to complete
	// after the flag.
	if !value {
		return nil
	}
	// Otherwise, suggest only to turn it off.
	return []string{"false"}
}

func parseString(s string) (string, error) { return s, nil }

func formatString(v string) string { return v }

func parseInt(s string) (int, error) { return strconv.Atoi(s) }

func formatInt(v int) string { return strconv.Itoa(v) }

func parseBool(s string) (bool, error) { return strconv.ParseBool(s) }

func formatBool(v bool) string { return strconv.FormatBool(v) }

func parseDuration(s string) (time.Duration, error) { return time.ParseDuration(s) }

func formatDuration(v time.Duration) string { return v.String() }
