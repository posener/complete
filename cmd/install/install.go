package install

import (
	"os"
	"path/filepath"
)

type installer interface {
	Install(cmd, bin string) error
	Uninstall(cmd, bin string) error
}

func Install(cmd string, asRoot bool) error {
	bin, err := getBinaryPath()
	if err != nil {
		return err
	}
	return getInstaller(asRoot).Install(cmd, bin)
}

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
	} else {
		return home{}
	}
}

func getBinaryPath() (string, error) {
	bin, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(bin)
}
