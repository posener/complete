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
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
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

func (fs *FlagSet) String(name string, value string, usage string, options ...predict.Option) *string {
	p := new(string)
	(*flag.FlagSet)(fs).Var(newStringValue(value, p, predict.Options(options...)), name, usage)
	return p
}

func (fs *FlagSet) Bool(name string, value bool, usage string, options ...predict.Option) *bool {
	p := new(bool)
	(*flag.FlagSet)(fs).Var(newBoolValue(value, p, predict.Options(options...)), name, usage)
	return p
}

func (fs *FlagSet) Int(name string, value int, usage string, options ...predict.Option) *int {
	p := new(int)
	(*flag.FlagSet)(fs).Var(newIntValue(value, p, predict.Options(options...)), name, usage)
	return p
}

func (fs *FlagSet) Duration(name string, value time.Duration, usage string, options ...predict.Option) *time.Duration {
	p := new(time.Duration)
	(*flag.FlagSet)(fs).Var(newDurationValue(value, p, predict.Options(options...)), name, usage)
	return p
}

var CommandLine = (*FlagSet)(flag.CommandLine)

// Parse parses command line arguments. It also performs bash completion when needed.
func Parse() {
	CommandLine.Complete()
	CommandLine.Parse(os.Args[1:])
}

func String(name string, value string, usage string, options ...predict.Option) *string {
	return CommandLine.String(name, value, usage, options...)
}

func Bool(name string, value bool, usage string, options ...predict.Option) *bool {
	return CommandLine.Bool(name, value, usage, options...)
}

func Int(name string, value int, usage string, options ...predict.Option) *int {
	return CommandLine.Int(name, value, usage, options...)
}

func Duration(name string, value time.Duration, usage string, options ...predict.Option) *time.Duration {
	return CommandLine.Duration(name, value, usage, options...)
}

type boolValue struct {
	v *bool
	predict.Config
}

func newBoolValue(val bool, p *bool, c predict.Config) *boolValue {
	*p = val
	return &boolValue{v: p, Config: c}
}

func (b *boolValue) Set(val string) error {
	v, err := strconv.ParseBool(val)
	*b.v = v
	if err != nil {
		return fmt.Errorf("bad value for bool flag")
	}
	return b.Check(val)
}

func (b *boolValue) Get() interface{} { return *b.v }

func (b *boolValue) String() string {
	if b == nil || b.v == nil {
		return strconv.FormatBool(false)
	}
	return strconv.FormatBool(*b.v)
}

func (b *boolValue) IsBoolFlag() bool { return true }

func (b *boolValue) Predict(prefix string) []string {
	if b.Predictor != nil {
		return b.Predictor.Predict(prefix)
	}
	// If false, typing the bool flag is expected to turn it on, so there is nothing to complete
	// after the flag.
	if !*b.v {
		return nil
	}
	// Otherwise, suggest only to turn it off.
	return []string{"false"}
}

type stringValue struct {
	v *string
	predict.Config
}

func newStringValue(val string, p *string, c predict.Config) *stringValue {
	*p = val
	return &stringValue{v: p, Config: c}
}

func (s *stringValue) Set(val string) error {
	*s.v = val
	return s.Check(val)
}

func (s *stringValue) Get() interface{} {
	return *s.v
}

func (s *stringValue) String() string {
	if s == nil || s.v == nil {
		return ""
	}
	return *s.v
}

func (s *stringValue) Predict(prefix string) []string {
	if s.Predictor != nil {
		return s.Predictor.Predict(prefix)
	}
	return []string{""}
}

type intValue struct {
	v *int
	predict.Config
}

func newIntValue(val int, p *int, c predict.Config) *intValue {
	*p = val
	return &intValue{v: p, Config: c}
}

func (i *intValue) Set(val string) error {
	v, err := strconv.ParseInt(val, 0, strconv.IntSize)
	*i.v = int(v)
	if err != nil {
		return fmt.Errorf("bad value for int flag")
	}
	return i.Check(val)
}

func (i *intValue) Get() interface{} { return *i.v }

func (i *intValue) String() string {
	if i == nil || i.v == nil {
		return strconv.Itoa(0)
	}
	return strconv.Itoa(*i.v)
}

func (s *intValue) Predict(prefix string) []string {
	if s.Predictor != nil {
		return s.Predictor.Predict(prefix)
	}
	return []string{""}
}

type durationValue struct {
	v *time.Duration
	predict.Config
}

func newDurationValue(val time.Duration, p *time.Duration, c predict.Config) *durationValue {
	*p = val
	return &durationValue{v: p, Config: c}
}

func (i *durationValue) Set(val string) error {
	v, err := time.ParseDuration(val)
	*i.v = v
	if err != nil {
		return fmt.Errorf("bad value for duration flag")
	}
	return i.Check(val)
}

func (i *durationValue) Get() interface{} { return *i.v }

func (i *durationValue) String() string {
	if i == nil || i.v == nil {
		return time.Duration(0).String()
	}
	return i.v.String()
}

func (s *durationValue) Predict(prefix string) []string {
	if s.Predictor != nil {
		return s.Predictor.Predict(prefix)
	}
	return []string{""}
}
