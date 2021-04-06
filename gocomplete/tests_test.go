package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"bou.ke/monkey"
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
			want: []string{
				"TestPredictions",
				"Example",
				"TestErrorSupression",
			},
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
	p := monkey.Patch(os.Exit, func(int) {})
	defer p.Unpatch()
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

func TestErrorSupression(t *testing.T) {
	defer monkey.Patch(os.Exit, func(int) {}).Unpatch()

	// Completion API environment variable names.
	const envLine, envPoint = "COMP_LINE", "COMP_POINT"

	// line should work out to
	//
	// * on most POSIX:
	//     go test /tmp/
	// * on MacOS X:
	//     go test /var/folders/<randomized_pathname>/T//
	// * on Windows:
	//     go test C:\Users\<username>\AppData\Local\Temp\
	//
	// which should trigger "failed importing directory: ... no
	// buildable Go source files..." error messages.
	var line = "go test " + os.TempDir() + string(os.PathSeparator)

	defer os.Unsetenv(envLine)
	defer os.Unsetenv(envPoint)
	os.Setenv(envLine, line)
	os.Setenv(envPoint, strconv.Itoa(len(line)))

	tests := []struct {
		verbose string
		wantErr bool
	}{{
		verbose: "",
		wantErr: false,
	}, {
		verbose: "1",
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(fmt.Sprintf(
			"%s=%q", envVerbose, tt.verbose,
		), func(t *testing.T) {
			// Discard completion (stdout).
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			defer w.Close()
			defer func(o *os.File) { os.Stdout = o }(os.Stdout)
			os.Stdout = w
			go io.Copy(ioutil.Discard, r)

			// "Redirect" stderr into a buffer.
			b := &strings.Builder{}
			log.SetOutput(b)

			defer os.Unsetenv(envVerbose)
			os.Setenv(envVerbose, tt.verbose)

			main()

			gotErr := b.Len() != 0
			if tt.wantErr && !gotErr {
				t.Fatal("want something in stderr, got nothing")
			} else if !tt.wantErr && gotErr {
				t.Fatalf("want nothing in stderr, got %d bytes",
					b.Len())
			}
		})
	}
}
