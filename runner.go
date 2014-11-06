package godo

import (
	"fmt"
	"os"
	"sync"

	flag "github.com/ogier/pflag"
)

var watching = flag.BoolP("watch", "w", false, "")
var help = flag.BoolP("help", "h", false, "")
var verbose = flag.Bool("verbose", false, "")
var version = flag.BoolP("version", "v", false, "")

// DebounceMs is the default time (1500 ms) to debounce task events in watch mode.
var DebounceMs int64
var waitgroup sync.WaitGroup
var waitExit bool

func init() {
	DebounceMs = 2000
}

// Usage prints a usage screen with task descriptions.
func Usage(tasks string) {
	// go's flag package prints ugly screen
	format := `godo %s - do task(s)

Usage: godo [flags] [task...]
  -h, --help     This screen
      --verbose  Log verbosely
  -v, --version  Print version
  -w, --watch    Watch task and dependencies`

	if tasks == "" {
		fmt.Printf(format, Version)
	} else {
		format += "\n\n%s"
		fmt.Printf(format, Version, tasks)
	}
}

// Godo runs a project of tasks.
func Godo(tasksFunc func(*Project)) {
	flag.Parse()

	project := NewProject(tasksFunc)

	if *help {
		Usage(project.usage())
		os.Exit(0)
	}

	if *version {
		fmt.Printf("godo %s", Version)
	}

	// Run each task including their dependencies.
	args := flag.Args()
	if len(args) == 0 {
		if project.Tasks["default"] != nil {
			project.Run("default")
		} else {
			Usage(project.usage())
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
