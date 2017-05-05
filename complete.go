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
	log func(format string, args ...interface{})
}

func New(c Command) *Completer {
	return &Completer{
		Command: c,
		log:     logger(),
	}
}

func (c *Completer) Complete() {
	args := getLine()
	c.log("Completing args: %s", args)

	options := c.complete(args)

	c.log("Completion: %s", options)
	output(options)
}

func (c *Completer) complete(args []string) []string {
	all, _ := c.options(args[:len(args)-1])
	return c.chooseRelevant(last(args), all)
}

func (c *Completer) chooseRelevant(last string, list []string) (opts []string) {
	if last == "" {
		return list
	}
	for _, sub := range list {
		if strings.HasPrefix(sub, last) {
			opts = append(opts, sub)
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
