package gosu

import (
	"github.com/mgutz/gosu/fsnotify"
)

// Context is the data passed to a task.
type Context struct {
	// Task is the currently running task.
	Task *Task

	// FileEvent is an event from the watcher with change details.
	FileEvent *fsnotify.FileEvent
}

// AnyFile returns either a non-DELETe FileEvent file or the WatchGlob patterns which
// can be used by goa.Load()
func (context *Context) AnyFile() []string {
	if context.FileEvent != nil && !context.FileEvent.IsDelete() {
		return []string{context.FileEvent.Name}
	}
	return context.Task.WatchGlobs
}
