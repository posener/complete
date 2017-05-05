package complete

import (
	"os"
	"sort"
	"testing"
)

func TestCompleter_Complete(t *testing.T) {
	t.Parallel()

	if testing.Verbose() {
		os.Setenv(envDebug, "1")
	}

	c := Completer{
		Command: Command{
			Sub: map[string]Command{
				"sub1": {
					Flags: map[string]Predicate{
						"-flag1": PredictAnything,
						"-flag2": PredictNothing,
					},
				},
				"sub2": {
					Flags: map[string]Predicate{
						"-flag2": PredictNothing,
						"-flag3": PredictSet("opt1", "opt2", "opt12"),
					},
					Args: PredictDirs("./tests/").Or(PredictFiles("./tests/*.md")),
				},
			},
			Flags: map[string]Predicate{
				"-h":       PredictNothing,
				"-global1": PredictAnything,
				"-o":       PredictFiles("./tests/*.txt"),
			},
		},
	}

	allGlobals := []string{}
	for sub := range c.Sub {
		allGlobals = append(allGlobals, sub)
	}
	for flag := range c.Flags {
		allGlobals = append(allGlobals, flag)
	}

	testTXTFiles := []string{"./tests/a.txt", "./tests/b.txt", "./tests/c.txt"}

	tests := []struct {
		args string
		want []string
	}{
		{
			args: "",
			want: allGlobals,
		},
		{
			args: "-",
			want: []string{"-h", "-global1", "-o"},
		},
		{
			args: "-h ",
			want: allGlobals,
		},
		{
			args: "-global1 ", // global1 is known follow flag
			want: []string{},
		},
		{
			args: "sub",
			want: []string{"sub1", "sub2"},
		},
		{
			args: "sub1",
			want: []string{"sub1"},
		},
		{
			args: "sub2",
			want: []string{"sub2"},
		},
		{
			args: "sub1 ",
			want: []string{"-flag1", "-flag2", "-h", "-global1", "-o"},
		},
		{
			args: "sub2 ",
			want: []string{"./tests", "-flag2", "-flag3", "-h", "-global1", "-o"},
		},
		{
			args: "sub2 tests",
			want: []string{"./tests", "./tests/readme.md", "./tests/dir"},
		},
		{
			args: "sub2 tests/re",
			want: []string{"./tests/readme.md"},
		},
		{
			args: "sub2 -flag2 ",
			want: []string{"./tests", "-flag2", "-flag3", "-h", "-global1", "-o"},
		},
		{
			args: "sub1 -fl",
			want: []string{"-flag1", "-flag2"},
		},
		{
			args: "sub1 -flag1",
			want: []string{"-flag1"},
		},
		{
			args: "sub1 -flag1 ",
			want: []string{}, // flag1 is unknown follow flag
		},
		{
			args: "sub1 -flag2 ",
			want: []string{"-flag1", "-flag2", "-h", "-global1", "-o"},
		},
		{
			args: "-no-such-flag",
			want: []string{},
		},
		{
			args: "-no-such-flag ",
			want: allGlobals,
		},
		{
			args: "no-such-command",
			want: []string{},
		},
		{
			args: "no-such-command ",
			want: allGlobals,
		},
		{
			args: "-o ",
			want: []string{},
		},
		{
			args: "-o ./tes",
			want: []string{},
		},
		{
			args: "-o tests/",
			want: testTXTFiles,
		},
		{
			args: "-o tests",
			want: testTXTFiles,
		},
		{
			args: "-o ./compl",
			want: []string{},
		},
		{
			args: "-o ./complete.go",
			want: []string{},
		},
		{
			args: "-o ./complete.go ",
			want: allGlobals,
		},
		{
			args: "-o sub2 -flag3 ",
			want: []string{"opt1", "opt2", "opt12"},
		},
		{
			args: "-o sub2 -flag3 opt1",
			want: []string{"opt1", "opt12"},
		},
		{
			args: "-o sub2 -flag3 opt",
			want: []string{"opt1", "opt2", "opt12"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {

			tt.args = "cmd " + tt.args
			os.Setenv(envComplete, tt.args)
			args := getLine()

			got := c.complete(args)

			sort.Strings(tt.want)
			sort.Strings(got)

			if !equalSlices(got, tt.want) {
				t.Errorf("failed '%s'\ngot = %s\nwant: %s", t.Name(), got, tt.want)
			}
		})
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
