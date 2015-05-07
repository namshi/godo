// This package adds some utility
// functions in order to handle
// logging in a fancy way.
package log

import (
	"fmt"
	golog "log"
	"os"

	"github.com/kvz/logstreamer"
	"github.com/mgutz/ansi"
)

// Returns a pair of loggers that will be
// used to log the output of remote commands
// executed via SSH.
//
// The first logger is used for the stdout, the
// second one is for the stderr.
func GetRemoteLoggers(server string) (*logstreamer.Logstreamer, *logstreamer.Logstreamer) {
	stdLogger := golog.New(os.Stdout, ansi.Color("   "+server+": ", "green"), 0)
	errLogger := golog.New(os.Stdout, ansi.Color("   "+server+": ", "red"), 0)

	logStreamerOut := logstreamer.NewLogstreamer(stdLogger, " ", false)
	logStreamerErr := logstreamer.NewLogstreamer(errLogger, " ", false)

	return logStreamerOut, logStreamerErr
}

// Prints an error message
func Err(message string, args ...string) {
	list := []interface{}{}

	for _, val := range args {
		list = append(list, val)
	}

	fmt.Printf(ansi.Color(message, "red+h"), list...)
}

// Prints an info message
func Info(message string, args ...string) {
	list := []interface{}{}

	for _, val := range args {
		val = ansi.Color(val, "blue+h")
		list = append(list, val)
	}

	fmt.Printf(message, list...)
}
