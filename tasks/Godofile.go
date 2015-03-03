package main

import (
	"fmt"

	"github.com/mgutz/goa"
	f "github.com/mgutz/goa/filter"
	"github.com/mgutz/str"
	. "gopkg.in/godo.v1"
)

func tasks(p *Project) {

	p.Task("dist", D{"test", "lint"})

	p.Task("install", func() {
		Run("go get github.com/golang/lint/golint")
		Run("go get github.com/mgutz/goa")
		Run("go get github.com/robertkrimen/godocdown/godocdown")
	})

	p.Task("lint", func() {
		Run("golint .")
		Run("gofmt -w -s .")
		Run("go vet .")
	})

	p.Task("readme", func() {
		Run("godocdown -o README.md")
		// add godoc
		goa.Pipe(
			f.Load("./README.md"),
			f.Str(str.ReplaceF("--", "\n[godoc](https://godoc.org/gopkg.in/godo.v1)\n", 1)),
			f.Write(),
		)
	})

	p.Task("test", func() {
		Run("go test")
	})

	p.Task("build", func() {
		Run("go install", In{"cmd/godo"})
	})

	p.Task("interactive", func() {
		Bash(`
			echo name?
			read name
			echo hello $name
		`)
	})

	p.Task("whoami", func() {
		Run("whoami")
	})

	pass := 0
	p.Task("err2", func() error {
		if pass == 2 {
			return fmt.Errorf("oh oh")
		}
		return nil
	})

	p.Task("err", D{"err2"}, func() error {
		pass++
		if pass == 1 {
			return nil
		}
		return fmt.Errorf("foo err")
	}).Watch("test/*.txt")

	p.Task("hello", Debounce(3000), W{"*.hello"}, func(c *Context) {
		name := c.Args.MayString("default value", "name", "n")
		fmt.Println("Hello", name)
	})

	p.Task("server", func() {
		Start("main.go", M{"$in": "cmd/example"})
	}).Watch("cmd/example/**/*.go")
}

func main() {
	Godo(tasks)
}
