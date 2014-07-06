package gosu

import (
	//"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mgutz/gosu/fsnotify"
)

// Context is the data passed to a task.
type Context struct {
	// Task is the currently running task.
	Task *Task

	// FileEvent is an event from the watcher with change details.
	FileEvent *fsnotify.FileEvent
}

// Files type is use to discern between files and dependencies when adding
// a task to the project.
type Files []string

// A Task encapsulates a handler that executes some use defined work.
type Task struct {
	Name           string
	Description    string
	Dependencies   []string
	Handler        func()
	ContextHandler func(*Context)

	// Sources are the files that need to be processed. For example `style.less`
	SourceFiles []*FileAsset
	SourceGlobs Files

	// Watches are the files are watched. On change the task is rerun. For example `**/*.less`
	// Usually Watches and Sources are the same.
	WatchFiles   []*FileAsset
	WatchGlobs   Files
	WatchRegexps []*RegexpInfo

	// Complete indicates whether this task has already ran. This flag is
	// ignored in watch mode.
	Complete bool
}

// Expands glob patterns.
func (task *Task) expandGlobs() {
	files, regexps, err := Glob(task.WatchGlobs)
	if err != nil {
		Errorf(task.Name, "%v", err)
		return
	}
	task.WatchRegexps = regexps
	task.WatchFiles = files
}

// Run runs all the dependencies of this task and when they have completed,
// runs this task.
func (task *Task) Run() {
	if !*watching && task.Complete {
		Debugf(task.Name, "Already ran\n")
		return
	}
	task.RunFromEvent(task.Name, nil)
}

// isWatchedFile determines if a FileEvent's file is a watched file
func (task *Task) isWatchedFile(e *fsnotify.FileEvent) bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}

	filename, err := filepath.Rel(cwd, e.Name)
	Debugf("task", "checking for match %s\n", filename)

	if err != nil {
		return false
	}
	matched := false
	for _, info := range task.WatchRegexps {
		if info.Negate {
			if matched {
				matched = !info.MatchString(filename)
				continue
			}
		} else if info.MatchString(filename) {
			Debugf("task", "matched %s\n", filename)
			matched = true
			continue
		}
	}
	return matched
}

// RunFromEvent runs this task when triggered from a watch.
// *e* FileEvent contains information about the file/directory which changed
// in watch mode.
func (task *Task) RunFromEvent(logName string, e *fsnotify.FileEvent) {
	start := time.Now()
	if len(task.WatchGlobs) > 0 && len(task.WatchFiles) == 0 {
		task.expandGlobs()
	}

	rebuilt := ""
	// Run this task only if the file matches watch Regexps
	if e != nil {
		rebuilt = "rebuilt "
		if !task.isWatchedFile(e) {
			return
		}
	}

	if task.Handler != nil {
		task.Handler()
	} else if task.ContextHandler != nil {
		context := &Context{Task: task}
		if e != nil {
			context.FileEvent = e
		}
		task.ContextHandler(context)
	} else if len(task.Dependencies) == 0 {
		Panicf(task.Name, "Handler, ContextHandler or Dependencies is required")
		// must be dependencies only
	}

	elapsed := time.Now().Sub(start)
	Infof(logName, "%s%vms\n", rebuilt, elapsed.Nanoseconds()/1e6)

	task.Complete = true
}
