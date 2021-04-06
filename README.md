# complete

[![Build Status](https://travis-ci.org/posener/complete.svg?branch=master)](https://travis-ci.org/posener/complete)
[![codecov](https://codecov.io/gh/posener/complete/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/complete)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/posener/complete)

Package complete is everything for bash completion and Go.

Writing bash completion scripts is a hard work, usually done in the bash scripting language.
This package provides:

* A library for bash completion for Go programs.

* A tool for writing bash completion script in the Go language. For any Go or non Go program.

* Bash completion for the `go` command line (See [./gocomplete](./gocomplete)).

* Library for bash-completion enabled flags (See [./compflag](./compflag)).

* Enables an easy way to install/uninstall the completion of the command.

The library and tools are extensible such that any program can add its one logic, completion types
or methologies.

## Go Command Bash Completion

[./gocomplete](./gocomplete) is the script for bash completion for the `go` command line. This is an example
that uses the `complete` package on the `go` command - the `complete` package can also be used to
implement any completions, see #usage.

Install:

1. Type in your shell:

```go
go get -u github.com/posener/complete/v2/gocomplete
COMP_INSTALL=1 gocomplete
```

2. Restart your shell

Uninstall by `COMP_UNINSTALL=1 gocomplete`

Features:

- Complete `go` command, including sub commands and flags.
- Complete packages names or `.go` files when necessary.
- Complete test names after `-run` flag.

## Complete Package

Supported shells:

- [x] bash
- [x] zsh
- [x] fish

The installation of completion for a command line tool is done automatically by this library by
running the command line tool with the `COMP_INSTALL` environment variable set. Uninstalling the
completion is similarly done by the `COMP_UNINSTALL` environment variable.
For example, if a tool called `my-cli` uses this library, the completion can install by running
`COMP_INSTALL=1 my-cli`.

## Usage

Add bash completion capabilities to any Go program. See [./example/command](./example/command).

```go
 import (
 	"flag"
 	"github.com/posener/complete/v2"
 	"github.com/posener/complete/v2/predict"
 )
 var (
 	// Add variables to the program.
 	name      = flag.String("name", "", "")
 	something = flag.String("something", "", "")
 	nothing   = flag.String("nothing", "", "")
 )
 func main() {
 	// Create the complete command.
 	// Here we define completion values for each flag.
 	cmd := &complete.Command{
	 	Flags: map[string]complete.Predictor{
 			"name":      predict.Set{"foo", "bar", "foo bar"},
 			"something": predict.Something,
 			"nothing":   predict.Nothing,
 		},
 	}
 	// Run the completion - provide it with the binary name.
 	cmd.Complete("my-program")
 	// Parse the flags.
 	flag.Parse()
 	// Program logic...
 }
```

This package also enables to complete flags defined by the standard library `flag` package.
To use this feature, simply call `complete.CommandLine` before `flag.Parse`. (See [./example/stdlib](./example/stdlib)).

```diff
 import (
 	"flag"
+	"github.com/posener/complete/v2"
 )
 var (
 	// Define flags here...
 	foo = flag.Bool("foo", false, "")
 )
 func main() {
 	// Call command line completion before parsing the flags - provide it with the binary name.
+	complete.CommandLine("my-program")
 	flag.Parse()
 }
```

If flag value completion is desired, it can be done by providing the standard library `flag.Var`
function a `flag.Value` that also implements the `complete.Predictor` interface. For standard
flag with values, it is possible to use the `github.com/posener/complete/v2/compflag` package.
(See [./example/compflag](./example/compflag)).

```diff
 import (
 	"flag"
+	"github.com/posener/complete/v2"
+	"github.com/posener/complete/v2/compflag"
 )
 var (
 	// Define flags here...
-	foo = flag.Bool("foo", false, "")
+	foo = compflag.Bool("foo", false, "")
 )
 func main() {
 	// Call command line completion before parsing the flags.
+	complete.CommandLine("my-program")
 	flag.Parse()
 }
```

Instead of calling both `complete.CommandLine` and `flag.Parse`, one can call just `compflag.Parse`
which does them both.

## Testing

For command line bash completion testing use the `complete.Test` function.

## Sub Packages

* [compflag](./compflag): Package compflag provides a handful of standard library-compatible flags with bash complition capabilities.

* [compflag/gen](./compflag/gen): Generates flags.go.

* [example/command](./example/command): command shows how to have bash completion to an arbitrary Go program using the `complete.Command` struct.

* [example/compflag](./example/compflag): compflag shows how to use the github.com/posener/complete/v2/compflag package to have auto bash completion for a defined set of flags.

* [example/stdlib](./example/stdlib): stdlib shows how to have flags bash completion to an arbitrary Go program that uses the standard library flag package.

* [gocomplete](./gocomplete): Package main is complete tool for the go command line

* [install](./install): Package install provide installation functions of command completion.

* [predict](./predict): Package predict provides helper functions for completion predictors.

## Examples

### OutputCapturing

ExampleComplete_outputCapturing demonstrates the ability to capture
the output of Complete() invocations, crucial for integration tests.

```golang
package main

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/posener/complete/v2/internal/arg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCmd = &Command{
	Flags: map[string]Predictor{"cmd-flag": nil},
	Sub: map[string]*Command{
		"flags": {
			Flags: map[string]Predictor{
				"values":    set{"a", "a a", "b"},
				"something": set{""},
				"nothing":   nil,
			},
		},
		"sub1": {
			Flags: map[string]Predictor{"flag1": nil},
			Sub: map[string]*Command{
				"sub11": {
					Flags: map[string]Predictor{"flag11": nil},
				},
				"sub12": {},
			},
			Args: set{"arg1", "arg2"},
		},
		"sub2": {},
		"args": {
			Args: set{"a", "a a", "b"},
		},
	},
}

func TestCompleter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args string
		want []string
	}{
		// Check empty flag name matching.

		{args: "flags ", want: []string{"-values", "-nothing", "-something", "-cmd-flag", "-h"}},
		{args: "flags -", want: []string{"-values", "-nothing", "-something", "-cmd-flag", "-h"}},
		{args: "flags --", want: []string{"--values", "--nothing", "--something", "--cmd-flag", "--help"}},
		// If started a flag with no matching prefix, expect to see all possible flags.
		{args: "flags -x", want: []string{"-values", "-nothing", "-something", "-cmd-flag", "-h"}},
		// Check prefix matching for chain of sub commands.
		{args: "sub1 sub11 -fl", want: []string{"-flag11", "-flag1"}},
		{args: "sub1 sub11 --fl", want: []string{"--flag11", "--flag1"}},

		// Test sub command completion.

		{args: "", want: []string{"flags", "sub1", "sub2", "args", "-h"}},
		{args: " ", want: []string{"flags", "sub1", "sub2", "args", "-h"}},
		{args: "f", want: []string{"flags"}},
		{args: "sub", want: []string{"sub1", "sub2"}},
		{args: "sub1", want: []string{"sub1"}},
		{args: "sub1 ", want: []string{"sub11", "sub12", "-h"}},
		// Suggest all sub commands if prefix is not known.
		{args: "x", want: []string{"flags", "sub1", "sub2", "args", "-h"}},

		// Suggest flag value.

		// A flag that has an empty completion should return empty completion. It "completes
		// something"... But it doesn't know what, so we should not complete anything else.
		{args: "flags -something ", want: []string{""}},
		{args: "flags -something foo", want: []string{""}},
		// A flag that have nil completion should complete all other options.
		{args: "flags -nothing ", want: []string{"-values", "-nothing", "-something", "-cmd-flag", "-h"}},
		// Trying to provide a value to the nothing flag should revert the phrase back to nothing.
		{args: "flags -nothing=", want: []string{}},
		// The flag value was not started, suggest all relevant values.
		{args: "flags -values ", want: []string{"a", "a\\ a", "b"}},
		{args: "flags -values a", want: []string{"a", "a\\ a"}},
		{args: "flags -values a\\", want: []string{"a\\ a"}},
		{args: "flags -values a\\ ", want: []string{"a\\ a"}},
		{args: "flags -values a\\ a", want: []string{"a\\ a"}},
		{args: "flags -values a\\ a ", want: []string{"-values", "-nothing", "-something", "-cmd-flag", "-h"}},
		{args: "flags -values \"a", want: []string{"\"a\"", "\"a a\""}},
		{args: "flags -values \"a ", want: []string{"\"a a\""}},
		{args: "flags -values \"a a", want: []string{"\"a a\""}},
		{args: "flags -values \"a a\"", want: []string{"\"a a\""}},
		{args: "flags -values \"a a\" ", want: []string{"-values", "-nothing", "-something", "-cmd-flag", "-h"}},

		{args: "flags -values=", want: []string{"a", "a\\ a", "b"}},
		{args: "flags -values=a", want: []string{"a", "a\\ a"}},
		{args: "flags -values=a\\", want: []string{"a\\ a"}},
		{args: "flags -values=a\\ ", want: []string{"a\\ a"}},
		{args: "flags -values=a\\ a", want: []string{"a\\ a"}},
		{args: "flags -values=a\\ a ", want: []string{"-values", "-nothing", "-something", "-cmd-flag", "-h"}},
		{args: "flags -values=\"a", want: []string{"\"a\"", "\"a a\""}},
		{args: "flags -values=\"a ", want: []string{"\"a a\""}},
		{args: "flags -values=\"a a", want: []string{"\"a a\""}},
		{args: "flags -values=\"a a\"", want: []string{"\"a a\""}},
		{args: "flags -values=\"a a\" ", want: []string{"-values", "-nothing", "-something", "-cmd-flag", "-h"}},

		// Complete positional arguments

		{args: "args ", want: []string{"-cmd-flag", "-h", "a", "a\\ a", "b"}},
		{args: "args a", want: []string{"a", "a\\ a"}},
		{args: "args a\\", want: []string{"a\\ a"}},
		{args: "args a\\ ", want: []string{"a\\ a"}},
		{args: "args a\\ a", want: []string{"a\\ a"}},
		{args: "args a\\ a ", want: []string{"-cmd-flag", "-h", "a", "a\\ a", "b"}},
		{args: "args \"a", want: []string{"\"a\"", "\"a a\""}},
		{args: "args \"a ", want: []string{"\"a a\""}},
		{args: "args \"a a", want: []string{"\"a a\""}},
		{args: "args \"a a\"", want: []string{"\"a a\""}},
		{args: "args \"a a\" ", want: []string{"-cmd-flag", "-h", "a", "a\\ a", "b"}},

		// Complete positional arguments from a parent command
		{args: "sub1 sub12 arg", want: []string{"arg1", "arg2"}},

		// Test help

		{args: "-", want: []string{"-h"}},
		{args: " -", want: []string{"-h"}},
		{args: "--", want: []string{"--help"}},
		{args: "-he", want: []string{"-help"}},
		{args: "-x", want: []string{"-help"}},
		{args: "flags -h", want: []string{"-h"}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			Test(t, testCmd, tt.args, tt.want)
		})
	}
}

func TestCompleter_error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args string
		err  string
	}{
		// Sub command already fully typed but unknown.
		{args: "x ", err: "unknown subcommand: x"},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			_, err := completer{Completer: testCmd, args: arg.Parse(tt.args)}.complete()
			require.Error(t, err)
			assert.Equal(t, tt.err, err.Error())
		})
	}
}

func TestComplete(t *testing.T) {
	defer func() {
		getEnv = os.Getenv
		exit = os.Exit
	}()

	in, out, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func(o *os.File) { os.Stdout = o }(os.Stdout)
	defer out.Close()
	os.Stdout = out
	go io.Copy(ioutil.Discard, in)

	tests := []struct {
		line, point string
		shouldExit  bool
		shouldPanic bool
		install     string
		uninstall   string
	}{
		{shouldExit: true, line: "cmd", point: "1"},
		{shouldExit: false, line: "", point: ""},
		{shouldPanic: true, line: "cmd", point: ""},
		{shouldPanic: true, line: "cmd", point: "a"},
		{shouldExit: true, line: "cmd", point: "4"},

		{shouldExit: true, install: "1"},
		{shouldExit: false, install: "a"},
		{shouldExit: true, uninstall: "1"},
		{shouldExit: false, uninstall: "a"},
	}

	for _, tt := range tests {
		t.Run(tt.line+"@"+tt.point, func(t *testing.T) {
			getEnv = func(env string) string {
				switch env {
				case "COMP_LINE":
					return tt.line
				case "COMP_POINT":
					return tt.point
				case "COMP_INSTALL":
					return tt.install
				case "COMP_UNINSTALL":
					return tt.uninstall
				case "COMP_YES":
					return "0"
				default:
					panic(env)
				}
			}
			isExit := false
			exit = func(int) {
				isExit = true
			}
			if tt.shouldPanic {
				assert.Panics(t, func() { testCmd.Complete("") })
			} else {
				testCmd.Complete("")
				assert.Equal(t, tt.shouldExit, isExit)
			}
		})
	}
}

// ExampleComplete_outputCapturing demonstrates the ability to capture
// the output of Complete() invocations, crucial for integration tests.
func main() {
	defer func(f func(int)) { exit = f }(exit)
	defer func(f getEnvFn) { getEnv = f }(getEnv)
	exit = func(int) {}

	// This is where the actual example starts:

	cmd := &Command{Sub: map[string]*Command{"bar": {}}}
	getEnv = promptEnv("foo b")

	Complete("foo", cmd)

}

type set []string

func (s set) Predict(_ string) []string {
	return s
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		s      string
		prefix string
		want   string
		wantOK bool
	}{
		{s: "ab", prefix: `b`, want: ``, wantOK: false},
		{s: "", prefix: `b`, want: ``, wantOK: false},
		{s: "ab", prefix: `a`, want: `ab`, wantOK: true},
		{s: "ab", prefix: `"'b`, want: ``, wantOK: false},
		{s: "ab", prefix: `"'a`, want: `"'ab'"`, wantOK: true},
		{s: "ab", prefix: `'"a`, want: `'"ab"'`, wantOK: true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"/"+tt.prefix, func(t *testing.T) {
			got, gotOK := hasPrefix(tt.s, tt.prefix)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOK, gotOK)
		})
	}
}

// getEnvFn emulates os.GetEnv by mapping one string to another.
type getEnvFn = func(string) string

// promptEnv returns getEnvFn that emulates the environment variables
// a shell would set when its prompt has the given contents.
var promptEnv = func(contents string) getEnvFn {
	return func(key string) string {
		switch key {
		case "COMP_LINE":
			return contents
		case "COMP_POINT":
			return strconv.Itoa(len(contents))
		}
		return ""
	}
}

```

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
