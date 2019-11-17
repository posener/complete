package main

import (
	"os"
	"sort"
	"testing"

	"github.com/posener/complete/v2"
)

func TestPredictions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		predictor complete.Predictor
		prefix    string
		want      []string
	}{
		{
			name:      "predict tests ok",
			predictor: predictTest,
			want:      []string{"TestPredictions", "Example"},
		},
		{
			name:      "predict benchmark ok",
			predictor: predictBenchmark,
			want:      []string{"BenchmarkFake"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.predictor.Predict(tt.prefix)
			if !equal(got, tt.want) {
				t.Errorf("Failed %s: got: %q, want: %q", t.Name(), got, tt.want)
			}
		})
	}
}

func BenchmarkFake(b *testing.B) {}

func Example() {
	os.Setenv("COMP_LINE", "go ru")
	os.Setenv("COMP_POINT", "5")
	main()
	// output: run
}

func equal(s1, s2 []string) bool {
	sort.Strings(s1)
	sort.Strings(s2)
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
