# Tasks

A task is a named unit of work which can be executed in series or parallel. Define smaller tasks to form larger tasks.

```go

var S = do.S, M = do.M, Context = do.Context

func tasks(p *do.Project) {
    p.Task("assets", nil, func(c *Context) {
        c.Run("webpack")
    })

    project.Task("build", S{"build"}, func(c *Context) {
        c.Run("go run main.go", M{"$in": "cmd/app"})
    })
}
```
