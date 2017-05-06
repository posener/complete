package complete

import (
	"os"
	"path/filepath"
)

// Predicate determines what terms can follow a command or a flag
type Predicate func(last string) []Option

// Or unions two predicate struct, so that the result predicate
// returns the union of their predication
func (p Predicate) Or(other Predicate) Predicate {
	if p == nil || other == nil {
		return nil
	}
	return func(last string) []Option { return append(p.predict(last), other.predict(last)...) }
}

func (p Predicate) predict(last string) []Option {
	if p == nil {
		return nil
	}
	return p(last)
}

var (
	PredictNothing Predicate = nil
)

func PredictAnything(last string) []Option { return nil }

func PredictSet(options ...string) Predicate {
	return func(last string) []Option {
		ret := make([]Option, len(options))
		for i := range options {
			ret[i] = Arg(options[i])
		}
		return ret
	}
}

func PredictDirs(last string) (options []Option) {
	dir := dirFromLast(last)
	return dirsAt(dir)
}

func dirsAt(path string) []Option {
	dirs := []string{}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if !filepath.IsAbs(path) {
		filesToRel(dirs)
	}
	return filesToOptions(dirs)
}

func PredictFiles(pattern string) Predicate {
	return func(last string) []Option {
		dir := dirFromLast(last)
		files, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			Log("failed glob operation with pattern '%s': %s", pattern, err)
		}
		if !filepath.IsAbs(pattern) {
			filesToRel(files)
		}
		return filesToOptions(files)
	}
}

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
		files[i] = "./" + rel
	}
	return
}

func filesToOptions(files []string) []Option {
	options := make([]Option, len(files))
	for i, f := range files {
		options[i] = ArgFileName(f)
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
