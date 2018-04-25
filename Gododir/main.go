package main

import (
	"fmt"

	// "github.com/mgutz/goa"
	// f "github.com/mgutz/goa/filter"

	do "github.com/davars/godo"
)

func tasks(p *do.Project) {
	p.Task("test", nil, func(c *do.Context) {
		c.Run("go test")
	})

	p.Task("test", do.S{"build"}, func(c *do.Context) {
		c.Run("go test")
	})

	p.Task("dist", do.S{"test", "lint"}, nil)

	p.Task("install", nil, func(c *do.Context) {
		c.Run("go get github.com/golang/lint/golint")
		// Run("go get github.com/mgutz/goa")
		c.Run("go get github.com/robertkrimen/godocdown/godocdown")
	})

	p.Task("lint", nil, func(c *do.Context) {
		c.Run("golint .")
		c.Run("gofmt -w -s .")
		c.Run("go vet .")
	})

	// p.Task("readme", func() {
	// 	Run("godocdown -o README.md")
	// 	// add godoc
	// 	goa.Pipe(
	// 		f.Load("./README.md"),
	// 		f.Str(str.ReplaceF("--", "\n[godoc](https://godoc.org/github.com/davars/godo)\n", 1)),
	// 		f.Write(),
	// 	)
	// })

	p.Task("build", nil, func(c *do.Context) {
		c.Run("go install", do.M{"$in": "cmd/godo"})
	})

	p.Task("interactive", nil, func(c *do.Context) {
		c.Bash(`
			echo name?
			read name
			echo hello $name
		`)
	})

	p.Task("whoami", nil, func(c *do.Context) {
		c.Run("whoami")
	})

	pass := 0
	p.Task("err2", nil, func(*do.Context) {
		if pass == 2 {
			do.Halt("oh oh")
		}
	})

	p.Task("err", do.S{"err2"}, func(*do.Context) {
		pass++
		if pass == 1 {
			return
		}
		do.Halt("foo err")
	}).Src("test/*.txt")

	p.Task("hello", nil, func(c *do.Context) {
		name := c.Args.AsString("default value", "name", "n")
		fmt.Println("Hello", name)
	}).Src("*.hello").Debounce(3000)

	p.Task("server", nil, func(c *do.Context) {
		c.Start("main.go", do.M{"$in": "cmd/example"})
	}).Src("cmd/example/**/*.go")

	p.Task("change-package", nil, func(c *do.Context) {
		// works on mac
		c.Run(`find . -name "*.go" -print | xargs sed -i "" 's|gopkg.in/godo.v1|github.com/davars/godo|g'`)
		// maybe linux?
		//Run(`find . -name "*.go" -print | xargs sed -i 's|gopkg.in/godo.v1|github.com/davars/godo|g'`)
	})
}

func main() {
	do.Godo(tasks)
}
