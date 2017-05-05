package complete

import (
	"os"
	"path/filepath"
)

type FlagOptions struct {
	HasFollow      bool
	FollowsOptions func() []Option
}

func (f *FlagOptions) follows() []Option {
	if f.FollowsOptions == nil {
		return nil
	}
	return f.FollowsOptions()
}

var (
	FlagNoFollow      = FlagOptions{}
	FlagUnknownFollow = FlagOptions{HasFollow: true}
)

func FlagFileFilter(pattern string) FlagOptions {
	return FlagOptions{
		HasFollow:      true,
		FollowsOptions: glob(pattern),
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
