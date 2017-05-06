package install

import "os"

type root struct{}

func (r root) Install(cmd string, bin string) error {
	completeLink := getBashCompletionDLink(cmd)
	err := r.Uninstall(cmd, bin)
	if err != nil {
		return err
	}
	return os.Symlink(bin, completeLink)
}

func (root) Uninstall(cmd string, bin string) error {
	completeLink := getBashCompletionDLink(cmd)
	if _, err := os.Stat(completeLink); err == nil {
		err := os.Remove(completeLink)
		if err != nil {
			return err
		}
	}
	return nil
}

func getBashCompletionDLink(cmd string) string {
	return "/etc/bash_completion.d/"+cmd
}

