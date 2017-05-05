package complete

import (
	"os"
	"sort"
	"testing"
)

func TestCompleter_Complete(t *testing.T) {
	t.Parallel()

	os.Setenv(envDebug, "1")

	c := Completer{
		Command: Command{
			Sub: map[string]Command{
				"sub1": {
					Flags: map[string]FlagOptions{
						"-flag1": FlagUnknownFollow,
						"-flag2": FlagNoFollow,
					},
				},
				"sub2": {
					Flags: map[string]FlagOptions{
						"-flag2": FlagNoFollow,
						"-flag3": FlagNoFollow,
					},
				},
			},
			Flags: map[string]FlagOptions{
				"-h":       FlagNoFollow,
				"-global1": FlagUnknownFollow,
			},
		},
		log: t.Logf,
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
			want: []string{"-h", "-global1"},
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
			want: []string{"-flag1", "-flag2", "-h", "-global1"},
		},
		{
			args: "sub2 ",
			want: []string{"-flag2", "-flag3", "-h", "-global1"},
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
			want: []string{"-flag1", "-flag2", "-h", "-global1"},
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
