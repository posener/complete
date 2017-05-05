package complete

import (
	"os"
	"path/filepath"
)

// Predicate determines what terms can follow a command or a flag
type Predicate struct {
	// ExpectsNothing determine if the predicate expects something after.
	// flags/commands that do not expect any specific argument should
	// leave it on false
	ExpectsNothing bool
	// Predictor is function that returns list of arguments that can
	// come after the flag/command
	Predictor func() []Option
}

// Or unions two predicate struct, so that the result predicate
// returns the union of their predication
func (p Predicate) Or(other Predicate) Predicate {
	return Predicate{
		ExpectsNothing: p.ExpectsNothing && other.ExpectsNothing,
		Predictor:      func() []Option { return append(p.predict(), other.predict()...) },
	}
}

func (p Predicate) predict() []Option {
	if p.Predictor == nil {
		return nil
	}
	return p.Predictor()
}

var (
	PredictNothing  = Predicate{ExpectsNothing: true}
	PredictAnything = Predicate{}
)

func PredictSet(options []string) Predicate {
	return Predicate{
		Predictor: func() []Option {
			ret := make([]Option, len(options))
			for i := range options {
				ret[i] = Arg(options[i])
			}
			return ret
		},
	}
}

func PredictFiles(pattern string) Predicate {
	return Predicate{Predictor: glob(pattern)}
}

func PredictDirs(path string) Predicate {
	return Predicate{Predictor: dirs(path)}
}

func dirs(path string) func() []Option {
	return func() (options []Option) {
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
}

func glob(pattern string) func() []Option {
	return func() []Option {
		files, err := filepath.Glob(pattern)
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
