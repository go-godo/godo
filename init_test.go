package godo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// amount of time to wait before a test project is ready
var testProjectDelay = 150 * time.Millisecond
var testWatchDelay time.Duration

func init() {
	// devnull, err := os.Open(os.DevNull)
	// if err != nil {
	// 	panic(err)
	// }
	// util.LogWriter = devnull

	// WatchDelay is the time to poll the file system
	SetWatchDelay(150 * time.Millisecond)
	//testWatchDelay = watchDelay + 250*time.Millisecond
	testWatchDelay = watchDelay*2 + 50

	// Debounce should be less han watch delay
	Debounce = 100 * time.Millisecond
	verbose = false
}

// Runs a project returning the error value from the task.
func runTask(tasksFn func(*Project), name string) (*Project, error) {
	proj := NewProject(tasksFn, func(status int) {
		//fmt.Println("exited with", status)
		panic("exited with code" + strconv.Itoa(status))
	}, nil)
	return proj, proj.Run(name)
}

// Runs tasksFn with command line arguments.
func execCLI(tasksFn func(*Project), argv []string, customExitFn func(int)) int {
	var code int
	var exitFn func(code int)
	if customExitFn == nil {
		exitFn = func(status int) {
			code = status
		}
	} else {
		exitFn = func(status int) {
			code = status
			customExitFn(status)
		}
	}
	godoExit(tasksFn, argv, exitFn)
	return code
}

func touch(path string, delta time.Duration) {
	if _, err := os.Stat(path); err == nil {
		tn := time.Now()
		tn = tn.Add(delta)
		err := os.Chtimes(path, tn, tn)
		if err != nil {
			fmt.Printf("Err touching %s\n", err.Error())
		}
		//fmt.Printf("touched %s %s\n", path, time.Now())
		return
	}
	os.MkdirAll(filepath.Dir(path), 0755)
	ioutil.WriteFile(path, []byte{}, 0644)
}

func touchTil(filename string, timeout time.Duration, cquit chan bool) {
	filename, _ = filepath.Abs(filename)
forloop:
	for {
		select {
		case <-cquit:
			break forloop
		case <-time.After(timeout):
			touch(filename, 1*time.Minute)
		}
	}
}

func sliceContains(slice []string, val string) bool {
	for _, it := range slice {
		if it == val {
			return true
		}
	}
	return false
}
