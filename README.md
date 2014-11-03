# godo

[godoc](https://godoc.org/github.com/go-godo/godo)

godo is a task runner and file watcher for golang in the spirit of
rake, gulp.

To install

    go get -u gopkg.in/godo.v1/cmd/godo

## Godofile

As an example, create a file **tasks/Godofile.go** with this content

    package main

    import (
        . "gopkg.in/godo.v1"
    )

    func Tasks(p *Project) {
        Env = "GOPATH=.vendor:$GOPATH OTHER_VAR=val"

        p.Task("default", D{"hello", "build"})

        p.Task("hello", func() {
            Bash("echo Hello $USER!")
        })

        p.Task("build", W{"**/*.go"}, func() {
            Run("GOOS=linux GOARCH=amd64 go build", In{"cmd/server"})
        })

        p.Task("views", W{"templates/**/*.go.html"}, func() {
            Run("razor templates")
        })

        p.Task("server", D{"views"}, W{"**/*.go"}, Debounce(3000), func() {
            // Start recompiles and restarts on changes when watching
            Start("main.go", In{"cmd/server"})
        })
    }

    func main() {
        Godo(Tasks)
    }


To run "server" task from parent dir of `tasks/`

    godo server

To rerun "server" and its dependencies whenever any `*.go.html`  or `*.go` file changes

    godo server --watch

To run the "default" task which runs "hello" and "views"

    godo

Task names may add a "?" suffix meaning only run once even when watching

    // build once regardless of number of dependents
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

### Bash

Bash functions uses the bash executable and thus may not run on all OS.

Run a bash script string. The script can be multine line with continutation.

    Bash(`
        echo -n $USER
        echo some really long \
            command
    `)

Run a bash script string and capture its output.

    output, err := BashOutput(`echo -n $USER`)

### Run

Run `go build` inside of cmd/app and set environment variables. Notice
environment variables are set the same way as in a shell.

    Run(`GOOS=linux GOARCH=amd64 go build`)

Run and capture output

    output, err := RunOutput("whoami")


### Start

Godo tracks the pid of the `Start()` async function to restart an application gracefully.

Start an async command. If executable has suffix ".go" then it will be "go install"ed then executed.
Use this for watching a server task.

    Start("main.go", In{"cmd/app"})

### Inside

If you need to run many commands in a directory, use `Inside` instead of
the `In` options.

    Inside("somedir", func() {
        Run("...")
        Bash("...")
    })

## Godofile run-time environment

Tasks often need to run in a known evironment.

To specify whether to inherit from parent's process environment,
set `InheritParentEnv`. This setting defaults to true

    InheritParentEnv = false

To specify the base environment for your tasks, set `Env`.
Separate with whitespace or newlines.

    Env = `
        GOPATH=.vendor:$GOPATH
        PG_USER="developer"
    `

Functions can add or override environment variables as part of the command string.
Note that environment variables are set similar to how you would set them in
a shell; however, the `Run` and `Start` functions do not use a shell.

    p.Task("build", func() {
        Run("GOOS=linux GOARCH=amd64 go build" )
    })

The effective environment for `Run` or `Start` is: `parent (if inherited) <- Env <- func overrides`

The effective environment for `Bash` is: `parent (if inherited) <- Env`

Note: Interpolation of `$VARIABLE` is always from parent environment even if
`InheritParentEnv` is `false`.

