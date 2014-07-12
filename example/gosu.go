package main

import (
	"fmt"

	"github.com/mgutz/gosu"
	"github.com/mgutz/gosu/util"
)

// ImportedProject could be an imported project from someone else's library
func ImportedProject(p *gosu.Project) {
	p.Task("sprite", func(c *gosu.Context) {
		fmt.Printf("creating sprite image\n")
	})
}

// Project is your local project. Define your tasks here.
func Project(p *gosu.Project) {
	// User other projects in namespace
	p.Use("ext", ImportedProject)

	p.Task("default", "Default task", []string{"views", "ext:sprite"})

	p.Task("views", gosu.Files{"views/**/*.go.html"}, func(c *gosu.Context) {
		util.Exec("razor views views")
	})

	p.Task("restart", "(Re)starts the app", gosu.Files{"**/*.go"}, func() {
		fmt.Printf("Restarting app")
		// (re)start your app
	})
}
