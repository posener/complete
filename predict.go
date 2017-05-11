package complete

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/posener/complete/match"
)

// Predictor implements a predict method, in which given
// command line arguments returns a list of options it predicts.
type Predictor interface {
	Predict(Args) []string
}

// PredictOr unions two predicate functions, so that the result predicate
// returns the union of their predication
func PredictOr(predictors ...Predictor) Predictor {
	return PredictFunc(func(a Args) (prediction []string) {
		for _, p := range predictors {
			if p == nil {
				continue
			}
			prediction = append(prediction, p.Predict(a)...)
		}
		return
	})
}

// PredictFunc determines what terms can follow a command or a flag
// It is used for auto completion, given last - the last word in the already
// in the command line, what words can complete it.
type PredictFunc func(Args) []string

// Predict invokes the predict function and implements the Predictor interface
func (p PredictFunc) Predict(a Args) []string {
	if p == nil {
		return nil
	}
	return p(a)
}

// PredictNothing does not expect anything after.
var PredictNothing Predictor

// PredictAnything expects something, but nothing particular, such as a number
// or arbitrary name.
var PredictAnything = PredictFunc(func(Args) []string { return nil })

// PredictSet expects specific set of terms, given in the options argument.
func PredictSet(options ...string) Predictor {
	return predictSet(options)
}

type predictSet []string

func (p predictSet) Predict(a Args) (prediction []string) {
	for _, m := range p {
		if match.Prefix(m, a.Last) {
			prediction = append(prediction, m)
		}
	}
	return
}

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
	return func(a Args) (prediction []string) {
		if strings.HasSuffix(a.Last, "/..") {
			return
		}
		dir := dirFromLast(a.Last)
		rel := !filepath.IsAbs(pattern)
		Log("looking for files in %s (last=%s)", dir, a.Last)
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
			if match.File(f, a.Last) {
				prediction = append(prediction, f)
			}
		}
		return
	}
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
