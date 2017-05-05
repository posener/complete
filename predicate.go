package complete

import (
	"os"
	"path/filepath"
)

// Predicate determines what terms can follow a command or a flag
type Predicate struct {
	// Expects determine if the predicate expects something after.
	// flags/commands that do not expect any specific argument should
	// leave it on false
	Expects bool
	// Predictor is function that returns list of arguments that can
	// come after the flag/command
	Predictor func() []Option
}

func (f *Predicate) predict() []Option {
	if f.Predictor == nil {
		return nil
	}
	return f.Predictor()
}

var (
	PredictNothing  = Predicate{Expects: false}
	PredictAnything = Predicate{Expects: true}
)

func PredictFiles(pattern string) Predicate {
	return Predicate{
		Expects:   true,
		Predictor: glob(pattern),
	}
}

func glob(pattern string) func() []Option {
	return func() []Option {
		files, err := filepath.Glob(pattern)
		if err != nil {
			logger("failed glob operation with pattern '%s': %s", pattern, err)
		}
		if !filepath.IsAbs(pattern) {
			filesToRel(files)
		}
		options := make([]Option, len(files))
		for i, f := range files {
			options[i] = ArgFileName(f)
		}
		return options
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
