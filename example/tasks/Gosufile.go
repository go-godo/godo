package tasks

import (
	"fmt"
	. "github.com/mgutz/gosu"
	"github.com/mgutz/gosu/util"
)

// ImportedTasks could be an imported project from someone else's library
func ImportedTasks(p *Project) {
	p.Task("sprite", func(c *Context) {
		fmt.Printf("creating sprite image\n")
	})
}

// Tasks is your local project. Define your tasks here.
func Tasks(p *Project) {
	// User other projects in namespace
	p.Use("ext", ImportedTasks)

	//p.Task("default", "Default task", []string{"views", "ext:sprite"})
	p.Task("default", Pre{"views", "ext:sprite"})

	p.Task("views", Watch{"views/**/*.go.html"}, func(c *Context) {
		util.Exec("razor views views")
	})

	p.Task("restart", "(Re)starts the app", Watch{"**/*.go"}, func() {
		fmt.Printf("Restarting app")
		// (re)start your app
	})

	p.Task("files", Watch{"**/*"}, func(c *Context) {
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
}
