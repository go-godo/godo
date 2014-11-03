package godo

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/mgutz/ansi"
	"gopkg.in/godo.v1/util"
)

var spawnedProcesses = make(map[string]*os.Process)

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
	return cmd, nil
}

func (gcmd *command) run() (output string, err error) {

	cmd, err := gcmd.toExecCmd()
	if err != nil {
		return
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
	waitExit = true
	waitgroup.Add(1)
	go func() {
		err = cmd.Start()
		if err != nil {
			return
		}
		spawnedProcesses[id] = cmd.Process
		c := make(chan error, 1)
		c <- cmd.Wait()
		_ = <-c
		waitgroup.Done()
	}()
	return nil
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
}
