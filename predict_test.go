package complete

import (
	"sort"
	"strings"
	"testing"
)

func TestPredicate(t *testing.T) {
	t.Parallel()
	initTests()

	tests := []struct {
		name string
		p    Predictor
		arg  string
		want []string
	}{
		{
			name: "set",
			p:    PredictSet("a", "b", "c"),
			want: []string{"a", "b", "c"},
		},
		{
			name: "set with does",
			p:    PredictSet("./..", "./x"),
			arg:  "./.",
			want: []string{"./.."},
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
			want: []string{"./", "./dir/", "./a.txt", "./b.txt", "./c.txt", "./.dot.txt"},
		},
		{
			name: "files/txt",
			p:    PredictFiles("*.txt"),
			arg:  "./dir/",
			want: []string{"./dir/"},
		},
		{
			name: "files/x",
			p:    PredictFiles("x"),
			arg:  "./dir/",
			want: []string{"./dir/", "./dir/x"},
		},
		{
			name: "files/*",
			p:    PredictFiles("x*"),
			arg:  "./dir/",
			want: []string{"./dir/", "./dir/x"},
		},
		{
			name: "files/md",
			p:    PredictFiles("*.md"),
			want: []string{"./", "./dir/", "./readme.md"},
		},
		{
			name: "dirs",
			p:    PredictDirs("*"),
			arg:  "./dir/",
			want: []string{"./dir/"},
		},
		{
			name: "dirs and files",
			p:    PredictFiles("*"),
			arg:  "./dir",
			want: []string{"./dir/", "./dir/x"},
		},
		{
			name: "dirs",
			p:    PredictDirs("*"),
			want: []string{"./", "./dir/"},
		},
		{
			name: "subdir",
			p:    PredictFiles("*"),
			arg:  "./dir/",
			want: []string{"./dir/", "./dir/x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"?arg='"+tt.arg+"'", func(t *testing.T) {

			matches := tt.p.Predict(newArgs(strings.Split(tt.arg, " ")))

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
