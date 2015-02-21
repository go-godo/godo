package godo

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/mgutz/str"
)

// In is used by Bash, Run and Start to set the working directory
type In []string

// Bash executes a bash script (string) with an option to set
// the working directory.
func Bash(script string, wd ...In) error {
	_, err := bash(false, script, wd)
	return err
}

// BashOutput is the same as Bash and it captures stdout and stderr.
func BashOutput(script string, wd ...In) (string, error) {
	return bash(true, script, wd)
}

// Run runs a command with an an option to set the working directory.
func Run(commandstr string, wd ...In) error {
	_, err := run(false, commandstr, wd)
	return err
}

// RunOutput is same as Run and it captures stdout and stderr.
func RunOutput(commandstr string, wd ...In) (string, error) {
	return run(true, commandstr, wd)
}

// Start starts an async command. If executable has suffix ".go" then it will
// be "go install"ed then executed. Use this for watching a server task.
//
// If Start is called with the same command it kills the previous process.
//
// The working directory is optional.
func Start(commandstr string, wd ...In) error {
	dir, err := getWd(wd)
	if err != nil {
		return err
	}

	executable, argv, env := splitCommand(commandstr)
	isGoFile := strings.HasSuffix(executable, ".go")
	if isGoFile {
		err = Run("go install", wd...)
		if err != nil {
			return err
		}
		executable = path.Base(dir)
	}

	cmd := &command{
		executable: executable,
		wd:         dir,
		env:        env,
		argv:       argv,
		commandstr: commandstr,
	}
	return cmd.runAsync()
}

// Bash executes a bash string. Use backticks for multiline. To execute as shell script,
// use Run("bash script.sh")
func bash(captureOutput bool, script string, wd []In) (output string, err error) {
	dir, err := getWd(wd)
	if err != nil {
		return
	}

	gcmd := &command{
		executable:    "bash",
		argv:          []string{"-c", script},
		wd:            dir,
		captureOutput: captureOutput,
		commandstr:    script,
	}

	return gcmd.run()
}

func run(captureOutput bool, commandstr string, wd []In) (output string, err error) {
	dir, err := getWd(wd)
	if err != nil {
		return
	}
	executable, argv, env := splitCommand(commandstr)

	cmd := &command{
		executable:    executable,
		wd:            dir,
		env:           env,
		argv:          argv,
		captureOutput: captureOutput,
		commandstr:    commandstr,
	}

	return cmd.run()
}

func getWd(wd []In) (string, error) {
	if len(wd) == 1 {
		return wd[0][0], nil
	}
	return os.Getwd()
}

func splitCommand(command string) (executable string, argv, env []string) {
	argv = str.ToArgv(command)
	for i, item := range argv {
		if strings.Contains(item, "=") {
			if env == nil {
				env = []string{item}
				continue
			}
			env = append(env, item)
		} else {
			executable = item
			argv = argv[i+1:]
			return
		}
	}

	executable = argv[0]
	argv = argv[1:]
	return
}

func toInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return result
}

// Inside temporarily changes the working directory and restores it when lambda
// finishes.
func Inside(dir string, lambda func()) error {
	olddir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	defer func() {
		os.Chdir(olddir)
	}()
	lambda()
	return nil
}
