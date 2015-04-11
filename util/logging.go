package util

import (
	"fmt"
	"io"

	"github.com/mattn/go-colorable"
	"github.com/mgutz/ansi"
)

var cyan func(string) string
var red func(string) string
var yellow func(string) string
var redInverse func(string) string
var gray func(string) string
var magenta func(string) string

// LogWriter is the writer to which the logs are written
var LogWriter io.Writer

func init() {
	ansi.DisableColors(false)
	cyan = ansi.ColorFunc("cyan")
	red = ansi.ColorFunc("red+b")
	yellow = ansi.ColorFunc("yellow+b")
	redInverse = ansi.ColorFunc("white:red")
	gray = ansi.ColorFunc("black+h")
	magenta = ansi.ColorFunc("magenta+h")
	LogWriter = colorable.NewColorableStdout()
}

// Debug writes a debug statement to stdout.
func Debug(group string, format string, any ...interface{}) {
	fmt.Fprint(LogWriter, gray(group)+" ")
	fmt.Fprintf(LogWriter, gray(format), any...)
}

// Info writes an info statement to stdout.
func Info(group string, format string, any ...interface{}) {
	fmt.Fprint(LogWriter, cyan(group)+" ")
	fmt.Fprintf(LogWriter, format, any...)
}

// Error writes an error statement to stdout.
func Error(group string, format string, any ...interface{}) error {
	fmt.Fprintf(LogWriter, red(group)+" ")
	fmt.Fprintf(LogWriter, red(format), any...)
	return fmt.Errorf(format, any...)
}

// Panic writes an error statement to stdout.
func Panic(group string, format string, any ...interface{}) {
	fmt.Fprintf(LogWriter, redInverse(group)+" ")
	fmt.Fprintf(LogWriter, redInverse(format), any...)
	panic("")
}

// Deprecate writes a deprecation warning.
func Deprecate(message string) {
	fmt.Fprintf(LogWriter, yellow("godo")+" "+message)
}
