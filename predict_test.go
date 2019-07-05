package complete

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestPredicate(t *testing.T) {
	t.Parallel()
	initTests()

	tests := []struct {
		name    string
		p       Predictor
		argList []string
		want    []string
	}{
		{
			name: "set",
			p:    PredictSet("a", "b", "c"),
			want: []string{"a", "b", "c"},
		},
		{
			name: "set/empty",
			p:    PredictSet(),
			want: []string{},
		},
		{
			name: "anything",
			p:    PredictAnything,
			want: []string{},
		},
		{
			name: "or: word with nil",
			p:    PredictOr(PredictSet("a"), nil),
			want: []string{"a"},
		},
		{
			name: "or: nil with word",
			p:    PredictOr(nil, PredictSet("a")),
			want: []string{"a"},
		},
		{
			name: "or: nil with nil",
			p:    PredictOr(PredictNothing, PredictNothing),
			want: []string{},
		},
		{
			name: "or: word with word with word",
			p:    PredictOr(PredictSet("a"), PredictSet("b"), PredictSet("c")),
			want: []string{"a", "b", "c"},
		},
		{
			name: "files/txt",
			p:    PredictFiles("*.txt"),
			want: []string{"./", "dir/", "outer/", "a.txt", "b.txt", "c.txt", ".dot.txt"},
		},
		{
			name:    "files/txt",
			p:       PredictFiles("*.txt"),
			argList: []string{"./dir/"},
			want:    []string{"./dir/"},
		},
		{
			name:    "complete files inside dir if it is the only match",
			p:       PredictFiles("foo"),
			argList: []string{"./dir/", "./d"},
			want:    []string{"./dir/", "./dir/foo"},
		},
		{
			name:    "complete files inside dir when argList includes file name",
			p:       PredictFiles("*"),
			argList: []string{"./dir/f", "./dir/foo"},
			want:    []string{"./dir/foo"},
		},
		{
			name:    "files/md",
			p:       PredictFiles("*.md"),
			argList: []string{""},
			want:    []string{"./", "dir/", "outer/", "readme.md"},
		},
		{
			name:    "files/md with ./ prefix",
			p:       PredictFiles("*.md"),
			argList: []string{".", "./"},
			want:    []string{"./", "./dir/", "./outer/", "./readme.md"},
		},
		{
			name:    "dirs",
			p:       PredictDirs("*"),
			argList: []string{"di", "dir", "dir/"},
			want:    []string{"dir/"},
		},
		{
			name:    "dirs with ./ prefix",
			p:       PredictDirs("*"),
			argList: []string{"./di", "./dir", "./dir/"},
			want:    []string{"./dir/"},
		},
		{
			name:    "predict anything in dir",
			p:       PredictFiles("*"),
			argList: []string{"dir", "dir/", "di"},
			want:    []string{"dir/", "dir/foo", "dir/bar"},
		},
		{
			name:    "predict anything in dir with ./ prefix",
			p:       PredictFiles("*"),
			argList: []string{"./dir", "./dir/", "./di"},
			want:    []string{"./dir/", "./dir/foo", "./dir/bar"},
		},
		{
			name:    "root directories",
			p:       PredictDirs("*"),
			argList: []string{""},
			want:    []string{"./", "dir/", "outer/"},
		},
		{
			name:    "root directories with ./ prefix",
			p:       PredictDirs("*"),
			argList: []string{".", "./"},
			want:    []string{"./", "./dir/", "./outer/"},
		},
		{
			name:    "nested directories",
			p:       PredictDirs("*.md"),
			argList: []string{"ou", "outer", "outer/"},
			want:    []string{"outer/", "outer/inner/"},
		},
		{
			name:    "nested directories with ./ prefix",
			p:       PredictDirs("*.md"),
			argList: []string{"./ou", "./outer", "./outer/"},
			want:    []string{"./outer/", "./outer/inner/"},
		},
		{
			name:    "nested inner directory",
			p:       PredictFiles("*.md"),
			argList: []string{"outer/i"},
			want:    []string{"outer/inner/", "outer/inner/readme.md"},
		},
	}

	for _, tt := range tests {

		// no args in argList, means an empty argument
		if len(tt.argList) == 0 {
			tt.argList = append(tt.argList, "")
		}

		for _, arg := range tt.argList {
			t.Run(tt.name+"/arg="+arg, func(t *testing.T) {

				matches := tt.p.Predict(newArgs(arg))

				sort.Strings(matches)
				sort.Strings(tt.want)

				got := strings.Join(matches, ",")
				want := strings.Join(tt.want, ",")

				if got != want {
					t.Errorf("failed %s\ngot = %s\nwant: %s", t.Name(), got, want)
				}
			})
		}
	}
}

func TestMatchFile(t *testing.T) {
	t.Parallel()

	// Change to tests directory for testing completion of
	// files and directories
	err := os.Chdir("../tests")
	if err != nil {
		panic(err)
	}

	type matcherTest struct {
		prefix string
		want   bool
	}

	tests := []struct {
		long  string
		tests []matcherTest
	}{
		{
			long: "file.txt",
			tests: []matcherTest{
				{prefix: "", want: true},
				{prefix: "f", want: true},
				{prefix: "./f", want: true},
				{prefix: "./.", want: false},
				{prefix: "file.", want: true},
				{prefix: "./file.", want: true},
				{prefix: "file.txt", want: true},
				{prefix: "./file.txt", want: true},
				{prefix: "other.txt", want: false},
				{prefix: "/other.txt", want: false},
				{prefix: "/file.txt", want: false},
				{prefix: "/fil", want: false},
				{prefix: "/file.txt2", want: false},
				{prefix: "/.", want: false},
			},
		},
		{
			long: "./file.txt",
			tests: []matcherTest{
				{prefix: "", want: true},
				{prefix: "f", want: true},
				{prefix: "./f", want: true},
				{prefix: "./.", want: false},
				{prefix: "file.", want: true},
				{prefix: "./file.", want: true},
				{prefix: "file.txt", want: true},
				{prefix: "./file.txt", want: true},
				{prefix: "other.txt", want: false},
				{prefix: "/other.txt", want: false},
				{prefix: "/file.txt", want: false},
				{prefix: "/fil", want: false},
				{prefix: "/file.txt2", want: false},
				{prefix: "/.", want: false},
			},
		},
		{
			long: "/file.txt",
			tests: []matcherTest{
				{prefix: "", want: true},
				{prefix: "f", want: false},
				{prefix: "./f", want: false},
				{prefix: "./.", want: false},
				{prefix: "file.", want: false},
				{prefix: "./file.", want: false},
				{prefix: "file.txt", want: false},
				{prefix: "./file.txt", want: false},
				{prefix: "other.txt", want: false},
				{prefix: "/other.txt", want: false},
				{prefix: "/file.txt", want: true},
				{prefix: "/fil", want: true},
				{prefix: "/file.txt2", want: false},
				{prefix: "/.", want: false},
			},
		},
		{
			long: "./",
			tests: []matcherTest{
				{prefix: "", want: true},
				{prefix: ".", want: true},
				{prefix: "./", want: true},
				{prefix: "./.", want: false},
			},
		},
	}

	for _, tt := range tests {
		for _, ttt := range tt.tests {
			name := fmt.Sprintf("long=%q&prefix=%q", tt.long, ttt.prefix)
			t.Run(name, func(t *testing.T) {
				got := matchFile(tt.long, ttt.prefix)
				if got != ttt.want {
					t.Errorf("Failed %s: got = %t, want: %t", name, got, ttt.want)
				}
			})
		}
	}
}
