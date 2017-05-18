package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/posener/complete"
)

const goListFormat = `{"Name": "{{.Name}}", "Path": "{{.Dir}}", "FilesString": "{{.GoFiles}}"}`

func predictPackages(a complete.Args) (prediction []string) {
	dir := a.Directory()
	pkgs := listPackages(dir)

	files := make([]string, 0, len(pkgs))
	for _, p := range pkgs {
		files = append(files, p.Path)
	}
	return complete.PredictFilesSet(files).Predict(a)
}

type pack struct {
	Name        string
	Path        string
	FilesString string
	Files       []string
}

func listPackages(dir string) (pkgs []pack) {
	dir = strings.TrimRight(dir, "/") + "/..."
	out, err := exec.Command("go", "list", "-f", goListFormat, dir).Output()
	if err != nil {
		return
	}
	lines := bytes.Split(out, []byte("\n"))
	for _, line := range lines {
		var p pack
		err := json.Unmarshal(line, &p)
		if err != nil {
			continue
		}
		// parse the FileString from a string "[file1 file2 file3]" to a list of files
		p.Files = strings.Split(strings.Trim(p.FilesString, "[]"), " ")
		pkgs = append(pkgs, p)
	}
	return
}
