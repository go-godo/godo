package godo

import "gopkg.in/godo.v1/watcher"

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
	if context.FileEvent != nil && context.FileEvent.Event != watcher.DELETED {
		return []string{context.FileEvent.Path}
	}
	return context.Task.WatchGlobs
}
