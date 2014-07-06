package gosu

import "flag"
import "os"

// Runs a project.
func Run(projectFunc func(*Project)) {
	project := NewProject()
	projectFunc(project)

	flag.Parse()

	if *help {
		project.Usage()
		os.Exit(0)
	}

	// Run each task including their dependencies.
	args := flag.Args()

	if len(args) == 0 {
		if project.Tasks["default"] != nil {
			project.Run("default")
		} else {
			flag.Usage = project.Usage
		}
	} else {
		for _, name := range flag.Args() {
			project.Run(name)
		}
	}

	if *watching {
		project.Watch(flag.Args())
	}
}
