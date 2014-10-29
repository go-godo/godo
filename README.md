# gosu

[godoc](https://godoc.org/github.com/mgutz/gosu)

gosu is a task runner and file watcher for golang in the spirit of
rake, gulp.

To install

    go get -u github.com/mgutz/gosu/cmd/gosu

As an example, create a file **"tasks/Gosufile.go"** with this content

    package main

    import (
        . "github.com/mgutz/gosu"
    )

    func Tasks(p *Project) {
        p.Task("default", D{"hello", "views"})

        p.Task("hello", func() {
            Run(`bash -c "echo Hello $USER!"`)
        })

        p.Task("views", "Compiles razor templates", W{"**/*.go.html"}, func(c *Context) {
            Run(`razor views`)
        })

        p.Task("server", D{"views"}, W{"**/*.go"}, Debounce(3000), func() {
            // Start recompiles and restarts on changes when watching
            Start("main.go", M{"Dir": "example"})
        })
    }

    func main() {
        Gosu(Tasks)
    }

To run "views" task from terminal

    gosu views

To rerun "views" whenever any `*.go.html` file changes

    gosu views --watch

To run the "default" task which runs "hello" and "views"

    gosu

Task names may add a "?" suffix meaning only run once even when watching

    // build once regardless of number of dependees
    p.Task("build?", func() {})

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

    func() {}           - Simple function handler
    func(c *Context) {} - Handler which accepts the current context

