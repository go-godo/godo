package tasks

import (
	"github.com/mgutz/goa"
	f "github.com/mgutz/goa/filter"
	. "github.com/mgutz/gosu"
	"github.com/mgutz/gosu/util"
	"github.com/mgutz/str"
)

// Tasks is local project.
func Tasks(p *Project) {

	p.Task("dist", Pre{"lint"})

	p.Task("install", func() {
		util.Run("go get github.com/golang/lint/golint")
		util.Run("go get github.com/mgutz/goa")
		util.Run("go get github.com/robertkrimen/godocdown/godocdown")

	})

	p.Task("lint", func() {
		util.Run("golint .")
		util.Run("gofmt -w -s .")
		util.Run("go vet .")
	})

	p.Task("readme", func() {
		util.Run("godocdown -o README.md")
		// add godoc
		goa.Pipe(
			f.Load("./README.md"),
			f.Str(str.ReplaceF("--", "\n[godoc](https://godoc.org/github.com/mgutz/gosu)\n", 1)),
			f.Str(str.BetweenF("", "## Usage")),
			f.Write(),
		)
	})
}
