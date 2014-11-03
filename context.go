package godo

import (
	"github.com/go-godo/godo/watcher"
	"gopkg.in/fsnotify.v1"
)

// Context is the data passed to a task.
type Context struct {
	// Task is the currently running task.
	Task *Task

	// FileEvent is an event from the watcher with change details.
	FileEvent *watcher.FileEvent
}

// AnyFile returns either a non-DELETe FileEvent file or the WatchGlob patterns which
// can be used by goa.Load()
func (context *Context) AnyFile() []string {
	if context.FileEvent != nil && context.FileEvent.Op != fsnotify.Remove {
		return []string{context.FileEvent.Name}
	}
	return context.Task.WatchGlobs
}
