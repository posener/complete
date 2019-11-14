package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"regexp"
)

func functionsInFile(path string, regexp *regexp.Regexp) (tests []string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		log.Printf("Failed parsing %s: %s", path, err)
		return nil
	}
	for _, d := range f.Decls {
		if f, ok := d.(*ast.FuncDecl); ok {
			name := f.Name.String()
			if regexp == nil || regexp.MatchString(name) {
				tests = append(tests, name)
			}
		}
	}
	return
}
