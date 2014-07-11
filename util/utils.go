package util

import (
	"os/exec"

	"github.com/mgutz/str"
)

// Exec is sugary way to execute a command. Exec converts cmd into an
// an executable and an argv.
func Exec(cmd string) {
	argv := str.ToArgv(cmd)
	executable := argv[0]
	argv = argv[1:]
	out, err := exec.Command(executable, argv...).CombinedOutput()
	if err != nil {
		Error("", "%s\n\n%s\n", cmd, string(out))
	}
}
