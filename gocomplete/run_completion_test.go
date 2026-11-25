package main

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func completionOutput(t *testing.T, line string, point int) []string {
	t.Helper()

	oldLine := os.Getenv("COMP_LINE")
	oldPoint := os.Getenv("COMP_POINT")
	defer os.Setenv("COMP_LINE", oldLine)
	defer os.Setenv("COMP_POINT", oldPoint)
	os.Setenv("COMP_LINE", line)
	os.Setenv("COMP_POINT", strconv.Itoa(point))

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	main()
	_ = w.Close()
	os.Stdout = oldStdout
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	return strings.Fields(string(out))
}

func TestGoRunUsesFileCompletionAfterPackage(t *testing.T) {
	if err := os.WriteFile("generic-file.txt", nil, 0o600); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("generic-file.txt")

	got := completionOutput(t, "go run . ", len("go run . "))
	foundGenericFile := false
	for _, item := range got {
		if item == "generic-file.txt" {
			foundGenericFile = true
		}
	}
	if !foundGenericFile {
		t.Fatalf("expected generic file completion after go run package, got %q", got)
	}

	got = completionOutput(t, "go run . -", len("go run . -"))
	for _, item := range got {
		if item == "-race" || item == "-exec" {
			t.Fatalf("go run flags were offered after the package argument: %q", got)
		}
	}
}
