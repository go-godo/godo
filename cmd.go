package gosu

import (
	//"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/mgutz/gosu/util"
	"github.com/mgutz/str"
)

var spawnedProcesses = make(map[string]*os.Process)

// Run is simple way to execute a CLI utility. `command` is parsed
// for arguments. args is optional and unparsed.
func Run(command string, options ...map[string]interface{}) {
	err := StartAsync(false, command, options...)
	if err != nil {
		util.Error("Run", "%s\n%+v", command, err)
	}
}

// StartAsync starts a process async or sync based on the first flag. If it is an async
// operation the process is tracked and killed if started again.
func StartAsync(isAsync bool, command string, options ...map[string]interface{}) error {
	existing := spawnedProcesses[command]
	argv := str.ToArgv(command)
	executable := argv[0]
	isGoFile := strings.HasSuffix(executable, ".go")
	if isGoFile {
		// install the executable which compiles files
		err := StartAsync(false, "go install", options...)
		if err != nil {
			return err
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		util.Error("Start", "Could not get work directory\n")
		return err
	}
	var env []string
	if len(options) == 1 {
		opts := options[0]
		if opts["Dir"] != nil {
			wd = opts["Dir"].(string)
		}
		if opts["Env"] != nil {
			env = opts["Env"].([]string)
		}
	}
	if isGoFile {
		executable = path.Base(wd)
	}

	argv = argv[1:]
	cmd := exec.Command(executable, argv...)
	cmd.Dir = wd
	if len(env) > 0 {
		cmd.Env = env
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if isAsync {
		// kills previously spawned process (if exists)
		killSpawned(command)
		err = cmd.Start()
		spawnedProcesses[command] = cmd.Process
	} else {
		err = cmd.Run()
	}
	if err != nil {
		util.Error("Start", "Could not start process %s\n", command)
		return err
	}

	if isAsync && existing == nil {
		waitgroup.Add(1)
	}
	return nil
}

// Start is a simple way to start a process or go file. If start is called with the same
// command it kills the previous process.
func Start(command string, options ...map[string]interface{}) {
	err := StartAsync(true, command, options...)
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
	_, err = process.Wait()
	if err != nil {
		util.Error("Start", "Error waiting %v\n", err)
		return
	}
}
