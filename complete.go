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
	logger("Completing args: %s", args)

	options := c.complete(args)

	logger("Completion: %s", options)
	output(options)
}

func (c *Completer) complete(args []string) []string {
	all, _ := c.options(args[:len(args)-1])
	return c.chooseRelevant(last(args), all)
}

func (c *Completer) chooseRelevant(last string, list []Option) (options []string) {
	for _, sub := range list {
		if sub.Matches(last) {
			options = append(options, sub.String())
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
