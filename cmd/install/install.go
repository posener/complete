package install

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type installer interface {
	Install(cmd, bin string) error
	Uninstall(cmd, bin string) error
}

// Install complete command given:
// cmd: is the command name
func Install(cmd string) error {
	shell := shellType()
	if shell == "" {
		return errors.New("must install through a terminatl")
	}
	i := getInstaller(shell)
	if i == nil {
		return fmt.Errorf("shell %s not supported", shell)
	}
	bin, err := getBinaryPath()
	if err != nil {
		return err
	}
	return i.Install(cmd, bin)
}

// Uninstall complete command given:
// cmd: is the command name
func Uninstall(cmd string) error {
	shell := shellType()
	if shell == "" {
		return errors.New("must uninstall through a terminatl")
	}
	i := getInstaller(shell)
	if i == nil {
		return fmt.Errorf("shell %s not supported", shell)
	}
	bin, err := getBinaryPath()
	if err != nil {
		return err
	}
	return i.Uninstall(cmd, bin)
}

func getInstaller(shell string) installer {
	switch shell {
	case "bash":
		return bash{}
	default:
		return nil
	}
}

func getBinaryPath() (string, error) {
	bin, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(bin)
}

func shellType() string {
	shell := os.Getenv("SHELL")
	return filepath.Base(shell)
}
