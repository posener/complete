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
		p    Predicate
		arg  string
		want []string
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
			name: "nothing",
			p:    PredictNothing,
			want: []string{},
		},
		{
			name: "or: word with nil",
			p:    PredictSet("a").Or(PredictNothing),
			want: []string{"a"},
		},
		{
			name: "or: nil with word",
			p:    PredictNothing.Or(PredictSet("a")),
			want: []string{"a"},
		},
		{
			name: "or: nil with nil",
			p:    PredictNothing.Or(PredictNothing),
			want: []string{},
		},
		{
			name: "or: word with word with word",
			p:    PredictSet("a").Or(PredictSet("b")).Or(PredictSet("c")),
			want: []string{"a", "b", "c"},
		},
		{
			name: "files/txt",
			p:    PredictFiles("*.txt"),
			want: []string{"./a.txt", "./b.txt", "./c.txt"},
		},
		{
			name: "files/txt",
			p:    PredictFiles("*.txt"),
			arg:  "./dir/",
			want: []string{},
		},
		{
			name: "files/x",
			p:    PredictFiles("x"),
			arg:  "./dir/",
			want: []string{"./dir/x"},
		},
		{
			name: "files/*",
			p:    PredictFiles("x*"),
			arg:  "./dir/",
			want: []string{"./dir/x"},
		},
		{
			name: "files/md",
			p:    PredictFiles("*.md"),
			want: []string{"./readme.md"},
		},
		{
			name: "dirs",
			p:    PredictDirs,
			arg:  "./dir/",
			want: []string{"./dir"},
		},
		{
			name: "dirs",
			p:    PredictDirs,
			want: []string{"./", "./dir"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/"+tt.arg, func(t *testing.T) {
			matchers := tt.p.predict(tt.arg)
			matchersString := []string{}
			for _, m := range matchers {
				matchersString = append(matchersString, m.String())
			}
			sort.Strings(matchersString)
			sort.Strings(tt.want)

			got := strings.Join(matchersString, ",")
			want := strings.Join(tt.want, ",")

			if got != want {
				t.Errorf("failed %s\ngot = %s\nwant: %s", tt.name, got, want)
			}
		})
	}
}
