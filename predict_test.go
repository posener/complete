package complete

import (
	"os"
	"sort"
	"strings"
	"testing"
)

func TestPredicate(t *testing.T) {
	t.Parallel()
	initTests()

	tests := []struct {
		name       string
		p          Predictor
		argList    []string
		want       []string
		prepEnv    func() (string, map[string]string, error)
		cleanEnv   func(dirTreeBase string)
		checkEqual func(dirTreeMappings map[string]string, got []string) bool
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
			name:    "predict anything in home directory with `~` prefix",
			p:       PredictFiles("*"),
			argList: []string{"~/foo"},
			want:    []string{"~/foo", "~/foo/foo.md", "~/foo/foo-dir"},
			prepEnv: func() (string, map[string]string, error) {
				basePath, dirTreeMappings, err := CreateDirTree(
					`~`,
					"foo",
					[]FileProperties{
						FileProperties{
							FilePath:         "foo.md",
							FileParent:       "",
							FileType:         RegularFile,
							ModificationType: CREATE,
						},
						FileProperties{
							FilePath:         "foo-dir",
							FileParent:       "",
							FileType:         Directory,
							ModificationType: CREATE,
						},
					},
				)
				return basePath, dirTreeMappings, err
			},
			cleanEnv: func(dirTreeBase string) {
				os.RemoveAll(dirTreeBase)
			},
			checkEqual: func(dirTreeMappings map[string]string, got []string) bool {
				want := []string{dirTreeMappings["foo"], dirTreeMappings["foo/foo.md"], dirTreeMappings["foo/-dir"]}
				sort.Strings(got)
				sort.Strings(want)
				gotStr := strings.Join(got, ",")
				wantStr := strings.Join(want, ",")
				return gotStr == wantStr
			},
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

	var basePath string
	var err error
	var dirTreeMappings map[string]string

	for _, tt := range tests {

		if tt.prepEnv != nil {
			if basePath, dirTreeMappings, err = tt.prepEnv(); err != nil {
				t.Errorf("error setting up env. Error %v", err)
			}
		}
		if tt.cleanEnv != nil {
			defer tt.cleanEnv(basePath)
		}

		// no args in argList, means an empty argument
		if len(tt.argList) == 0 {
			tt.argList = append(tt.argList, "")
		}

		for _, arg := range tt.argList {
			t.Run(tt.name+"/arg="+arg, func(t *testing.T) {

				matches := tt.p.Predict(newArgs(arg))

				sort.Strings(matches)
				sort.Strings(tt.want)

				if tt.checkEqual != nil {
					tt.checkEqual(dirTreeMappings, matches)
					return
				}
				got := strings.Join(matches, ",")
				want := strings.Join(tt.want, ",")

				if got != want {
					t.Errorf("failed %s\ngot = %s\nwant: %s", t.Name(), got, want)
				}
			})
		}
	}
}
