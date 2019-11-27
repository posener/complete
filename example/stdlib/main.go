// stdlib shows how to have flags bash completion to an arbitrary Go program that uses the standard
// library flag package.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/posener/complete/v2"
)

var (
	// Add variables to the program.
	name      = flag.String("name", "", "Give your name")
	something = flag.String("something", "", "Expect somthing, but we don't know what, so no other completion options will be provided.")
	nothing   = flag.String("nothing", "", "Expect nothing after flag, so other completion can be provided.")
)

func main() {
	// Run the completion. Notice that since we are using standard library flags, only the flag
	// names will be completed and not their values.
	complete.CommandLine()

	// Parse the flags.
	flag.Parse()

	// Program logic.
	if *name == "" {
		fmt.Println("Your name is missing")
		os.Exit(1)
	}

	fmt.Println("Hi,", name)
}
