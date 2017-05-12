package complete

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/posener/complete/match"
)

// PredictDirs will search for directories in the given started to be typed
// path, if no path was started to be typed, it will complete to directories
// in the current working directory.
func PredictDirs(pattern string) Predictor {
	return files(pattern, false)
}

// PredictFiles will search for files matching the given pattern in the started to
// be typed path, if no path was started to be typed, it will complete to files that
// match the pattern in the current working directory.
// To match any file, use "*" as pattern. To match go files use "*.go", and so on.
func PredictFiles(pattern string) Predictor {
	return files(pattern, true)
}

func files(pattern string, allowFiles bool) PredictFunc {

	// search for files according to arguments,
	// if only one directory has matched the result, search recursively into
	// this directory to give more results.
	return func(a Args) (prediction []string) {
		last := a.Last
		for {

			prediction = predictFiles(last, pattern, allowFiles)

			// if the number of prediction is not 1, we either have many results or
			// have no results, so we return it.
			if len(prediction) != 1 {
				return
			}

			// if the result is only one item, we might want to recursively check
			// for more accurate results.
			if prediction[0] == last { // avoid loop forever
				return
			}

			// only try deeper, if the one item is a directory
			if stat, err := os.Stat(prediction[0]); err != nil || !stat.IsDir() {
				return
			}

			last = prediction[0]
		}
	}
}

func predictFiles(last string, pattern string, allowFiles bool) (prediction []string) {
	if strings.HasSuffix(last, "/..") {
		return
	}

	dir := dirFromLast(last)
	rel := !filepath.IsAbs(pattern)
	files := listFiles(dir, pattern)

	// get wording directory for relative name
	workDir, err := os.Getwd()
	if err != nil {
		workDir = ""
	}

	// add dir if match
	files = append(files, dir)

	// add all matching files to prediction
	for _, f := range files {
		if stat, err := os.Stat(f); err != nil || (!stat.IsDir() && !allowFiles) {
			continue
		}

		// change file name to relative if necessary
		if rel && workDir != "" {
			f = toRel(workDir, f)
		}

		// test matching of file to the argument
		if match.File(f, last) {
			prediction = append(prediction, f)
		}
	}
	return

}

func listFiles(dir, pattern string) []string {
	m := map[string]bool{}
	if files, err := filepath.Glob(filepath.Join(dir, pattern)); err == nil {
		for _, f := range files {
			m[f] = true
		}
	}
	if dirs, err := ioutil.ReadDir(dir); err == nil {
		for _, d := range dirs {
			if d.IsDir() {
				m[d.Name()] = true
			}
		}
	}
	list := make([]string, 0, len(m))
	for k := range m {
		list = append(list, k)
	}
	return list
}

// toRel changes a file name to a relative name
func toRel(wd, file string) string {
	abs, err := filepath.Abs(file)
	if err != nil {
		return file
	}
	rel, err := filepath.Rel(wd, abs)
	if err != nil {
		return file
	}
	if rel != "." {
		rel = "./" + rel
	}
	if info, err := os.Stat(rel); err == nil && info.IsDir() {
		rel += "/"
	}
	return rel
}

// dirFromLast gives the directory of the current written
// last argument if it represents a file name being written.
// in case that it is not, we fall back to the current directory.
func dirFromLast(last string) string {
	if info, err := os.Stat(last); err == nil && info.IsDir() {
		return last
	}
	dir := filepath.Dir(last)
	_, err := os.Stat(dir)
	if err != nil {
		return "./"
	}
	return dir
}
