package godo

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/godo.v1/util"
	"gopkg.in/godo.v1/watcher"
)

// TaskFunction is the signature of the function used to define a type.
// type TaskFunc func(string, ...interface{}) *Task
// type UseFunc func(string, interface{})

// A Task is an operation performed on a user's project directory.
type Task struct {
	Name         string
	description  string
	Dependencies []string
	Handler      Handler

	// Watches are the files are watched. On change the task is rerun. For example `**/*.less`
	// Usually Watches and Sources are the same.
	WatchFiles   []*FileAsset
	WatchGlobs   []string
	WatchRegexps []*RegexpInfo

	// computed based on dependencies
	EffectiveWatchRegexps []*RegexpInfo
	EffectiveWatchGlobs   []string

	// Complete indicates whether this task has already ran. This flag is
	// ignored in watch mode.
	Complete bool
	debounce int64
	RunOnce  bool
}

// Expands glob patterns.
func (task *Task) expandGlobs() {
	files, regexps, err := Glob(task.WatchGlobs)
	if err != nil {
		util.Error(task.Name, "%v", err)
		return
	}
	task.WatchRegexps = regexps
	task.WatchFiles = files
}

// Run runs all the dependencies of this task and when they have completed,
// runs this task.
func (task *Task) Run() {
	if !watching && task.Complete {
		util.Debug(task.Name, "Already ran\n")
		return
	}
	task.RunWithEvent(task.Name, nil)
}

// isWatchedFile determines if a FileEvent's file is a watched file
func (task *Task) isWatchedFile(e *watcher.FileEvent) bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}

	filename, err := filepath.Rel(cwd, e.Path)
	filename = filepath.ToSlash(filename)
	//util.Debug("task", "checking for match %s\n", filename)
	if err != nil {
		return false
	}
	matched := false
	for _, info := range task.EffectiveWatchRegexps {
		if info.Negate {
			if matched {
				matched = !info.MatchString(filename)
				//util.Debug("task", "negated match? %s %s\n", filename, matched)
				continue
			}
		} else if info.MatchString(filename) {
			matched = true
			//util.Debug("task", "matched %s %s\n", filename, matched)
			continue
		}
	}
	return matched
}

// RunWithEvent runs this task when triggered from a watch.
// *e* FileEvent contains information about the file/directory which changed
// in watch mode.
func (task *Task) RunWithEvent(logName string, e *watcher.FileEvent) error {
	if task.RunOnce && task.Complete {
		//util.Debug(task.Name, "Already ran\n")
		return nil
	}

	start := time.Now()
	if len(task.WatchGlobs) > 0 && len(task.WatchFiles) == 0 {
		task.expandGlobs()
		if len(task.WatchFiles) == 0 {
			util.Error("task", "\""+task.Name+"\" '%v' did not match any files\n", task.WatchGlobs)
		}
	}
	// Run this task only if the file matches watch Regexps
	rebuilt := ""
	if e != nil {
		rebuilt = "rebuilt "
		if !task.isWatchedFile(e) {
			return nil
		}
		if verbose {
			util.Debug(logName, "%s\n", e.String())
		}
	}

	var err error
	log := true
	if task.Handler != nil {
		context := Context{Task: task, Args: contextArgm}
		err = task.Handler.Handle(&context)
		if err != nil {
			return fmt.Errorf("%q: %s", logName, err.Error())
		}

	} else if len(task.Dependencies) > 0 {
		// no need to log if just dependency
		log = false
	} else {
		util.Info(task.Name, "Ignored. Task does not have a handler or dependencies.\n")
		return nil
	}

	elapsed := time.Now().Sub(start)
	if log {
		util.Info(logName, "%s%vms\n", rebuilt, elapsed.Nanoseconds()/1e6)
	}

	task.Complete = true
	return nil
}

// Debounce is minimum milliseconds before task can run again
func (task *Task) Debounce(ms int64) *Task {
	if ms > 0 {
		task.debounce = ms
	}
	return task
}

// Watch a set of glob file patterns.
func (task *Task) Watch(globs ...string) *Task {
	if len(globs) > 0 {
		task.WatchGlobs = globs
	}
	return task
}

// Description sets the description for the task.
func (task *Task) Description(desc string) *Task {
	if desc != "" {
		task.description = desc
	}
	return task
}
