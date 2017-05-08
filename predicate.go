package complete

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/posener/complete/match"
)

// Predicate determines what terms can follow a command or a flag
// It is used for auto completion, given last - the last word in the already
// in the command line, what words can complete it.
type Predicate func(last string) []match.Matcher

// Or unions two predicate functions, so that the result predicate
// returns the union of their predication
func (p Predicate) Or(other Predicate) Predicate {
	if p == nil {
		return other
	}
	if other == nil {
		return p
	}
	return func(last string) []match.Matcher { return append(p.predict(last), other.predict(last)...) }
}

func (p Predicate) predict(last string) []match.Matcher {
	if p == nil {
		return nil
	}
	return p(last)
}

// PredictNothing does not expect anything after.
var PredictNothing Predicate

// PredictAnything expects something, but nothing particular, such as a number
// or arbitrary name.
func PredictAnything(last string) []match.Matcher { return nil }

// PredictSet expects specific set of terms, given in the options argument.
func PredictSet(options ...string) Predicate {
	return func(last string) []match.Matcher {
		ret := make([]match.Matcher, len(options))
		for i := range options {
			ret[i] = match.Prefix(options[i])
		}
		return ret
	}
}

// PredictDirs will search for directories in the given started to be typed
// path, if no path was started to be typed, it will complete to directories
// in the current working directory.
func PredictDirs(last string) (options []match.Matcher) {
	path := dirFromLast(last)
	dirs := []string{}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	// if given path is not absolute, return relative paths
	if !filepath.IsAbs(path) {
		filesToRel(dirs)
	}
	return filesToMatchers(dirs)
}

// PredictFiles will search for files matching the given pattern in the started to
// be typed path, if no path was started to be typed, it will complete to files that
// match the pattern in the current working directory.
// To match any file, use "*" as pattern. To match go files use "*.go", and so on.
func PredictFiles(pattern string) Predicate {
	return func(last string) []match.Matcher {
		dir := dirFromLast(last)
		files, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			Log("failed glob operation with pattern '%s': %s", pattern, err)
		}
		if !filepath.IsAbs(pattern) {
			filesToRel(files)
		}
		return filesToMatchers(files)
	}
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
		if rel == "." {
			rel = ""
		}
		if !strings.HasPrefix(rel, ".") {
			rel = "./" + rel
		}
		files[i] = rel
	}
	return
}

func filesToMatchers(files []string) []match.Matcher {
	options := make([]match.Matcher, len(files))
	for i, f := range files {
		options[i] = match.File(f)
	}
	return options
}

// dirFromLast gives the directory of the current written
// last argument if it represents a file name being written.
// in case that it is not, we fall back to the current directory.
func dirFromLast(last string) string {
	dir := filepath.Dir(last)
	_, err := os.Stat(dir)
	if err != nil {
		return "./"
	}
	return dir
}
