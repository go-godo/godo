package gosu

import (
	"fmt"
	"github.com/mgutz/ansi"
	"runtime"
)

var cyan = ansi.ColorFunc("cyan")
var red = ansi.ColorFunc("red+b")
var redInverse = ansi.ColorFunc("white:red")
var gray = ansi.ColorFunc("black+h")
var magenta = ansi.ColorFunc("magenta+h")

func init() {
	if runtime.GOOS == "windows" {
		ansi.DisableColors(true)
	}
}

func Debugf(group string, format string, any ...interface{}) {
	fmt.Print(gray(group) + " ")
	fmt.Printf(gray(format), any...)
}

func Infof(group string, format string, any ...interface{}) {
	fmt.Print(cyan(group) + " ")
	fmt.Printf(format, any...)
}

func Errorf(group string, format string, any ...interface{}) {
	fmt.Print(red(group) + " ")
	fmt.Printf(red(format), any...)
}

func Panicf(group string, format string, any ...interface{}) {
	fmt.Print(redInverse(group) + " ")
	fmt.Printf(redInverse(format), any...)
}
