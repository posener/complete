package complete

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)


func logger() func(format string, args ...interface{}) {
	var logfile io.Writer = ioutil.Discard
	if os.Getenv(envDebug) != "" {
		logfile = os.Stderr
	}
	return log.New(logfile, "complete ", log.Flags()).Printf
}
