package tasks

import (
	"fmt"

	"github.com/mgutz/goa"
	f "github.com/mgutz/goa/filter"
	. "github.com/mgutz/gosu"
	"github.com/mgutz/gosu/util"
	"github.com/mgutz/str"
)

// Project is local project.
func Tasks(task TaskFunc) {

	task("hello", func() {
		util.Exec(`bash -c "echo Hello $USER!"`)
	})

	task("hello2", func() {
		fmt.Println(Hello("foobar"))
	})

	task("files", Files{"**/*"}, func(c *Context) {
		if c.FileEvent == nil {
			for _, f := range c.Task.WatchFiles {
				// f.FileInfo and f.Path
				fmt.Printf("File: %s\n", f.Path)
			}
		} else {
			// change event when watching
			fmt.Printf("%v\n", c.FileEvent)
		}
	})

	task("dist", Pre{"lint", "readme"})

	task("lint", func() {
		util.Exec("golint .")
		util.Exec("gofmt -w -s .")
		util.Exec("go vet .")
	})

	task("readme", func() {
		util.Exec("godocdown -o README.md")
		// add godoc
		goa.Pipe(
			f.Load("./README.md"),
			f.Str(str.ReplaceF("--", "\n[godoc](https://godoc.org/github.com/mgutz/gosu)\n", 1)),
			f.Write(),
		)
	})
}
