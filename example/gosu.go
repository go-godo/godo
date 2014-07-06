package main

import (
	"fmt"
	"os/exec"

	"github.com/mgutz/gosu"
)

// ImportedProject could be an imported project from someone else's library
func ImportedProject(p *gosu.Project) {
	p.Task("sprite", func(c *gosu.Context) {
		fmt.Printf("creating sprite image ...\n")
	})
}

// Define your project's task inside a function
func Project(p *gosu.Project) {
	// User other projects in namespace
	p.Use("ext", ImportedProject)

	p.Task("default", "Default task", []string{"styles", "views", "ext:sprite"})

	p.Task("styles", gosu.Files{"public/css/*.less"}, func(c *gosu.Context) {
		if c.FileEvent != nil {
			// inspect watcher file events
		}
		exec.Command("lessc", "public/css/styles.less", "public/css/styles.css").Run()
	})

	p.Task("views", gosu.Files{"views/**/*.go.html"}, func() {
		exec.Command("gorazor", "views", "views").Run()
	})

	p.Task("restart", "(Re)starts the app", gosu.Files{"**/*.go"}, func() {
		// (re)start your app
	})
}

func main() {
	gosu.Run(Project)
}
