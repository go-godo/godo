# Getting Started

Create a file `Gododir/main.go` with this content

```go
import (
    "fmt"
    do "github.com/mgutz/godo/v2"
)

func tasks(p *do.Project) {
    p.Task("hello", nil, func(c *do.Context) {
        name := c.Args.MayString("world", "name", "n")
        fmt.Println("Hello", name, "!")
    }
}
```

From your terminal run

```
# prints "Hello world!"
godo hello
# prints "Hello gopher!"
godo hello -- n="gopher"
# prints "Hello gopher!"
godo hello -- name="gopher"
```

Let's create a non-trivial example, in which we want to run
tests whenever any go file changes.

```go
func tasks(p *do.Project) {
    p.Task("clean", nil, func(c *do.Context) {
        c.Run("rm -rf tmp")
    }

    p.Task("assets", nil, func(c *do.Context) {
        // Version is from external version.go
        versionDir := "dist/public/v" + Version
        c.Bash(`
            set -e
            mkdir -p {{.versionDir}}
            browserify . -o {{.versionDir}}
        `, M{"versionDir": versionDir})
    }

    p.Task("build", nil, func(c *do.Context) {
        c.Run("go build", M{"$in", "cmd/app"})
    }.Src("cmd/app/**/*.go")

    p.Task("test", nil, func(c *do.Context) {
        c.Run("go test")
    }.Src("**/*.go")


    // S==Series P==Parallel
    p.Task("default", S{"clean", P{"build", "assets"}, "test"}, nil)
}
```

From your terminal run

```sh
godo -w
```

That simple statement does the following

*   godo runs "default" task. Godo will use the "default" task in the absence of a task name.
*   The "default" task defines a set of dependencies qualified by the order in which they should be executed. The dependency

    ```go
    S{P{"clean", "build"}, "test"}
    ```

    means. Run the following in a series.

    1.  Run "clean" and "build" in parallel
    2.  Run "test"






