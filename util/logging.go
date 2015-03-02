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
var writer io.Writer

func init() {
	ansi.DisableColors(false)
	cyan = ansi.ColorFunc("cyan")
	red = ansi.ColorFunc("red+b")
	yellow = ansi.ColorFunc("yellow+b")
	redInverse = ansi.ColorFunc("white:red")
	gray = ansi.ColorFunc("black+h")
	magenta = ansi.ColorFunc("magenta+h")
	writer = colorable.NewColorableStdout()
}

// Debug writes a debug statement to stdout.
func Debug(group string, format string, any ...interface{}) {
	fmt.Fprint(writer, gray(group)+" ")
	fmt.Fprintf(writer, gray(format), any...)
}

// Info writes an info statement to stdout.
func Info(group string, format string, any ...interface{}) {
	fmt.Fprint(writer, cyan(group)+" ")
	fmt.Fprintf(writer, format, any...)
}

// Error writes an error statement to stdout.
func Error(group string, format string, any ...interface{}) {
	fmt.Fprintf(writer, red(group)+" ")
	fmt.Fprintf(writer, red(format), any...)
}

// Panic writes an error statement to stdout.
func Panic(group string, format string, any ...interface{}) {
	fmt.Fprintf(writer, redInverse(group)+" ")
	fmt.Fprintf(writer, redInverse(format), any...)
}

// Deprecate writes a deprecation warning.
func Deprecate(message string) {
	fmt.Fprintf(writer, yellow("godo")+" "+message)
}
