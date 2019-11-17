package main

import (
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/posener/complete/v2/predict"
)

// predictPackages completes packages in the directory pointed by a.Last
// and packages that are one level below that package.
func predictPackages(prefix string) (prediction []string) {
	prediction = []string{prefix}
	lastPrediction := ""
	for len(prediction) == 1 && (lastPrediction == "" || lastPrediction != prediction[0]) {
		// if only one prediction, predict files within this prediction,
		// for example, if the user entered 'pk' and we have a package named 'pkg',
		// which is the only package prefixed with 'pk', we will automatically go one
		// level deeper and give the user the 'pkg' and all the nested packages within
		// that package.
		lastPrediction = prediction[0]
		prefix = prediction[0]
		prediction = predictLocalAndSystem(prefix)
	}
	return
}

func predictLocalAndSystem(prefix string) []string {
	localDirs := predict.FilesSet(listPackages(directory(prefix))).Predict(prefix)
	// System directories are not actual file names, for example: 'github.com/posener/complete' could
	// be the argument, but the actual filename is in $GOPATH/src/github.com/posener/complete'. this
	// is the reason to use the PredictSet and not the PredictDirs in this case.
	s := systemDirs(prefix)
	sysDirs := predict.Set(s).Predict(prefix)
	return append(localDirs, sysDirs...)
}

// listPackages looks in current pointed dir and in all it's direct sub-packages
// and return a list of paths to go packages.
func listPackages(dir string) (directories []string) {
	// add subdirectories
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("failed reading directory %s: %s", dir, err)
		return
	}

	// build paths array
	paths := make([]string, 0, len(files)+1)
	for _, f := range files {
		if f.IsDir() {
			paths = append(paths, filepath.Join(dir, f.Name()))
		}
	}
	paths = append(paths, dir)

	// import packages according to given paths
	for _, p := range paths {
		pkg, err := build.ImportDir(p, 0)
		if err != nil {
			log.Printf("failed importing directory %s: %s", p, err)
			continue
		}
		directories = append(directories, pkg.Dir)
	}
	return
}

func systemDirs(dir string) (directories []string) {
	// get all paths from GOPATH environment variable and use their src directory
	paths := findGopath()
	for i := range paths {
		paths[i] = filepath.Join(paths[i], "src")
	}

	// normalize the directory to be an actual directory since it could be with an additional
	// characters after the last '/'.
	if !strings.HasSuffix(dir, "/") {
		dir = filepath.Dir(dir)
	}

	for _, basePath := range paths {
		path := filepath.Join(basePath, dir)
		files, err := ioutil.ReadDir(path)
		if err != nil {
			// path does not exists
			continue
		}
		// add the base path as one of the completion options
		switch dir {
		case "", ".", "/", "./":
		default:
			directories = append(directories, dir)
		}
		// add all nested directories of the base path
		// go supports only packages and not go files within the GOPATH
		for _, f := range files {
			if !f.IsDir() {
				continue
			}
			directories = append(directories, filepath.Join(dir, f.Name())+"/")
		}
	}
	return
}

func findGopath() []string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// By convention
		// See rationale at https://github.com/golang/go/issues/17262
		usr, err := user.Current()
		if err != nil {
			return nil
		}
		usrgo := filepath.Join(usr.HomeDir, "go")
		return []string{usrgo}
	}
	listsep := string([]byte{os.PathListSeparator})
	entries := strings.Split(gopath, listsep)
	return entries
}

func directory(prefix string) string {
	if info, err := os.Stat(prefix); err == nil && info.IsDir() {
		return fixPathForm(prefix, prefix)
	}
	dir := filepath.Dir(prefix)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return "./"
	}
	return fixPathForm(prefix, dir)
}

// fixPathForm changes a file name to a relative name
func fixPathForm(last string, file string) string {
	// get wording directory for relative name
	workDir, err := os.Getwd()
	if err != nil {
		return file
	}

	abs, err := filepath.Abs(file)
	if err != nil {
		return file
	}

	// if last is absolute, return path as absolute
	if filepath.IsAbs(last) {
		return fixDirPath(abs)
	}

	rel, err := filepath.Rel(workDir, abs)
	if err != nil {
		return file
	}

	// fix ./ prefix of path
	if rel != "." && strings.HasPrefix(last, ".") {
		rel = "./" + rel
	}

	return fixDirPath(rel)
}

func fixDirPath(path string) string {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() && !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}
