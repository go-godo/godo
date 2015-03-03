package godo

import (
	"fmt"
	"os"
	"sync"

	"github.com/mgutz/minimist"
	"gopkg.in/godo.v1/util"
)

var watching bool
var help bool
var verbose bool
var version bool
var deprecatedWarnings bool

// DebounceMs is the default time (1500 ms) to debounce task events in watch mode.
var DebounceMs int64
var waitgroup sync.WaitGroup
var waitExit bool
var argm minimist.ArgMap
var contextArgm minimist.ArgMap

func init() {
	DebounceMs = 2000

}

// Usage prints a usage screen with task descriptions.
func Usage(tasks string) {
	// go's flag package prints ugly screen
	format := `godo %s - do task(s)

Usage: godo [flags] [task...]
  -D             Print deprecated warnings
  -h, --help     This screen
      --rebuild  Rebuild Godofile
  -v  --verbose  Log verbosely
  -V, --version  Print version
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
	godo(tasksFunc, nil)
}

func godo(tasksFunc func(*Project), argv []string) {
	if argv == nil {
		argm = minimist.Parse()
	} else {
		argm = minimist.ParseArgv(argv)
	}

	help = argm.ZeroBool("help", "h", "?")
	verbose = argm.ZeroBool("verbose", "v")
	version = argm.ZeroBool("version", "V")
	watching = argm.ZeroBool("watch", "w")
	deprecatedWarnings = argm.ZeroBool("D")
	contextArgm = minimist.ParseArgv(argm.Unparsed())

	project := NewProject(tasksFunc)

	if help {
		Usage(project.usage())
		os.Exit(0)
	}

	if version {
		fmt.Printf("godo %s\n", Version)
		os.Exit(0)
	}

	// Run each task including their dependencies.
	args := []string{}
	for _, v := range argm.Leftover() {
		args = append(args, fmt.Sprintf("%v", v))
	}

	if len(args) == 0 {
		if project.Tasks["default"] != nil {
			args = append(args, "default")
		} else {
			Usage(project.usage())
			os.Exit(0)
		}
	}

	// quick fix to make cascading watch work on default task
	if len(args) == 1 && args[0] == "default" {
		args = project.Tasks["default"].Dependencies
	}

	for _, name := range args {
		err := project.Run(name)
		if err != nil {
			util.Error("ERR", "%s\n", err.Error())
			os.Exit(1)
		}
	}

	if watching {
		if project.Watch(args, true) {
			waitgroup.Add(1)
			waitExit = true
		} else {
			fmt.Println("Nothing to watch. Use W{} or Watch{} to specify glob patterns")
			os.Exit(0)
		}
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
