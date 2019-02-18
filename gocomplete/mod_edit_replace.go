package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/posener/complete"
)

var modEditReplaceRegexp = regexp.MustCompile(strings.Replace(strings.Replace(strings.Replace(modE_R_RE, " ", "", -1), "\n", "", -1), "\t", "", -1))
var modE_R_RE = `(?:-replace=)?
(?:
	(?P<old>[^@=]+)
	(?:
		(?:(@)(?P<v_old>[^=]+))?
		(?:=
			(?:
				(?P<new>[^@]+)
				(?:(@)(?P<v_new>.+)?)?
			)?
		)?
	)? 
)?`

type replaceArgs struct {
	Old   string
	OldAt bool
	OldV  string
	Eq    bool
	New   string
	NewAt bool
	NewV  string
}

func (rs replaceArgs) String() string {
	s := rs.Old
	if rs.OldV != "" {
		s += "@" + rs.OldV
	}
	s += "=" + rs.New
	if rs.NewV != "" {
		s += "@" + rs.NewV
	}
	return s
}

func parseReplace(last string) (ret replaceArgs) {
	defer func() {
		if x := recover(); x != nil {
			ret = replaceArgs{}
		}
	}()
	if !modEditReplaceRegexp.MatchString(last) {
		panic("")
	}
	subs := modEditReplaceRegexp.FindAllStringSubmatch(last, 5)
	return replaceArgs{
		Old:   subs[0][1],
		OldAt: subs[0][2] != "",
		OldV:  subs[0][3],
		Eq:    strings.Contains(last, "="),
		New:   subs[0][4],
		NewAt: subs[0][5] != "",
		NewV:  subs[0][6],
	}
}

func sliceContains(src []string, match string) bool {
	for _, st := range src {
		if st == match {
			return true
		}
	}
	return false
}

// var modEditReplaceRegexp = regexp.MustCompile("-replace=[^@=]+?P<v_old>@[^=]+)=)?) )?")

// from go help mod edit
// The -replace=old[@v]=new[@v] and -dropreplace=old[@v] flags
// add and drop a replacement of the given module path and version pair.
// If the @v in old@v is omitted, the replacement applies to all versions
// with the old module path. If the @v in new@v is omitted, the new path
// should be a local module root directory, not a module path.
// Note that -replace overrides any existing replacements for old[@v].
func predictModEditReplace(a complete.Args) (prediction []string) {
	fmt.Fprintf(os.Stderr, "XXXX %+v\n", a)
	if sliceContains(a.Completed, "-replace") {
		cline := strings.Split(os.Getenv("COMP_LINE"), " ")
		rArgs := parseReplace(cline[len(cline)-1])
		if rArgs.Old == "" || (rArgs.OldV == "" && !rArgs.Eq) {
			return rArgs.old()
		}
		if rArgs.New == "" && rArgs.Eq {
			return rArgs.newPart()
		}
	}
	return complete.PredictFiles("go.mod").Predict(a)
}

func (rArgs *replaceArgs) old() []string {
	fmt.Fprintf(os.Stderr, "YYYY %#+v\n", rArgs)
	cmd := exec.Command("go", "mod", "edit", "-json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		complete.Log("go mod error %v\n", err)
		return nil
	}
	gomod := GoMod{}
	err = json.Unmarshal(out, &gomod)
	if err != nil {
		complete.Log("go mod error %v\n%v\n", string(out), err)
		return nil
	}
	req := make([]string, 0)
	for i := range gomod.Require {
		req = append(req, gomod.Require[i].Path+"=")
		if rArgs.OldV == "" {
			req = append(req, gomod.Require[i].Path+"@")
		}
	}
	return req
}

func (rArgs *replaceArgs) newPart() []string {
	return nil
}

type Module struct {
	Path    string
	Version string
}

type GoMod struct {
	Module  Module
	Require []Require
	Exclude []Module
	Replace []Replace
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

type Replace struct {
	Old Module
	New Module
}
