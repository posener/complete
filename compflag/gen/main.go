// Generates flags.go.
package main

import (
	"log"
	"strings"

	"github.com/posener/autogen"
)

//go:generate go run .

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

func main() {
	err := autogen.Execute(flags)
	if err != nil {
		log.Fatal(err)
	}
}
