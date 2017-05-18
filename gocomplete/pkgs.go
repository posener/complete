package main

import (
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/posener/complete"
)

func predictPackages(a complete.Args) (prediction []string) {
	for {
		prediction = complete.PredictFilesSet(listPackages(a)).Predict(a)

		// if the number of prediction is not 1, we either have many results or
		// have no results, so we return it.
		if len(prediction) != 1 {
			return
		}

		// if the result is only one item, we might want to recursively check
		// for more accurate results.
		if prediction[0] == a.Last {
			return
		}

		// only try deeper, if the one item is a directory
		if stat, err := os.Stat(prediction[0]); err != nil || !stat.IsDir() {
			return
		}

		a.Last = prediction[0]
	}
}

func listPackages(a complete.Args) (dirctories []string) {
	dir := a.Directory()
	complete.Log("listing packages in %s", dir)
	// import current directory
	pkg, err := build.ImportDir(dir, 0)
	if err != nil {
		complete.Log("failed importing directory %s: %s", dir, err)
		return
	}
	dirctories = append(dirctories, pkg.Dir)

	// import subdirectories
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		complete.Log("failed reading directory %s: %s", dir, err)
		return
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		pkg, err := build.ImportDir(filepath.Join(dir, f.Name()), 0)
		if err != nil {
			complete.Log("failed importing subdirectory %s: %s", filepath.Join(dir, f.Name()), err)
			continue
		}
		dirctories = append(dirctories, pkg.Dir)
	}
	return
}
