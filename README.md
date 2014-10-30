# gosu

[godoc](https://godoc.org/github.com/mgutz/gosu)

gosu is a task runner and file watcher for golang in the spirit of
rake, gulp.

To install

    go get -u github.com/mgutz/gosu/cmd/gosu


BREAKING CHANGES

I apologize about the exec API changing so much. I will settle in the next
couple of days and release a v2 on gopkg.in with a promise to maintain
API compatbility for each version.

## Gosufile

As an example, create a file **tasks/Gosufile.go** with this content

    package main

    import (
        . "github.com/mgutz/gosu"
    )

    func Tasks(p *Project) {
        Env = `
            GOPATH=.vendor:$GOPATH
        `

        p.Task("default", D{"hello", "views"})

        p.Task("hello", func() {
            Bash(`
                echo Hello $USER! \
                     A beautiful day to ya
                printenv
            `)
        })

        p.Task("build", W{"**/*.go"}, func() {
            Run("GOOS=linux GOARCH=amd64 go build", In{"cmd/app"})
        })

        p.Task("server", D{"views"}, W{"**/*.go"}, Debounce(3000), func() {
            // Start recompiles and restarts on changes when watching
            Start("main.go", In{"example"})
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

Gosu provides simple exec functions. They are included as part of Gosu package as they
are frequently used in tasks. Moreover, Gosu tracks the pid of the `Start()` async function
to restart an application gracefully.

Run a bash script string. The script can be multine line with continutation.

    Bash(`
        echo -n $USER
        echo some really long \
            command
    `)

Run a bash script string and capture its output.

    output, err := BashOutput(`echo -n $USER`)

Run `go build` inside of cmd/app and set environment variables. Notice
environment variables are set the same way as in a shell.

    Run(`GOOS=linux GOARCH=amd64 go build`)

Run and capture output

    output, err := RunOutput("whoami")

Start an async command. If executable has suffix ".go" then it will be "go install"ed then executed.
Use this for watching a server task.

    Start("main.go", In{"cmd/app"})

If you need to run many commands in a directory

    Inside("somedir", func() {
        Run("...")
        Bash("...")
    })

## Gosufile run-time environment

Tasks often need to run in a known evironment.

To specify whether to inherit from parent's environment, set `InheritParentEnv`.
This setting defaults to true

    InheritParentEnv = false

To specify the base environment for your tasks, set `ENV`.
Separate with whitespace or newlines.

    Env = `
        GOPATH=.vendor:$GOPATH
    `

Funcs can add or override environment variables as part of the command string.

    p.Task("build", func() {
        Run("GOOS=linux GOARCH=amd64 go build" )
    })

The effective environment is: parent (if inherited) <- Env <- func's overrides

Note: Interpolation of `$VARIABLE` is always from parent environment.

