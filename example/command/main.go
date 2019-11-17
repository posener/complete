// command shows how to have bash completion to an arbitrary Go program using the `complete.Command`
// struct.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
)

var (
	// Add variables to the program.
	name      = flag.String("name", "", "Give your name")
	something = flag.String("something", "", "Expect somthing, but we don't know what, so no other completion options will be provided.")
	nothing   = flag.String("nothing", "", "Expect nothing after flag, so other completion can be provided.")
)

func main() {
	// Create the complete command.
	// Here we define completion values for each flag.
	cmd := &complete.Command{
		Flags: map[string]complete.Predictor{
			"name":      predict.Set{"foo", "bar", "foo bar"},
			"something": predict.Something,
			"nothing":   predict.Nothing,
		},
	}

	// Run the completion.
	cmd.Complete("command")

	// Parse the flags.
	flag.Parse()

	// Program logic.
	if *name == "" {
		fmt.Println("Your name is missing")
		os.Exit(1)
	}

	fmt.Println("Hi,", name)
}
