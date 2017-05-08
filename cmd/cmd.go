// Package cmd used for command line options for the complete tool
package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/posener/complete/cmd/install"
)

// Run is used when running complete in command line mode.
// this is used when the complete is not completing words, but to
// install it or uninstall it.
func Run(cmd string) {
	c := parseFlags(cmd)
	err := c.validate()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	if !c.yes && !prompt(c.action(), cmd) {
		fmt.Println("Cancelling...")
		os.Exit(2)
	}
	fmt.Println(c.action() + "ing...")
	if c.install {
		err = install.Install(cmd)
	} else {
		err = install.Uninstall(cmd)
	}
	if err != nil {
		fmt.Printf("%s failed! %s\n", c.action(), err)
		os.Exit(3)
	}
	fmt.Println("Done!")
}

// prompt use for approval
func prompt(action, cmd string) bool {
	fmt.Printf("%s completion for %s? ", action, cmd)
	var answer string
	fmt.Scanln(&answer)

	switch strings.ToLower(answer) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

// config for command line
type config struct {
	install   bool
	uninstall bool
	yes       bool
}

// create a config from command line arguments
func parseFlags(cmd string) config {
	var c config
	flag.BoolVar(&c.install, "install", false,
		fmt.Sprintf("Install completion for %s command", cmd))
	flag.BoolVar(&c.uninstall, "uninstall", false,
		fmt.Sprintf("Uninstall completion for %s command", cmd))
	flag.BoolVar(&c.yes, "y", false, "Don't prompt user for typing 'yes'")
	flag.Parse()
	return c
}

// validate the config
func (c config) validate() error {
	if c.install && c.uninstall {
		return errors.New("Install and uninstall are exclusive")
	}
	if !c.install && !c.uninstall {
		return errors.New("Must specify -install or -uninstall")
	}
	return nil
}

// action name according to the config values.
func (c config) action() string {
	if c.install {
		return "Install"
	}
	return "Uninstall"
}
