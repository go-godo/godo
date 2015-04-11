package main

import (
	"fmt"

	// "github.com/mgutz/goa"
	// f "github.com/mgutz/goa/filter"

	. "gopkg.in/godo.v2"
)

func tasks(p *Project) {
	p.Task("test", nil, func(*Context) {
		Run("go test")
	})

	p.Task("test", S{"build"}, func(c *Context) {
		c.Run("go test")
	})

	p.Task("dist", S{"test", "lint"}, nil)

	p.Task("install", nil, func(c *Context) {
		c.Run("go get github.com/golang/lint/golint")
		// Run("go get github.com/mgutz/goa")
		c.Run("go get github.com/robertkrimen/godocdown/godocdown")
	})

	p.Task("lint", nil, func(c *Context) {
		c.Run("golint .")
		c.Run("gofmt -w -s .")
		c.Run("go vet .")
	})

	// p.Task("readme", func() {
	// 	Run("godocdown -o README.md")
	// 	// add godoc
	// 	goa.Pipe(
	// 		f.Load("./README.md"),
	// 		f.Str(str.ReplaceF("--", "\n[godoc](https://godoc.org/gopkg.in/godo.v2)\n", 1)),
	// 		f.Write(),
	// 	)
	// })

	p.Task("build", nil, func(c *Context) {
		c.Run("go install", M{"$in": "cmd/godo"})
	})

	p.Task("interactive", nil, func(c *Context) {
		c.Bash(`
			echo name?
			read name
			echo hello $name
		`)
	})

	p.Task("whoami", nil, func(c *Context) {
		Run("whoami")
	})

	pass := 0
	p.Task("err2", nil, func(*Context) {
		if pass == 2 {
			Halt("oh oh")
		}
	})

	p.Task("err", S{"err2"}, func(*Context) {
		pass++
		if pass == 1 {
			return
		}
		Halt("foo err")
	}).Src("test/*.txt")

	p.Task("hello", nil, func(c *Context) {
		name := c.Args.MayString("default value", "name", "n")
		fmt.Println("Hello", name)
	}).Src("*.hello").Debounce(3000)

	p.Task("server", nil, func(*Context) {
		Start("main.go", M{"$in": "cmd/example"})
	}).Src("cmd/example/**/*.go")

	p.Task("change-package", nil, func(*Context) {
		// works on mac
		Run(`find . -name "*.go" -print | xargs sed -i "" 's|gopkg.in/godo.v1|gopkg.in/godo.v2|g'`)
		// maybe linux?
		//Run(`find . -name "*.go" -print | xargs sed -i 's|gopkg.in/godo.v1|gopkg.in/godo.v2|g'`)
	})

}

func main() {
	Godo(tasks)
}
