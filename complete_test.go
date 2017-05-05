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
						"-flag3": PredictNothing,
					},
				},
			},
			Flags: map[string]Predicate{
				"-h":       PredictNothing,
				"-global1": PredictAnything,
				"-o":       PredictFiles("./gocomplete/*.go"),
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
			want: []string{"-flag2", "-flag3", "-h", "-global1", "-o"},
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
			want: []string{"./gocomplete/complete.go"},
		},
		{
			args: "-o goco",
			want: []string{"./gocomplete/complete.go"},
		},
		{
			args: "-o ./goco",
			want: []string{"./gocomplete/complete.go"},
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
