package complete

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/posener/complete/v2/internal/arg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCmd = &Command{
	Flags: map[string]Predictor{"cmd-flag": nil},
	Sub: map[string]*Command{
		"flags": &Command{
			Flags: map[string]Predictor{
				"values":    set{"a", "a a", "b"},
				"something": set{""},
				"nothing":   nil,
			},
		},
		"sub1": &Command{
			Flags: map[string]Predictor{"flag1": nil},
			Sub: map[string]*Command{
				"sub11": &Command{
					Flags: map[string]Predictor{"flag11": nil},
				},
				"sub12": &Command{},
			},
			Args: set{"arg1", "arg2"},
		},
		"sub2": &Command{},
		"args": &Command{
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
		out = os.Stdout
	}()

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
		{shouldPanic: true, line: "cmd", point: "4"},

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
			out = ioutil.Discard
			if tt.shouldPanic {
				assert.Panics(t, func() { testCmd.Complete("") })
			} else {
				testCmd.Complete("")
				assert.Equal(t, tt.shouldExit, isExit)
			}
		})
	}
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
