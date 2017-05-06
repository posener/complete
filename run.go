package complete

import (
	"fmt"
	"os"
	"strings"
)

const (
	envComplete = "COMP_LINE"
	envDebug    = "COMP_DEBUG"
)

// Run get a command, get the typed arguments from environment
// variable, and print out the complete options
func Run(c Command) {
	args := getLine()
	Log("Completing args: %s", args)

	options := complete(c, args)

	Log("Completion: %s", options)
	output(options)
}

// complete get a command an command line arguments and returns
// matching completion options
func complete(c Command, args []string) (matching []string) {
	options, _ := c.options(args[:len(args)-1])

	// choose only matching options
	l := last(args)
	for _, option := range options {
		if option.Matches(l) {
			matching = append(matching, option.String())
		}
	}
	return
}

func getLine() []string {
	line := os.Getenv(envComplete)
	if line == "" {
		panic("should be run as a complete script")
	}
	return strings.Split(line, " ")
}

func last(args []string) (last string) {
	if len(args) > 0 {
		last = args[len(args)-1]
	}
	return
}

func output(options []string) {
	// stdout of program defines the complete options
	for _, option := range options {
		fmt.Println(option)
	}
}
