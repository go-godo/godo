# gosu

*gosu* is a build tool for Go in the spirit of Rake, Gulp, Projmate ...
*gosu* supports watching, globbing, tasks and modular projects.

_Asset.pipeline is at [goa](http://github.com/mgutz/goa)_

## Install

```sh
go get github.com/mgutz/gosu/gosu
```

## Example

Create a file `tasks/Gosufile.go`. Go is strict about package main and having multiple
packages in the same dir. Learn to compromise and Go loves you back.

```go
package tasks

import (
    "fmt"
    "github.com/mgutz/gosu"
    "github.com/mgutz/gosu/util"
)

func Project(p *gosu.Project) {
    p.Task("default", []string{"files", "hello"})

    p.Task("hello", func() {
        util.Exec(`bash -c "echo Hello $USER!"`)
    })

    p.Task("files", gosu.Files{"**/*"}, func(c *gosu.Context) {
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
```

To run default task: `gosu`

To run a single task:  `gosu files`

To run and watch a task: `gosu files --watch`

To see all tasks: `gosu --help`

To build your own utility

1.  Project must be inside `package main`
2.  Add this code

    fun main() {
        gosu.Run(Project)
    }


## Syntax

### Adding tasks

`Project.Task` has variable arguments of type `interface{}` for convenience.

```go
func (project *Project) Task(name string, args ...interface{})
```

To add a default task, which runs when a task name is not provided on the command line.
The best practice is to use the "default" task to define the most frequently used
dependencies. Avoid defining a handler for "default"

```go
p.Task("default", []string{"clean", "stylesheets", "views"})
```

To add a task with description and Handler

```go
// description is displayed in the Tasks help screen
p.Task("name", "description", func() {
    // work here
})
```

To add a task with description and ContextHandler

```go
p.Task("name", "description", func(c *gosu.Context) {
    // use context to get info about c.FileEvent or c.Task
})
```

To add a task with Dependencies only

```go
// run dep1, dep2, name in sequence
p.Task("name", []string{"dep1", "dep2"})
```

To enable watching on a task, add glob patterns for the files to be watched

```go
// watches all files ending with "go.html"
p.Task("views", gosu.Files{"./views/**/*.go.html"}, func() {
    // ...
})
```

All tasks MUST have a Handler, ContextHandler or Dependencies.

### Glob Patterns

```
/**/   - match zero or more directories
{a,b}  - match a or b, no spaces
*      - match any non-separator char
?      - match a single non-separator char
**/    - match any directory, start of pattern only
/**    - match any this directory, end of pattern only
!      - removes files from resultset, start of pattern only
```

### Import Project

A large project can be broken into multiple projects or projects can be
imported from other packages. Imported projects MUST be namespaced to avoid
conflicts with tasks in your project.

```go
import (
    "github.com/acme/project"
)

func Project(p *Project) {
    // Use  it within this project and assign namespace "ns"
    p.Use("ns", project.Project)

    // Add as dependency, note the namespace
    p.Task("default", []string{"ns:sprite"})
}
```

## FAQ

If you are receiving weird events (Mac Users) please read [fsnotify](https://github.com/howeyc/fsnotify) FAQ

## LICENSE

The MIT License

