package gosu

import (
	//"log"

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
	WatchFiles []*FileAsset
	WatchGlobs Files

	// Complete indicates whether this task has already ran. This flag is
	// ignored in watch mode.
	Complete bool
}

// Expands glob patterns.
func (self *Task) expandGlobs() {
	files, err := Glob(self.WatchGlobs)
	if err != nil {
		Errorf(self.Name, "%v", err)
		return
	}
	self.WatchFiles = files
}

// Run runs all the dependencies of this task and when they have completed,
// runs this task.
func (self *Task) Run() {
	if !*watching && self.Complete {
		Debugf(self.Name, "Already ran\n")
		return
	}
	self.RunFromEvent(nil)
}

// Run runs this task when triggered from a watch.
// *e* contains information about the file/directory which changed when
// watching.
func (self *Task) RunFromEvent(e *fsnotify.FileEvent) {
	if len(self.WatchGlobs) > 0 && len(self.WatchFiles) == 0 {
		self.expandGlobs()
	}

	if self.Handler != nil {
		self.Handler()
	} else if self.ContextHandler != nil {
		context := &Context{Task: self}
		if e != nil {
			context.FileEvent = e
		}
		self.ContextHandler(context)
	} else if len(self.Dependencies) == 0 {
		Panicf(self.Name, "Handler, ContextHandler or Dependencies is required")
		// must be dependencies only
	}

	self.Complete = true
}
