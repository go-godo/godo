package tasks

import (
	"fmt"

	. "github.com/mgutz/gosu"
	"github.com/mgutz/gosu/util"
)

// ImportedProject could be an imported project from someone else's library
func ImportedTasks(task TaskFunc) {
	task("sprite", func(c *Context) {
		fmt.Printf("creating sprite image\n")
	})
}

// Project is your local project. Define your tasks here.
func Tasks(task TaskFunc, use UseFunc) {
	// User other projects in namespace
	use("ext", ImportedTasks)

	//p.Task("default", "Default task", []string{"views", "ext:sprite"})
	task("default", "Default task", Pre{"views", "ext:sprite"})

	task("views", Watch{"views/**/*.go.html"}, func(c *Context) {
		util.Exec("razor views views")
	})

	task("restart", "(Re)starts the app", Watch{"**/*.go"}, func() {
		fmt.Printf("Restarting app")
		// (re)start your app
	})
}
