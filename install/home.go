package install

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type home struct{}

func (home) Install(cmd, bin string) error {
	bashRCFileName, err := bashRCFileName()
	if err != nil {
		return err
	}
	completeCmd := completeCmd(cmd, bin)
	if isInFile(bashRCFileName, completeCmd) {
		return errors.New("Already installed in ~/.bashrc")
	}

	bashRC, err := os.OpenFile(bashRCFileName, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		return err
	}
	defer bashRC.Close()
	_, err = bashRC.WriteString(fmt.Sprintf("\n%s\n", completeCmd))
	return err
}

func (home) Uninstall(cmd, bin string) error {
	bashRC, err := bashRCFileName()
	if err != nil {
		return err
	}
	backup := bashRC + ".bck"
	err = copyFile(bashRC, backup)
	if err != nil {
		return err
	}
	completeCmd := completeCmd(cmd, bin)
	if !isInFile(bashRC, completeCmd) {
		return errors.New("Does not installed in ~/.bashrc")
	}
	temp, err := uninstallToTemp(bashRC, completeCmd)
	if err != nil {
		return err
	}

	err = copyFile(temp, bashRC)
	if err != nil {
		return err
	}

	return os.Remove(backup)
}

func completeCmd(cmd, bin string) string {
	return fmt.Sprintf("complete -C %s %s", bin, cmd)
}

func bashRCFileName() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(u.HomeDir, ".bashrc"), nil
}

func isInFile(name string, lookFor string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	r := bufio.NewReader(f)
	prefix := []byte{}
	for {
		line, isPrefix, err := r.ReadLine()
		if err == io.EOF {
			return false
		}
		if err != nil {
			return false
		}
		if isPrefix {
			prefix = append(prefix, line...)
			continue
		}
		line = append(prefix, line...)
		if string(line) == lookFor {
			return true
		}
		prefix = prefix[:0]
	}
	return false
}

func uninstallToTemp(bashRCFileName, completeCmd string) (string, error) {
	rf, err := os.Open(bashRCFileName)
	if err != nil {
		return "", err
	}
	defer rf.Close()
	wf, err := ioutil.TempFile("/tmp", "bashrc-")
	if err != nil {
		return "", err
	}
	defer wf.Close()

	r := bufio.NewReader(rf)
	prefix := []byte{}
	for {
		line, isPrefix, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if isPrefix {
			prefix = append(prefix, line...)
			continue
		}
		line = append(prefix, line...)
		str := string(line)
		if str == completeCmd {
			continue
		}
		wf.WriteString(str + "\n")
		prefix = prefix[:0]
	}
	return wf.Name(), nil
}

func copyFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
