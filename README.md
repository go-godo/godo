# godo

[godoc](https://godoc.org/gopkg.in/godo.v1)

godo is a task runner and file watcher for golang in the spirit of
rake, gulp.

To install

    go get -u gopkg.in/godo.v1/cmd/godo

## Godofile

**godo** runs either `Gododir/Godofile.go` or `tasks/Godofile.go`.

As an example, create a file **Gododir/Godofile.go** with this content

```go
package main

import (
	"fmt"
    . "gopkg.in/godo.v1"
)

func tasks(p *Project) {
    Env = `GOPATH=.vendor::$GOPATH`

    p.Task("default", D{"hello", "build"})

    p.Task("hello", func(c *Context) {
        name := c.Args.ZeroString("name", "n")
        if name == "" {
            Bash("echo Hello $USER!")
        } else {
            fmt.Println("Hello", name)
        }
    })

    p.Task("assets?", func() {
        // The "?" tells Godo to run this task ONLY ONCE regardless of
        // how many tasks depend on it. In this case watchify watches
        // on its own.
		Run("watchify public/js/index.js d -o dist/js/app.bundle.js")
    }).Watch("public/**/*.{css,js,html}")

    p.Task("build", D{"views", "assets"}, func() error {
        return Run("GOOS=linux GOARCH=amd64 go build", In{"cmd/server"})
    }).Watch("**/*.go")

    p.Task("server", D{"views", "assets"}, func() {
        // rebuilds and restarts when a watched file changes
        Start("main.go", M{"$in": "cmd/server"})
    }).Watch("server/**/*.go", "cmd/server/*.{go,json}").
       Debounce(3000)

    p.Task("views", func() error {
        return Run("razor templates")
    }).Watch("templates/**/*.go.html")
}

func main() {
    Godo(tasks)
}
```

To run "server" task from parent dir of `tasks/`

    godo server

To rerun "server" and its dependencies whenever any of their watched files change

    godo server --watch

To run the "default" task which runs "hello" and "build"

    godo

Task names may add a "?" suffix to execute only once even when watching

```go
// build once regardless of number of dependents
p.Task("assets?", func() {})
```

Task options

    D{} or Dependencies{} - dependent tasks which run before task
    W{} or Watches{}      - array of glob file patterns to watch

        /**/   - match zero or more directories
        {a,b}  - match a or b, no spaces
        *      - match any non-separator char
        ?      - match a single non-separator char
        **/    - match any directory, start of pattern only
        /**    - match any in this directory, end of pattern only
        !      - removes files from result set, start of pattern only

Task handlers

    func()                  - Simple function handler, don't care about return
    func() error            - Simple function handler
    func(c *Context)        - Task with context, don't care about return
    func(c *Context) error  - Task with context

Any error return in task or its dependencies stops the pipeline and
`godo` exits with status code of 1, except in watch mode.

### Task Arguments

Task arguments follow POSIX style flag convention
(unlike go's built-in flag package). Any command line arguments
succeeding `--` are passed to each task. Note, arguments before `--`
are reserved for `godo`.

As an example,

```go
p.Task("hello", func(c *Context) {
    // "(none)" is the default value
    msg := c.Args.MayString("(none)", "message", "msg", "m")
    var name string
    if len(c.Args.Leftover() == 1) {
        name = c.Args.Leftover()[0]
    }
    fmt.Println(msg, name)
})
```

running

```sh
# prints "(none)"
godo hello

# prints "Hello dude" using POSIX style flags
godo hello -- dude --message Hello
godo hello -- dude --msg Hello
godo hello -- -m Hello dude
```

Args functions are categorized as

*  `Must*` - Argument must be set by user or panic.

    ```go
c.Args.MustInt("number", "n")
```

*  `May*` - If argument is not set, default to first value.

    ```go
    // defaults to 100
    c.Args.MayInt(100, "number", "n")
```

*  `Zero*` - If argument is not set, default to zero value.

    ```go
// defaults to 0
c.Args.ZeroInt("number", "n")
```

## godobin

`godo` compiles `Godofile.go` to `godobin-VERSION` (`godobin-VERSION.exe` on Windows) whenever
`Godofile.go` changes. The binary file is built into the same directory as
`Godofile.go` and should be ignored by adding the path `godobin*` to `.gitignore`.

## Exec functions

All of these functions accept a `map[string]interface{}` or `M` for
options. Option keys that start with `"$"` are reserved for `godo`.
Other fields can be used as context for template.

### Bash

Bash functions uses the bash executable and may not run on all OS.

Run a bash script string. The script can be multiline line with continutation.

```go
Bash(`
    echo -n $USER
    echo some really long \
        command
`)
```

Bash can use Go templates

```go
Bash(`echo -n {{.name}}`, M{"name": "mario", "$in": "cmd/bar"})
```

Run a bash script and capture STDOUT and STDERR.

```go
output, err := BashOutput(`echo -n $USER`)
```

### Run

Run `go build` inside of cmd/app and set environment variables.

```go
Run(`GOOS=linux GOARCH=amd64 go build`, M{"$in": "cmd/app"})
```

Run can use Go templates

```go
Run(`echo -n {{.name}}`, M{"name": "mario", "$in": "cmd/app"})
```

Run and capture STDOUT and STDERR

```go
output, err := RunOutput("whoami")
```

### Start

Start an async command. If the executable has suffix ".go" then it will be "go install"ed then executed.
Use this for watching a server task.

```go
Start("main.go", M{"$in": "cmd/app"})
```

Godo tracks the process ID of started processes in the map `Processes` to restart the app gracefully.

If `Start` is not working as expected with go files, verify `go install -a main.go` works in the 
directory of the file, replacing "main.go" with your file. Most issues around (re)starts are 
due to being outside of `GOPATH` or multiple packages in a directory. Set `Env` if `GOPATH`
needs to be adjusted.

### Inside

To run many commands inside a directory, use `Inside` instead of the `$in` option.
`Inside` changes the working directory.

```go
Inside("somedir", func() {
    Run("...")
    Bash("...")
})
```

## User Input

To get plain string

```go
user := Prompt("user: ")
```

To get password

```go
password := PromptPassword("password: ")
```

## Godofile Run-Time Environment

To specify whether to inherit from parent's process environment,
set `InheritParentEnv`. This setting defaults to true

```go
InheritParentEnv = false
```

To specify the base environment for your tasks, set `Env`.
Separate with whitespace or newlines.

```go
Env = `
    GOPATH=.vendor::$GOPATH
    PG_USER=mario
`
```

TIP: Set the `Env` when using a dependency manager like `godep`

```go
wd, _ := os.Getwd()
ws := path.Join(wd, "Godeps/_workspace")
Env = fmt.Sprintf("GOPATH=%s::$GOPATH", ws)
```
Functions can add or override environment variables as part of the command string.
Note that environment variables are set before the executable similar to a shell;
however, the `Run` and `Start` functions do not use a shell.

```go
p.Task("build", func() {
    Run("GOOS=linux GOARCH=amd64 go build" )
})
```

The effective environment for exec functions is: `parent (if inherited) <- Env <- func parsed env`

Paths should use `::` as a cross-platform path list separator. On Windows `::` is replaced with `;`.
On Mac and linux `::` is replaced with `:`.
