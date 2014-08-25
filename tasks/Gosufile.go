package tasks

import (
	"github.com/mgutz/goa"
	f "github.com/mgutz/goa/filter"
	. "github.com/mgutz/gosu"
	"github.com/mgutz/str"
)

// Tasks is local project.
func Tasks(p *Project) {

	p.Task("dist", D{"lint", "readme"})

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
			f.Str(str.ReplaceF("--", "\n[godoc](https://godoc.org/github.com/mgutz/gosu)\n", 1)),
			f.Str(str.BetweenF("", "## Usage")),
			f.Write(),
		)
	})
}
