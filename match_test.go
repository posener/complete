package complete

import "testing"

func TestMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		m      Matcher
		prefix string
		want   bool
	}{
		{
			m:      MatchPrefix("abcd"),
			prefix: "",
			want:   true,
		},
		{
			m:      MatchPrefix("abcd"),
			prefix: "ab",
			want:   true,
		},
		{
			m:      MatchPrefix("abcd"),
			prefix: "ac",
			want:   false,
		},
		{
			m:      MatchPrefix(""),
			prefix: "ac",
			want:   false,
		},
		{
			m:      MatchPrefix(""),
			prefix: "",
			want:   true,
		},
		{
			m:      MatchFileName("file.txt"),
			prefix: "",
			want:   true,
		},
		{
			m:      MatchFileName("./file.txt"),
			prefix: "",
			want:   true,
		},
		{
			m:      MatchFileName("./file.txt"),
			prefix: "f",
			want:   true,
		},
		{
			m:      MatchFileName("./file.txt"),
			prefix: "file.",
			want:   true,
		},
		{
			m:      MatchFileName("./file.txt"),
			prefix: "./f",
			want:   true,
		},
		{
			m:      MatchFileName("./file.txt"),
			prefix: "other.txt",
			want:   false,
		},
		{
			m:      MatchFileName("./file.txt"),
			prefix: "/file.txt",
			want:   false,
		},
		{
			m:      MatchFileName("/file.txt"),
			prefix: "file.txt",
			want:   false,
		},
		{
			m:      MatchFileName("/file.txt"),
			prefix: "./file.txt",
			want:   false,
		},
		{
			m:      MatchFileName("/file.txt"),
			prefix: "/file.txt",
			want:   true,
		},
		{
			m:      MatchFileName("/file.txt"),
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
