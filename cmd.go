package gosu

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/mgutz/gosu/util"
	"github.com/mgutz/str"
)

var spawnedProcesses = make(map[string]*os.Process)

// Run is simple way to execute a CLI utility. `command` is parsed
// for arguments. args is optional and unparsed.
func Run(command string, options ...map[string]interface{}) (string, error) {
	return StartAsync(false, command, options...)
}

func mapToEnv(m map[string]string) []string {
	env := make([]string, len(m))

	i := 0
	for k, v := range m {
		env[i] = k + "=" + v
		i++
	}
	return env
}

func mergeEnv(pairs []string) []string {
	m := map[string]string{}

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		m[pair[0]] = pair[1]
	}

	for _, kv := range pairs {
		pair := strings.Split(kv, "=")
		// ignore non key=value strings
		if len(pair) == 2 {
			m[pair[0]] = pair[1]
		}
	}

	return mapToEnv(m)
}

// StartAsync starts a process async or sync based on the first flag. If it is an async
// operation the process is tracked and killed if started again.
func StartAsync(isAsync bool, command string, options ...map[string]interface{}) (output string, err error) {
	//existing := spawnedProcesses[command]
	argv := str.ToArgv(command)
	executable := argv[0]

	isGoFile := strings.HasSuffix(executable, ".go")
	if isGoFile {
		// install the executable which compiles files
		_, err = StartAsync(false, "go install "+executable, options...)
		if err != nil {
			return
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		util.Error("Start", "Could not get work directory\n")
		return
	}

	var childEnv []string
	if len(options) == 1 {
		opts := options[0]
		if opts["Dir"] != nil {
			wd = opts["Dir"].(string)
		}

		vars, ok := opts["Env"].([]string)
		if ok {
			childEnv = mergeEnv(vars)
		}
	}
	if isGoFile {
		executable = path.Base(wd)
	}

	argv = argv[1:]
	cmd := exec.Command(executable, argv...)
	cmd.Dir = wd
	if childEnv != nil {
		cmd.Env = childEnv
	}

	cmd.Stdin = os.Stdin
	if isAsync {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// kills previously spawned process (if exists)
		killSpawned(command)
		waitExit = true
		waitgroup.Add(1)
		go func() {
			err = cmd.Start()
			if err != nil {
				return
			}
			spawnedProcesses[command] = cmd.Process
			c := make(chan error, 1)
			c <- cmd.Wait()
			_ = <-c
			waitgroup.Done()
		}()
		return "", nil
	}

	var recorder bytes.Buffer
	outWrapper := newFileWrapper(os.Stdout, &recorder, "")
	errWrapper := newFileWrapper(os.Stderr, &recorder, ansi.ColorCode("red+b"))
	cmd.Stdout = outWrapper
	cmd.Stderr = errWrapper
	err = cmd.Run()
	return recorder.String(), err
}

// Start is a simple way to start a process or go file. If start is called with the same
// command it kills the previous process.
func Start(command string, options ...map[string]interface{}) {
	_, err := StartAsync(true, command, options...)
	if err != nil {
		util.Error("Start", "%s\n%+v\n", command, err)
	}
}

func toInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return result
}

func killSpawned(command string) {
	process := spawnedProcesses[command]
	if process == nil {
		return
	}

	err := process.Kill()
	if err != nil {
		util.Error("Start", "Could not kill existing process %+v\n", process)
		return
	}
	// _, err = process.Wait()
	// if err != nil {
	// 	util.Error("Start", "Error waiting %v\n", err)
	// 	return
	// }
	// waitgroup.Done()
}
