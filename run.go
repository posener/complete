// Package complete provides a tool for bash writing bash completion in go.
//
// Writing bash completion scripts is a hard work. This package provides an easy way
// to create bash completion scripts for any command, and also an easy way to install/uninstall
// the completion of the command.
package complete

import (
	"fmt"
	"os"
	"strings"

	"github.com/posener/complete/cmd"
)

const (
	envComplete = "COMP_LINE"
	envDebug    = "COMP_DEBUG"
)

// Run get a command, get the typed arguments from environment
// variable, and print out the complete options
func Run(c Command) {
	args, ok := getLine()
	if !ok {
		cmd.Run(c.Name)
		return
	}
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
		if option.Match(l) {
			matching = append(matching, option.String())
		}
	}
	return
}

func getLine() ([]string, bool) {
	line := os.Getenv(envComplete)
	if line == "" {
		return nil, false
	}
	return strings.Split(line, " "), true
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
