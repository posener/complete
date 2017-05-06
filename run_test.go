package complete

import (
	"os"
	"sort"
	"testing"
)

func TestCompleter_Complete(t *testing.T) {
	t.Parallel()

	// Set debug environment variable so logs will be printed
	if testing.Verbose() {
		os.Setenv(envDebug, "1")
	}

	// Change to tests directory for testing completion of files and directories
	err := os.Chdir("./tests")
	if err != nil {
		t.Fatal(err)
	}

	c := Command{
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
				Args: Predicate(PredictDirs).Or(PredictFiles("*.md")),
			},
		},
		Flags: map[string]Predicate{
			"-h":       PredictNothing,
			"-global1": PredictAnything,
			"-o":       PredictFiles("*.txt"),
		},
	}

	allGlobals := []string{}
	for sub := range c.Sub {
		allGlobals = append(allGlobals, sub)
	}
	for flag := range c.Flags {
		allGlobals = append(allGlobals, flag)
	}

	testTXTFiles := []string{"./a.txt", "./b.txt", "./c.txt"}

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
			want: []string{"./", "./dir", "./readme.md", "-flag2", "-flag3", "-h", "-global1", "-o"},
		},
		{
			args: "sub2 ./",
			want: []string{"./", "./readme.md", "./dir"},
		},
		{
			args: "sub2 re",
			want: []string{"./readme.md"},
		},
		{
			args: "sub2 -flag2 ",
			want: []string{"./", "./dir", "./readme.md", "-flag2", "-flag3", "-h", "-global1", "-o"},
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
			want: testTXTFiles,
		},
		{
			args: "-o ./no-su",
			want: []string{},
		},
		{
			args: "-o ./",
			want: testTXTFiles,
		},
		{
			args: "-o ",
			want: testTXTFiles,
		},
		{
			args: "-o ./read",
			want: []string{},
		},
		{
			args: "-o ./readme.md",
			want: []string{},
		},
		{
			args: "-o ./readme.md ",
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
			args, _ := getLine()

			got := complete(c, args)

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
