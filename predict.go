package complete

import (
	"os"
	"path/filepath"

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
	p := predictSet{}
	for _, o := range options {
		p = append(p, match.Prefix(o))
	}
	return p
}

type predictSet []match.Prefix

func (p predictSet) Predict(a Args) (prediction []string) {
	for _, m := range p {
		if m.Match(a.Last) {
			prediction = append(prediction, m.String())
		}
	}
	return
}

// PredictDirs will search for directories in the given started to be typed
// path, if no path was started to be typed, it will complete to directories
// in the current working directory.
func PredictDirs(pattern string) Predictor {
	return files(pattern, true, false)
}

// PredictFiles will search for files matching the given pattern in the started to
// be typed path, if no path was started to be typed, it will complete to files that
// match the pattern in the current working directory.
// To match any file, use "*" as pattern. To match go files use "*.go", and so on.
func PredictFiles(pattern string) Predictor {
	return files(pattern, false, true)
}

// PredictFilesOrDirs any file or directory that matches the pattern
func PredictFilesOrDirs(pattern string) Predictor {
	return files(pattern, true, true)
}

func files(pattern string, allowDirs, allowFiles bool) PredictFunc {
	return func(a Args) (prediction []string) {
		dir := dirFromLast(a.Last)
		Log("looking for files in %s (last=%s)", dir, a.Last)
		files, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			Log("failed glob operation with pattern '%s': %s", pattern, err)
		}
		if allowDirs {
			files = append(files, dir)
		}
		files = selectByType(files, allowDirs, allowFiles)
		if !filepath.IsAbs(pattern) {
			filesToRel(files)
		}
		// add all matching files to prediction
		for _, f := range files {
			if m := match.File(f); m.Match(a.Last) {
				prediction = append(prediction, m.String())
			}
		}
		return
	}
}

func selectByType(names []string, allowDirs bool, allowFiles bool) []string {
	filtered := make([]string, 0, len(names))
	for _, name := range names {
		stat, err := os.Stat(name)
		if err != nil {
			continue
		}
		if (stat.IsDir() && !allowDirs) || (!stat.IsDir() && !allowFiles) {
			continue
		}
		filtered = append(filtered, name)
	}
	return filtered
}

// filesToRel, change list of files to their names in the relative
// to current directory form.
func filesToRel(files []string) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	for i := range files {
		abs, err := filepath.Abs(files[i])
		if err != nil {
			continue
		}
		rel, err := filepath.Rel(wd, abs)
		if err != nil {
			continue
		}
		if rel != "." {
			rel = "./" + rel
		}
		if info, err := os.Stat(rel); err == nil && info.IsDir() {
			rel += "/"
		}
		files[i] = rel
	}
	return
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
