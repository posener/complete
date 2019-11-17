// Package compflag provides a handful of standard library-compatible flags with bash complition capabilities.
//
// Usage
//
// 	import "github.com/posener/complete/compflag"
//
// 	var (
// 		// Define flags...
// 		foo = compflag.String("foo", "", "")
// 	)
//
// 	func main() {
// 		compflag.Parse("my-program")
// 		// Main function.
// 	}
//
// Alternatively, the library can just be used with the standard library flag package:
//
// 	import (
// 		"flag"
// 		"github.com/posener/complete/compflag"
// 	)
//
// 	var (
// 		// Define flags...
// 		foo = compflag.String("foo", "", "")
// 		bar = flag.String("bar", "", "")
// 	)
//
// 	func main() {
// 		complete.CommandLine("my-program")
// 		flag.ParseArgs()
// 		// Main function.
// 	}
package compflag

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/posener/complete"
)

// Flag options.
type Option func(*options)

// OptValues allows to set a desired set of valid values for the flag.
func OptValues(values ...string) Option {
	return func(o *options) { o.values = values }
}

// OptCheck enforces the valid values on the predicted flag.
func OptCheck() Option {
	return func(o *options) { o.check = true }
}

type options struct {
	values []string
	check  bool
}

func config(fs ...Option) options {
	var op options
	for _, f := range fs {
		f(&op)
	}
	return op
}

// FlagSet is bash completion enabled flag.FlagSet.
type FlagSet flag.FlagSet

// Parse parses command line arguments.
func (fs *FlagSet) Parse(args []string) error {
	return (*flag.FlagSet)(fs).Parse(args)
}

// Complete performs bash completion if needed.
func (fs *FlagSet) Complete(name string) {
	complete.Complete(name, complete.FlagSet((*flag.FlagSet)(CommandLine)))
}

func (fs *FlagSet) String(name string, value string, usage string, options ...Option) *string {
	p := new(string)
	(*flag.FlagSet)(fs).Var(newStringValue(value, p, config(options...)), name, usage)
	return p
}

func (fs *FlagSet) Bool(name string, value bool, usage string, options ...Option) *bool {
	p := new(bool)
	(*flag.FlagSet)(fs).Var(newBoolValue(value, p, config(options...)), name, usage)
	return p
}

func (fs *FlagSet) Int(name string, value int, usage string, options ...Option) *int {
	p := new(int)
	(*flag.FlagSet)(fs).Var(newIntValue(value, p, config(options...)), name, usage)
	return p
}

func (o options) checkValue(v string) error {
	if !o.check || len(o.values) == 0 {
		return nil
	}
	for _, vv := range o.values {
		if v == vv {
			return nil
		}
	}
	return fmt.Errorf("not in allowed values: %s", strings.Join(o.values, ","))
}

var CommandLine = (*FlagSet)(flag.CommandLine)

// Parse parses command line arguments. It also performs bash completion when needed.
func Parse(name string) {
	CommandLine.Complete(name)
	CommandLine.Parse(os.Args[1:])
}

func String(name string, value string, usage string, options ...Option) *string {
	return CommandLine.String(name, value, usage, options...)
}

func Bool(name string, value bool, usage string, options ...Option) *bool {
	return CommandLine.Bool(name, value, usage, options...)
}

func Int(name string, value int, usage string, options ...Option) *int {
	return CommandLine.Int(name, value, usage, options...)
}

type boolValue struct {
	v *bool
	options
}

func newBoolValue(val bool, p *bool, o options) *boolValue {
	*p = val
	return &boolValue{v: p, options: o}
}

func (b *boolValue) Set(val string) error {
	v, err := strconv.ParseBool(val)
	*b.v = v
	if err != nil {
		return fmt.Errorf("bad value for bool flag")
	}
	return b.checkValue(val)
}

func (b *boolValue) Get() interface{} { return bool(*b.v) }

func (b *boolValue) String() string {
	if b == nil || b.v == nil {
		return strconv.FormatBool(false)
	}
	return strconv.FormatBool(bool(*b.v))
}

func (b *boolValue) IsBoolFlag() bool { return true }

func (b *boolValue) Predict(_ string) []string {
	if b.values != nil {
		return b.values
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
	options
}

func newStringValue(val string, p *string, o options) *stringValue {
	*p = val
	return &stringValue{v: p, options: o}
}

func (s *stringValue) Set(val string) error {
	*s.v = val
	return s.options.checkValue(val)
}

func (s *stringValue) Get() interface{} {
	return string(*s.v)
}

func (s *stringValue) String() string {
	if s == nil || s.v == nil {
		return ""
	}
	return string(*s.v)
}

func (s *stringValue) Predict(_ string) []string {
	if s.values != nil {
		return s.values
	}
	return []string{""}
}

type intValue struct {
	v *int
	options
}

func newIntValue(val int, p *int, o options) *intValue {
	*p = val
	return &intValue{v: p, options: o}
}

func (i *intValue) Set(val string) error {
	v, err := strconv.ParseInt(val, 0, strconv.IntSize)
	*i.v = int(v)
	if err != nil {
		return fmt.Errorf("bad value for int flag")
	}
	return i.checkValue(val)
}

func (i *intValue) Get() interface{} { return int(*i.v) }

func (i *intValue) String() string {
	if i == nil || i.v == nil {
		return strconv.Itoa(0)
	}
	return strconv.Itoa(int(*i.v))
}

func (s *intValue) Predict(_ string) []string {
	if s.values != nil {
		return s.values
	}
	return []string{""}
}
