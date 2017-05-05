package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/posener/complete"
)

func predictTest(testType string) complete.Predicate {
	return complete.Predicate{
		Predictor: func() []complete.Option {
			tests := testNames(testType)
			options := make([]complete.Option, len(tests))
			for i := range tests {
				options[i] = complete.Arg(tests[i])
			}
			return options
		},
	}
}

// get all test names in current directory
func testNames(testType string) (tests []string) {
	filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		// if not a test file, skip
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		// inspect test file and append all the test names
		tests = append(tests, testsInFile(testType, path)...)
		return nil
	})
	return
}

func testsInFile(testType, path string) (tests []string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		complete.Log("Failed parsing %s: %s", path, err)
		return nil
	}
	for _, d := range f.Decls {
		if f, ok := d.(*ast.FuncDecl); ok {
			name := f.Name.String()
			if strings.HasPrefix(name, testType) {
				tests = append(tests, name)
			}
		}
	}
	return
}
