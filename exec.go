package godo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/mgutz/str"
	"gopkg.in/godo.v1/util"
)

// In is used by Bash, Run and Start to set the working directory.
// This is DEPRECATED use M{"$in": "somedir"} instead.
type In []string

// Bash executes a bash script (string) with an option to set
// the working directory.
func Bash(script string, options ...interface{}) error {
	_, err := bash(false, script, options)
	return err
}

// BashOutput is the same as Bash and it captures stdout and stderr.
func BashOutput(script string, options ...interface{}) (string, error) {
	return bash(true, script, options)
}

// Run runs a command with an an option to set the working directory.
func Run(commandstr string, options ...interface{}) error {
	_, err := run(false, commandstr, options)
	return err
}

// RunOutput is same as Run and it captures stdout and stderr.
func RunOutput(commandstr string, options ...interface{}) (string, error) {
	return run(true, commandstr, options)
}

// Start starts an async command. If executable has suffix ".go" then it will
// be "go install"ed then executed. Use this for watching a server task.
//
// If Start is called with the same command it kills the previous process.
//
// The working directory is optional.
func Start(commandstr string, options ...interface{}) error {
	m := getOptionsMap(options)
	dir, err := getWorkingDir(m)
	if err != nil {
		return nil
	}

	if strings.Contains(commandstr, "{{") {
		commandstr, err = util.StrTemplate(commandstr, m)
		if err != nil {
			return err
		}
	}

	executable, argv, env := splitCommand(commandstr)
	isGoFile := strings.HasSuffix(executable, ".go")
	if isGoFile {
		err = Run("go install -a", options...)
		if err != nil {
			return err
		}
		executable = filepath.Base(dir)
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

func getWorkingDir(m map[string]interface{}) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", nil
	}

	var wd string
	if m != nil {
		if d, ok := m["$in"].(string); ok {
			wd = d
		}
	}
	if wd != "" {
		path := filepath.Join(pwd, wd)
		_, err := os.Stat(path)
		if err == nil {
			return filepath.Join(path), nil
		}
		return "", fmt.Errorf("working dir does not exist: %s", path)
	}
	return pwd, nil
}

// Bash executes a bash string. Use backticks for multiline. To execute as shell script,
// use Run("bash script.sh")
func bash(captureOutput bool, script string, options []interface{}) (output string, err error) {
	m := getOptionsMap(options)
	dir, err := getWorkingDir(m)
	if err != nil {
		return
	}

	if strings.Contains(script, "{{") {
		script, err = util.StrTemplate(script, m)
		if err != nil {
			return "", err
		}
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

func run(captureOutput bool, commandstr string, options []interface{}) (output string, err error) {
	m := getOptionsMap(options)
	dir, err := getWorkingDir(m)
	if err != nil {
		return
	}

	if strings.Contains(commandstr, "{{") {
		commandstr, err = util.StrTemplate(commandstr, m)
		if err != nil {
			return "", err
		}
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

	s, err := cmd.run()
	return s, err
}

func getOptionsMap(args []interface{}) map[string]interface{} {
	if len(args) == 0 {
		return nil
	} else if m, ok := args[0].(map[string]interface{}); ok {
		return m
	} else if m, ok := args[0].(M); ok {
		return m
	} else if in, ok := args[0].(In); ok {
		// legacy functions used to pass in In{}
		return map[string]interface{}{"$in": in[0]}
	}
	return nil
}

// func getWd(wd []In) (string, error) {
// 	if len(wd) == 1 {
// 		return wd[0][0], nil
// 	}
// 	return os.Getwd()
// }

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

// Prompt prompts user for input with default value.
func Prompt(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text
}

// PromptPassword prompts user for password input.
func PromptPassword(prompt string) string {
	fmt.Printf(prompt)
	return string(gopass.GetPasswd())
}
