package install

import (
	"os"
	"path/filepath"
)

type installer interface {
	Install(cmd, bin string) error
	Uninstall(cmd, bin string) error
}

// Install complete command given:
// cmd: is the command name
// asRoot: if true the completion will be installed in /etc/bash_complete.d
// otherwise the complete command will be added to the ~/.bashrc file.
func Install(cmd string, asRoot bool) error {
	bin, err := getBinaryPath()
	if err != nil {
		return err
	}
	return getInstaller(asRoot).Install(cmd, bin)
}

// Uninstall complete command given:
// cmd: is the command name
// asRoot: if true the completion will be removed from /etc/bash_complete.d
// otherwise the complete command will be removed from the ~/.bashrc file.
func Uninstall(cmd string, asRoot bool) error {
	bin, err := getBinaryPath()
	if err != nil {
		return err
	}
	return getInstaller(asRoot).Uninstall(cmd, bin)
}

func getInstaller(asRoot bool) installer {
	if asRoot {
		return root{}
	}
	return home{}
}

func getBinaryPath() (string, error) {
	bin, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(bin)
}
