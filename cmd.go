package complete

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/posener/complete/install"
)

func runCommandLine(cmd string) {
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
		err = install.Install(cmd, c.root)
	} else {
		err = install.Uninstall(cmd, c.root)
	}
	if err != nil {
		fmt.Printf("%s failed! %s\n", c.action(), err)
		os.Exit(3)
	}
	fmt.Println("Done!")
}

func prompt(action, cmd string) bool {
	fmt.Printf("%s bash completion for %s? ", action, cmd)
	var answer string
	fmt.Scanln(&answer)

	switch strings.ToLower(answer) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

type config struct {
	install   bool
	uninstall bool
	root      bool
	yes       bool
}

func parseFlags(cmd string) config {
	var c config
	flag.BoolVar(&c.install, "install", false,
		fmt.Sprintf("Install bash completion for %s command", cmd))
	flag.BoolVar(&c.uninstall, "uninstall", false,
		fmt.Sprintf("Uninstall bash completion for %s command", cmd))
	flag.BoolVar(&c.root, "root", false,
		"(Un)Install as root:\n"+
			"            (Un)Install at /etc/bash_completion.d/ (user should have write permissions to that directory).\n"+
			"            If not set, a complete command will be added(removed) to ~/.bashrc")
	flag.BoolVar(&c.yes, "y", false, "Don't prompt user for typing 'yes'")
	flag.Parse()
	return c
}

func (c config) validate() error {
	if c.install && c.uninstall {
		return errors.New("Install and uninstall are exclusive")
	}
	if !c.install && !c.uninstall {
		return errors.New("Must specify -install or -uninstall")
	}
	return nil
}

func (c config) action() string {
	if c.install {
		return "Install"
	}
	return "Uninstall"
}
