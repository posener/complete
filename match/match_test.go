package match

import (
	"os"
	"testing"
)

func TestMatch(t *testing.T) {
	t.Parallel()

	// Change to tests directory for testing completion of
	// files and directories
	err := os.Chdir("../tests")
	if err != nil {
		panic(err)
	}

	tests := []struct {
		m      Matcher
		prefix string
		want   bool
	}{
		{
			m:      Prefix("abcd"),
			prefix: "",
			want:   true,
		},
		{
			m:      Prefix("abcd"),
			prefix: "ab",
			want:   true,
		},
		{
			m:      Prefix("abcd"),
			prefix: "ac",
			want:   false,
		},
		{
			m:      Prefix(""),
			prefix: "ac",
			want:   false,
		},
		{
			m:      Prefix(""),
			prefix: "",
			want:   true,
		},
		{
			m:      File("file.txt"),
			prefix: "",
			want:   true,
		},
		{
			m:      File("./file.txt"),
			prefix: "",
			want:   true,
		},
		{
			m:      File("./file.txt"),
			prefix: "f",
			want:   true,
		},
		{
			m:      File("./file.txt"),
			prefix: "file.",
			want:   true,
		},
		{
			m:      File("./file.txt"),
			prefix: "./f",
			want:   true,
		},
		{
			m:      File("./file.txt"),
			prefix: "other.txt",
			want:   false,
		},
		{
			m:      File("./file.txt"),
			prefix: "/file.txt",
			want:   false,
		},
		{
			m:      File("/file.txt"),
			prefix: "file.txt",
			want:   false,
		},
		{
			m:      File("/file.txt"),
			prefix: "./file.txt",
			want:   false,
		},
		{
			m:      File("/file.txt"),
			prefix: "/file.txt",
			want:   true,
		},
		{
			m:      File("/file.txt"),
			prefix: "/fil",
			want:   true,
		},
	}

	for _, tt := range tests {
		name := tt.m.String() + "/" + tt.prefix
		t.Run(name, func(t *testing.T) {
			got := tt.m.Match(tt.prefix)
			if got != tt.want {
				t.Errorf("Failed %s: got = %t, want: %t", name, got, tt.want)
			}
		})
	}
}
