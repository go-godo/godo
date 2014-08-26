package gosu

import (
	"os"
	//"time"
	"sync"

	//"github.com/mgutz/gosu/util"
	flag "github.com/ogier/pflag"
)

var watching = flag.Bool("watch", false, "Watch task and dependencies")
var help = flag.Bool("help", false, "View this usage screen")
var verbose = flag.Bool("verbose", false, "View more info like which file changed")

// DebounceMs is the default time (1500 ms) to debounce task events in watch mode.
var DebounceMs int64
var waitgroup sync.WaitGroup
var waitExit bool

func init() {
	DebounceMs = 3000
}

// Gosu runs a project of tasks.
func Gosu(tasksFunc func(*Project)) {
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

	if waitExit {
		waitgroup.Wait()
	}
}

// MustNotError checks if error is not nil. If it is not nil it will panic.
func mustNotError(err error) {
	if err != nil {
		panic(err)
	}
}
