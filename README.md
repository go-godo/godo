# gosu

[godoc](https://godoc.org/github.com/mgutz/gosu)

gosu is a task runner and file watcher for golang in the spirit of
rake, gulp.

To install

    go get -u github.com/mgutz/gosu/cmd/gosu

## Gosufile

As an example, create a file **"tasks/Gosufile.go"** with this content

    package main

    import (
        . "github.com/mgutz/gosu"
    )

    func Tasks(p *Project) {
        p.Task("default", D{"hello", "views"})

        p.Task("hello", func() {
            Bash(`
                echo Hello $USER!
                echo A beautiful day \
                    to ya
            `)
        })

        p.Task("views", "Compiles razor templates", W{"templates/**/*.go.html"}, func(c *Context) {
            Inside("templates", func() {
                Run("razor views")
            })
            Bash(`pwd`)
        })

        p.Task("server", D{"views"}, W{"**/*.go"}, Debounce(3000), func() {
            // Start recompiles and restarts on changes when watching
            Start("main.go", &Cmd{Wd: "example"})
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

## Exec functions


Gosu provides simple exec functions. They are included as part of Gosu package because they
are frequently used in tasks. Moreover, Gosu tracks the PID of the `Start()` async function
to restart an application gracefully.

Run a bash string. Can be multine line with continutation.

    Bash(`
        echo -n $USERj
        echo some really long \
            command
    `)

Run a bash string script and capture its output.

    output, err := BashOutput(`echo -n $USER`)

Run main executable inside of cmd/app and set environment var FOO

    Run("main", &Cmd{Wd: "cmd/app", Env: []string{"FOO=bar"})

Run and capture output

    output, err := RunOutput('whoami')

Start an async command. If executable has suffix ".go" then it will be "go install"ed then executed.
Use this for watching a server task.

    Start("main.go", &Cmd{Wd: "cmd/app")
