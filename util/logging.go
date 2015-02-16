package util

import (
	"fmt"
	"runtime"

	"github.com/mgutz/ansi"
)

var cyan = ansi.ColorFunc("cyan")
var red = ansi.ColorFunc("red+b")
var yellow = ansi.ColorFunc("yellow+b")
var redInverse = ansi.ColorFunc("white:red")
var gray = ansi.ColorFunc("black+h")
var magenta = ansi.ColorFunc("magenta+h")

func init() {
	if runtime.GOOS == "windows" {
		ansi.DisableColors(true)
	}
}

// Debug writes a debug statement to stdout.
func Debug(group string, format string, any ...interface{}) {
	fmt.Print(gray(group) + " ")
	fmt.Printf(gray(format), any...)
}

// Info writes an info statement to stdout.
func Info(group string, format string, any ...interface{}) {
	fmt.Print(cyan(group) + " ")
	fmt.Printf(format, any...)
}

// Error writes an error statement to stdout.
func Error(group string, format string, any ...interface{}) {
	fmt.Printf(red(group) + " ")
	fmt.Printf(red(format), any...)
}

// Panic writes an error statement to stdout.
func Panic(group string, format string, any ...interface{}) {
	fmt.Printf(redInverse(group) + " ")
	fmt.Printf(redInverse(format), any...)
}

// Deprecate writes a deprecation warning.
func Deprecate(message string) {
	fmt.Printf(yellow("godo") + " " + message)
}
