// compflag shows how to use the github.com/posener/complete/v2/compflag package to have auto bash
// completion for a defined set of flags.
package main

import (
	"fmt"
	"os"

	"github.com/posener/complete/v2/compflag"
	"github.com/posener/complete/v2/predict"
)

var (
	// Add variables to the program. Since we are using the compflag library, we can pass options to
	// enable bash completion to the flag values.
	name      = compflag.String("name", "", "Give your name", predict.OptValues("foo", "bar", "foo bar"))
	something = compflag.String("something", "", "Expect somthing, but we don't know what, so no other completion options will be provided.", predict.OptValues(""))
	nothing   = compflag.String("nothing", "", "Expect nothing after flag, so other completion can be provided.")
)

func main() {
	// Parse flags and perform bash completion if needed.
	compflag.Parse()

	// Program logic.
	if *name == "" {
		fmt.Println("Your name is missing")
		os.Exit(1)
	}

	fmt.Println("Hi,", name)
}
