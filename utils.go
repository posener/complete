package complete

import (
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// fixPathForm changes a file name to a relative name in accordance with isExpandLast which if true,
// will have the form of absolute path for cases like `~` which are not recognized by anyone other than
// the shell
func fixPathForm(last string, isExpandLast bool, file string) string {
	// get wording directory for relative name
	workDir, err := os.Getwd()
	if err != nil {
		return file
	}

	if strings.Contains(last, "~") {
		if isExpandLast {
			path, err := homedir.Expand(last)
			if err != nil {
				path = last
			}
			return fixDirPath(path)
		} else {
			// else here although not required, is added here to stress the orthogonality
			// nature of this block with the block above
			return unResolveHome(last, file)
		}
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

func unResolveHome(last string, path string) string {
	if strings.Contains(last, "~") {

		// Resolve `~` to complete resolvable path
		lastEq, _ := homedir.Expand(last)
		// Get parent equivalent for cases of partially complete last
		lastParentEq, _ := homedir.Expand(filepath.Dir(lastEq))

		/*
			The following might be the possible cases to unresolve absolute path to path completion of path with `~`:
			* Case when the last as entered in terminal is partial
				** This is when, parent dir of incompletely entered last, is same as parent dir of possible path
				   Because, possible path is in this case expected to have contents after resolving the incomplete last on terminal
			* Case when the last as entered in terminal is complete
				** Resolved form of complete last is equal to path
				   This is possible because we include the same directory also as one of the possible auto-completion solutions
				** Resolved form of complete last is equal to parent of resolved auto-completion path
				   This happens because normally, except for the last itself in the auto-complete suggestions,
				   all other paths are expected to be only 1 step ahead of last at any given auto-completion step

			According to the cases above, we replace the maximum possible path out of the obtained from passed potential
			auto-completion path with the maximum available form of last so that we retain `~` in the suggestions provided
		*/
		if fixDirPath(lastEq) == fixDirPath(path) {
			path = strings.Replace(path, fixDirPath(path), fixDirPath(last), 1)
		} else if fixDirPath(lastEq) == fixDirPath(filepath.Dir(path)) {
			path = strings.Replace(path, fixDirPath(lastEq), fixDirPath(last), 1)
		} else if fixDirPath(lastParentEq) == fixDirPath(filepath.Dir(path)) {
			path = strings.Replace(path, fixDirPath(filepath.Dir(lastEq)), fixDirPath(filepath.Dir(last)), 1)
		}
	}
	return path
}

func fixDirPath(path string) string {
	tmpPath, _ := homedir.Expand(path)
	info, err := os.Stat(tmpPath)
	if err == nil && info.IsDir() && !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}
