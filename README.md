# gosu

[godoc](https://godoc.org/github.com/mgutz/gosu)

    import "github.com/mgutz/gosu"

Package gosu is a project build toolkit for Go in the spirit of Rake, Grunt and
others. Gosu supports watching, file globs, tasks and modular projects.

Gosu requires a tasks configuration function, in which task are registered and
other tasks imported.

To install

    go get -u github.com/mgutz/gosu/gosu

As an example, create a file 'tasks/Gosufile.go'

    package tasks

    import (
        . "github.com/mgutz/gosu"
    )

    func Tasks(p *Project) {
        p.Task("default", Pre{"hello, "views"})

        p.Task("hello", func() {
            util.Exec(`bash -c "echo Hello $USER!"`)
        })

        p.Task("views", Watch{"**/*.go.html"}, func(c *Context) {
            if c.FileEvent == nil {
                for _, f := range c.Task.WatchFiles {
                    // f.FileInfo and f.Path
                    fmt.Printf("File: %s\n", f.Path)
                }
            } else {
                // change event when watching
                fmt.Printf("%v\n", c.FileEvent)
            }
        })
    }

To run "views"

    gosu views

To run the "default" task which runs the dependencies "hello", "views"

    gosu

Note the "views" task specifies "**/*.go.html", which is a glob pattern to watch
any file with .go.html extension. To rerun "views" whenever any file changes,
run gosu in watch mode

    gosu --watch

## Usage

#### func  Glob

```go
func Glob(patterns []string) ([]*FileAsset, []*RegexpInfo, error)
```
Glob returns files and dirctories that match patterns. Patterns must use
slashes, even Windows.

Special chars.

    /**/   - match zero or more directories
    {a,b}  - match a or b, no spaces
    *      - match any non-separator char
    ?      - match a single non-separator char
    **/    - match any directory, start of pattern only
    /**    - match any this directory, end of pattern only
    !      - removes files from resultset, start of pattern only

#### func  Globexp

```go
func Globexp(glob string) *regexp.Regexp
```
Globexp builds a regular express from from extended glob pattern and then
returns a Regexp object from the pattern.

#### func  Run

```go
func Run(tasksFunc func(*Project))
```
Run runs a project of tasks.

#### type Context

```go
type Context struct {
	// Task is the currently running task.
	Task *Task

	// FileEvent is an event from the watcher with change details.
	FileEvent *watcher.FileEvent
}
```

Context is the data passed to a task.

#### func (*Context) AnyFile

```go
func (context *Context) AnyFile() []string
```
AnyFile returns either a non-DELETe FileEvent file or the WatchGlob patterns
which can be used by goa.Load()

#### type FileAsset

```go
type FileAsset struct {
	os.FileInfo
	// Path to asset
	Path string
}
```

FileAsset contains file information and path from globbing.

#### type Files

```go
type Files []string
```

Files type is use to discern between files and dependencies when adding a task
to the project.

#### type M

```go
type M map[string]interface{}
```

M is generic string to interface alias

#### type Pre

```go
type Pre []string
```

Pre are dependencies which are run before a task.

#### type Project

```go
type Project struct {
	Tasks     map[string]*Task
	Namespace map[string]*Project
}
```

Project is a container for tasks.

#### func  NewProject

```go
func NewProject(tasksFunc func(*Project)) *Project
```
NewProject creates am empty project ready for tasks.

#### func (*Project) Define

```go
func (project *Project) Define(fn func(*Project))
```
Define defines tasks

#### func (*Project) Run

```go
func (project *Project) Run(name string)
```
Run runs a task by name.

#### func (*Project) Task

```go
func (project *Project) Task(name string, args ...interface{}) *Task
```
Task adds a task to the project.

#### func (*Project) Usage

```go
func (project *Project) Usage()
```
Usage prints usage about the app and tasks.

#### func (*Project) Use

```go
func (project *Project) Use(namespace string, tasksFunc func(*Project))
```
Use uses another project's task within a namespace.

#### func (*Project) Watch

```go
func (project *Project) Watch(names []string)
```
Watch watches the Files of a task and reruns the task on a watch event. Any
direct dependency is also watched.

#### type RegexpInfo

```go
type RegexpInfo struct {
	*regexp.Regexp
	Negate bool
}
```

RegexpInfo contains additional info about the Regexp created by a glob pattern.

#### type Task

```go
type Task struct {
	Name           string
	Description    string
	Dependencies   []string
	Handler        func()
	ContextHandler func(*Context)

	// Watches are the files are watched. On change the task is rerun. For example `**/*.less`
	// Usually Watches and Sources are the same.
	WatchFiles   []*FileAsset
	WatchGlobs   []string
	WatchRegexps []*RegexpInfo

	// Complete indicates whether this task has already ran. This flag is
	// ignored in watch mode.
	Complete bool
}
```

A Task is an operation performed on a user's project directory.

#### func (*Task) Run

```go
func (task *Task) Run()
```
Run runs all the dependencies of this task and when they have completed, runs
this task.

#### func (*Task) RunWithEvent

```go
func (task *Task) RunWithEvent(logName string, e *watcher.FileEvent)
```
RunWithEvent runs this task when triggered from a watch. *e* FileEvent contains
information about the file/directory which changed in watch mode.

#### type Watch

```go
type Watch []string
```

Watch type defines the glob patterns to use for watching.
