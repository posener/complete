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

type Completer struct {
	Command
}

func New(c Command) *Completer {
	return &Completer{Command: c}
}

func (c *Completer) Complete() {
	args := getLine()
	Log("Completing args: %s", args)

	options := c.complete(args)

	Log("Completion: %s", options)
	output(options)
}

func (c *Completer) complete(args []string) []string {
	all, _ := c.options(args[:len(args)-1])
	return c.chooseRelevant(last(args), all)
}

func (c *Completer) chooseRelevant(last string, options []Option) (relevant []string) {
	for _, option := range options {
		if option.Matches(last) {
			relevant = append(relevant, option.String())
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
