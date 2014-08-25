package gosu

import (
	//"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/mgutz/gosu/util"
	"github.com/mgutz/str"
)

var spawnedProcesses = make(map[string]*os.Process)

func init() {
	//	setupSignals()
}

// RunError is a simple way to execute a CLI utility.
func RunError(command string, options ...map[string]interface{}) error {
	argv := str.ToArgv(command)
	executable := argv[0]
	argv = argv[1:]
	// for _, arg := range args {
	// 	argv = append(argv, arg)
	// }
	cmd := exec.Command(executable, argv...)

	if len(options) == 1 {
		opts := options[0]
		if opts["Dir"] != nil {
			cmd.Dir = opts["Dir"].(string)
		}
		if opts["Env"] != nil {
			cmd.Env = opts["Env"].([]string)
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Run is simple way to execute a CLI utility. `command` is parsed
// for arguments. args is optional and unparsed.
func Run(command string, options ...map[string]interface{}) {
	err := RunError(command, options...)
	if err != nil {
		util.Error("Run", "%s\n%+v", command, err)
	}
}

// StartError is a simple way to start a process. If start is called with the same
// command it will kill the previous process.
func StartError(command string, options ...map[string]interface{}) error {
	argv := str.ToArgv(command)
	executable := argv[0]
	argv = argv[1:]
	// for _, arg := range args {
	// 	argv = append(argv, arg)
	// }
	cmd := exec.Command(executable, argv...)
	if len(options) == 1 {
		opts := options[0]
		if opts["Dir"] != nil {
			cmd.Dir = opts["Dir"].(string)
		}
		if opts["Env"] != nil {
			cmd.Env = opts["Env"].([]string)
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// kills previously (if any) spawned process
	killSpawned(command)

	err := cmd.Start()
	spawnedProcesses[command] = cmd.Process
	if err != nil {
		util.Error("Start", "Could not start process %s\n", command)
		return err
	}
	return nil
}

// Start is a simple way to start a process. If start is called with the same
// command it will kill the previous process.
func Start(command string, options ...map[string]interface{}) {
	err := StartError(command, options...)
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

func setupSignals() {
	sigc := make(chan os.Signal, 1)
	//signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	signal.Notify(sigc, syscall.SIGINT)
	go func() {
		<-sigc
		for command := range spawnedProcesses {
			killSpawned(command)
		}
	}()
}
