package gosu

import (
	"os"

	flag "github.com/ogier/pflag"
)

var watching = flag.Bool("watch", false, "Watch task and dependencies")
var help = flag.Bool("help", false, "View this usage screen")
var verbose = flag.Bool("verbose", false, "View more info like which file changed")
var debounceMilliseconds = flag.Int64("debounce", 2000, "Default number of milliseconds to debounce tasks")

// Run runs a project of tasks.
func Run(tasksFunc func(*Project)) {
	project := NewProject(tasksFunc)

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
		project.Watch(flag.Args(), true)
	}
}

// MustNotError checks if error is not nil. If it is not nil it will panic.
func MustNotError(err error) {
	if err != nil {
		panic(err)
	}
}
