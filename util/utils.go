package util

import (
	"os/exec"
	"strings"
)

// Exec is sugary way to execute a command. Exec splits a string
// on " " so do not use for complicated args.
//
//      Exec("java -jar plantuml.jar some-file.uml")
func Exec(cmd string) {
	args := strings.Split(cmd, " ")
	executable := args[0]
	args = args[1:]
	exec.Command(executable, args...).Run()
}
