# gosu

[godoc](https://godoc.org/github.com/mgutz/gosu)

    import "github.com/mgutz/gosu"

Package gosu is a project build toolkit for Go in the spirit of Rake, Grunt and
others. Gosu supports watching, file globs, tasks and modular projects.

Gosu requires a tasks configuration function, in which task are registered and
other tasks imported.

To install

    go get -u github.com/mgutz/gosu
    go get -u github.com/mgutz/gosu/cmd/gosu

As an example, create a file 'tasks/Gosufile.go'

    package tasks

    import (
        . "github.com/mgutz/gosu"
    )

    func Tasks(p *Project) {
        p.Task("default", D{"hello, "views"})

        p.Task("hello", func() {
            Run(`bash -c "echo Hello $USER!"`)
        })

        p.Task("views", W{"**/*.go.html"}, func(c *Context) {
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

        p.Task("server", D{"views"}, W{"**/*.go}, Debounce(3000), func() {
            // DO NOT use "go run", it creates a child process that is difficult to kill
            Run("go build -o example main.go", M{"Dir": "example"})
            Start("example", M{"Dir": "example"})
        })
    }

    func main() {
        Gosu(Tasks)
    }

Task options

    D{} or Dependencies{} - dependent tasks which run before task
    Debounce              - minimum milliseconds before task can run again
    W{} or Watches{}      - array of glob file patterns to watch

        /**/   - match zero or more directories
        {a,b}  - match a or b, no spaces
        *      - match any non-separator char
        ?      - match a single non-separator char
        **/    - match any directory, start of pattern only
        /**    - match any this directory, end of pattern only
        !      - removes files from resultset, start of pattern only

Task handlers

    func() {} - Simple function handler
    func(c *Context) {} - Handler which accepts the current context

To run "views" from terminal

    gosu views

To run the "default" task which runs the dependencies "hello", "views"

    gosu

Note the "views" task specifies W{"**/*.go.html"}, which is a glob pattern to
watch any file with .go.html extension. To rerun "views" whenever any file
changes, run gosu in watch mode

    gosu --watch

