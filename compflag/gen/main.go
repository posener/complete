// Generates flags.go.
package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/posener/script"
)

const tmplGlob = "gen/*.go.gotmpl"

type flag struct {
	Name          string
	Type          string
	IsBool        bool
	CustomPredict bool
}

func (f flag) NewInternalTypeFuncName() string {
	return "new" + strings.Title(f.InternalTypeName())
}

func (f flag) InternalTypeName() string {
	return strings.ToLower(f.Name[:1]) + f.Name[1:] + "Value"
}

var flags = []flag{
	{Name: "String", Type: "string"},
	{Name: "Bool", Type: "bool", IsBool: true, CustomPredict: true},
	{Name: "Int", Type: "int"},
	{Name: "Duration", Type: "time.Duration"},
}

var tmpl = template.Must(template.ParseGlob(tmplGlob))

func main() {
	for _, t := range tmpl.Templates() {
		fileName := outFileName(t.Name())
		f, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		log.Printf("Writing %s", fileName)
		err = t.Execute(f, flags)
		if err != nil {
			panic(err)
		}

		// Format the file.
		err = script.ExecHandleStderr(os.Stderr, "goimports", "-w", fileName).ToStdout()
		if err != nil {
			panic(err)
		}
	}
}

func outFileName(templateName string) string {
	name := filepath.Base(templateName)
	// Remove .gotmpl suffix.
	return name[:strings.LastIndex(name, ".")]
}
