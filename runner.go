package gosu

import (
	"os"

	flag "github.com/ogier/pflag"
)

var watching = flag.Bool("watch", false, "Watch task and dependencies")
var help = flag.Bool("help", false, "View this usage screen")
var verbose = flag.Bool("verbose", false, "View more info like which file changed")

// Run runs a project as defined within projectFunc handler function.
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
