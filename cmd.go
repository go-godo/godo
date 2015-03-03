package godo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/mgutz/ansi"
	"gopkg.in/godo.v1/util"
)

// Proccesses are the processes spawned by Start()
var Processes = make(map[string]*os.Process)

type command struct {
	// original command string
	commandstr string
	// parsed executable
	executable string
	// parsed argv
	argv []string
	// parsed env
	env []string
	// working directory
	wd string
	// whether to capture output
	captureOutput bool
	// the output recorder
	recorder bytes.Buffer
}

func (gcmd *command) toExecCmd() (cmd *exec.Cmd, err error) {
	cmd = exec.Command(gcmd.executable, gcmd.argv...)
	if gcmd.wd != "" {
		cmd.Dir = gcmd.wd
	}

	cmd.Env = effectiveEnv(gcmd.env)
	cmd.Stdin = os.Stdin

	if gcmd.captureOutput {
		outWrapper := newFileWrapper(os.Stdout, &gcmd.recorder, "")
		errWrapper := newFileWrapper(os.Stderr, &gcmd.recorder, ansi.ColorCode("red+b"))
		cmd.Stdout = outWrapper
		cmd.Stderr = errWrapper
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if verbose {
		if Env != "" {
			util.Debug("#", "Env: %s\n", Env)
		}
		util.Debug("#", "%s\n", gcmd.commandstr)
	}

	return cmd, nil
}

func (gcmd *command) run() (string, error) {
	var err error
	cmd, err := gcmd.toExecCmd()
	if err != nil {
		return "", err
	}

	err = cmd.Run()
	if gcmd.captureOutput {
		return gcmd.recorder.String(), err
	}
	return "", err

}

func (gcmd *command) runAsync() (err error) {
	cmd, err := gcmd.toExecCmd()
	if err != nil {
		return
	}

	id := gcmd.commandstr

	// kills previously spawned process (if exists)
	killSpawned(id)
	waitgroup.Add(1)
	waitExit = true
	go func() {
		err = cmd.Start()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		Processes[id] = cmd.Process
		if verbose {
			util.Debug("#", "Processes[%q] added\n", id)
		}
		c := make(chan error, 1)

		c <- cmd.Wait()
		_ = <-c
		waitgroup.Done()
	}()
	return nil
}

func killSpawned(command string) {
	process := Processes[command]
	if process == nil {
		return
	}

	err := process.Kill()
	delete(Processes, command)
	if err != nil {
		util.Error("Start", "Could not kill existing process %+v\n%s\n", process, err.Error())
		return
	}
	if verbose {
		util.Debug("#", "Processes[%q] killed\n", command)
	}
}
