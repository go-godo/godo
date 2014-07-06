# gosu

*gosu* is a build tool for Go in the spirit of Rake, Gulp, Projmate ...
*gosu* supports watching, globbing, tasks and importing other projects.

_Asset.pipeline is in the works. Much simpler than Gulp or Grunt._

## Install

    go get github.com/mgutz/gosu

## Example

*gosu* does not have its own executable. Instead, use *gosu* to build your
project build tool.

```go
import (
    gosu"github.com/mgutz/gosu"
)

func Project(p *gosu.Project) {
    p.Task("default", "Runs all tasks" string[]{"stylesheets", "app"})

    p.Task("stylesheets", gosu.Files{"public/css/**/*.less"}, func(c *gosu.Context) {
        if c.FileEvent != nil {
            // c.FileEvent contains change event from watch
        }
        for _, f := range c.Task.WatchFiles {
            // f.FileInfo and f.Path
        }
    })

    p.Task("app", "(Re)runs the app", gosu.Files{"**/*.go"}, func() {
        // use any restart package here
    })
}

func main() {
    gosu.Run(Project)
}
```

To run default task: `go run main.go`

To run a single task:  `go run main.go stylesheets`

To run and watch a task: `go run main.go --watch stylesheets`

Build your utility:

    go build -o gosu main.go    # name it whatever you want
    gosu --watch stylesheets         # profit

## Syntax

All file patterns MUST start with a directory:

-   OK  `"./test.go"`
-   Bad `"test.go"`

### Adding tasks

`Project.Task` has variable arguments of type `interface{}` for usability.

Tasks MUST define Handler, ContextHandler or have Dependencies

To add a default task, which runs when a task on command-line is not provided

```go
p.Task("default", string[]{"clean", "stylesheets", "views"})
```

To add a task with description and Handler

```go
p.Task("name", "description", func() {
    // work here
})
```

To add a task with description and ContextHandler

```
p.Task("name", "description", func(c *gosu.Context) {
    // use context to get info about c.FileEvent or c.Task
})
```

To add a task with Dependencies

```go
p.Task("name", string[]{"dep1", "dep2"})
```

To support watching, add glob patterns

```
p.Task("views", gosu.Files{"./views/**/*.go.html"}, func() {
    // compile templates
})
```

### Glob Patterns

    /**/   - match zero or more directories
    {a,b}  - match a or b, no spaces
    *      - match any non-separator char
    ?      - match a single non-separator char
    **/    - match any directory, start of pattern only
    /**    - match any this directory, end of pattern only
    !      - removes files from resultset, start of pattern only

### Import another project

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

If you are receiving wierd events, please read [fsnotify](https://github.com/howeyc/fsnotify)


## LICENSE

The MIT License

