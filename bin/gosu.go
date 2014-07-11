package main

import (
	"github.com/mgutz/goa"
	f "github.com/mgutz/goa/filter"
	"github.com/mgutz/gosu"
	"github.com/mgutz/gosu/util"
	"github.com/mgutz/str"
)

func Project(p *gosu.Project) {

	p.Task("dist", []string{"lint", "readme"})

	p.Task("lint", func() {
		util.Exec("golint .")
		util.Exec("gofmt -w -s .")
		util.Exec("go vet .")
	})

	p.Task("readme", func() {
		util.Exec("godocdown -o README.md")
		// add godoc
		goa.Pipe(
			f.Load("./README.md"),
			f.Str(str.ReplaceF("--", "\n[godoc](https://godoc.org/github.com/mgutz/gosu)\n", 1)),
			f.Write(),
		)
	})
}

func main() {
	gosu.Run(Project)
}
