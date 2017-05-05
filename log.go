package complete

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var logger = getLogger()

func getLogger() func(format string, args ...interface{}) {
	var logfile io.Writer = ioutil.Discard
	if os.Getenv(envDebug) != "" {
		logfile = os.Stderr
	}
	return log.New(logfile, "complete ", log.Flags()).Printf
}
