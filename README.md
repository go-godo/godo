# godo

[godoc](https://godoc.org/gopkg.in/godo.v1)

godo is a task runner and file watcher for golang in the spirit of
rake, gulp.

To install

	go get -u gopkg.in/godo.v1/cmd/godo

godo requires that the go tool and compiler installed and assumes that a valid `GOPATH` is set.

## Godofile

**godo** runs tasks defined in `Gododir`. The `Gododir` can contain multiple `.go` files which should form a valid go package.

_The legacy `tasks/Godofile.go` is still supported, but not recommended._

As an example, create a file **Gododir/myTasks.go** with this content
```go
package tasks

import (
	. "gopkg.in/godo.v1"
)

func Tasks(p *Project) {
	Env = "GOPATH=.vendor::$GOPATH PG_PASSWORD=dev"

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
```

To run "server" task from parent dir of `Gododir`

	godo server

To rerun "server" and its dependencies whenever any `*.go.html`  or `*.go` file changes

	godo server --watch

To run the "default" task which runs "hello" and "views"

	godo

Task names may add a "?" suffix to execute only once even when watching

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
		/**    - match any in this directory, end of pattern only
		!      - removes files from resultset, start of pattern only

Task handlers

	func() {}           - Simple function handler
	func(c *Context) {} - Handler which accepts the current context

## Exec functions

### Bash

Bash functions uses the bash executable and may not run on all OS.

Run a bash script string. The script can be multine line with continutation.

	Bash(`
		echo -n $USER
		echo some really long \
			command
	`)

Run a bash script and capture STDOUT and STDERR.

	output, err := BashOutput(`echo -n $USER`)

### Run

Run `go build` inside of cmd/app and set environment variables.

	Run(`GOOS=linux GOARCH=amd64 go build`, In{"cmd/app"})

Run and capture STDOUT and STDERR

	output, err := RunOutput("whoami")

### Start

Start an async command. If the executable has suffix ".go" then it will be "go install"ed then executed.
Use this for watching a server task.

	Start("main.go", In{"cmd/app"})

Godo tracks the pid of the `Start` async function to restart an application gracefully.

### Inside

To run many commands inside a directory, use `Inside` instead of the `In` option.
`Inside` changes the working directory.

	Inside("somedir", func() {
		Run("...")
		Bash("...")
	})

## Godofile Run-Time Environment

To specify whether to inherit from parent's process environment,
set `InheritParentEnv`. This setting defaults to true

	InheritParentEnv = false

To specify the base environment for your tasks, set `Env`.
Separate with whitespace or newlines.

	Env = `
		GOPATH=.vendor::$GOPATH
		PG_USER="developer"
	`

Functions can add or override environment variables as part of the command string.
Note that environment variables are set before the executable similar to a shell;
however, the `Run` and `Start` functions do not use a shell.

	p.Task("build", func() {
		Run("GOOS=linux GOARCH=amd64 go build" )
	})

The effective environment for exec functions is: `parent (if inherited) <- Env <- func parsed env`

Paths should use `::` as a cross-platform path list separator. On Windows `::` is replaced with `;`.
On Mac and linux `::` is replaced with `:`.

