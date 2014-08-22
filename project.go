package gosu

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mgutz/gosu/util"
	"github.com/mgutz/gosu/watcher"
	"github.com/mgutz/str"
)

// M is generic string to interface alias
type M map[string]interface{}

// Project is a container for tasks.
type Project struct {
	Tasks     map[string]*Task
	Namespace map[string]*Project
}

// NewProject creates am empty project ready for tasks.
func NewProject(tasksFunc func(*Project)) *Project {
	project := &Project{Tasks: map[string]*Task{}}
	project.Namespace = map[string]*Project{}
	project.Namespace[""] = project
	project.Define(tasksFunc)
	return project
}

func (project *Project) mustTask(name string) (*Project, *Task) {
	namespace, taskName := project.namespaceTaskName(name)

	proj := project.Namespace[namespace]
	if proj == nil {
		util.Panic("project", "Could not find project having namespace \"%s\"\n", namespace)
	}

	task := proj.Tasks[taskName]
	if task == nil {
		util.Error("ERR", `"%s" task is not defined`+"\n", name)
		os.Exit(1)
	}
	return proj, task
}

func (project *Project) namespaceTaskName(name string) (namespace string, taskName string) {
	namespace = ""
	taskName = name
	if strings.Contains(name, ":") {
		parts := strings.Split(name, ":")
		namespace = parts[0]
		taskName = parts[1]
	}
	return
}

// Run runs a task by name.
func (project *Project) Run(name string) {
	project.run(name, name, nil)
}

// RunWithEvent runs a task by name and adds FileEvent e to the context.
func (project *Project) runWithEvent(name string, logName string, e *watcher.FileEvent) {
	project.run(name, logName, e)
}

// run runs the project, executing any tasks named on the command line.
func (project *Project) run(name string, logName string, e *watcher.FileEvent) error {
	_, task := project.mustTask(name)

	// Run each task including their dependencies.
	for _, depName := range task.Dependencies {
		namespace, taskName := project.namespaceTaskName(depName)
		proj := project.Namespace[namespace]
		if proj == nil {
			fmt.Errorf("Project was not loaded for \"%s\" task", name)
		}
		project.Namespace[namespace].runWithEvent(taskName, name+">"+depName, e)
	}
	task.RunWithEvent(logName, e)
	return nil
}

// Usage prints usage about the app and tasks.
func (project *Project) Usage() {
	fmt.Printf("Usage: %s [flags] [task...]\n\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Printf("\nTasks\n\n")

	names := []string{}
	m := map[string]*Task{}
	for ns, proj := range project.Namespace {
		if ns != "" {
			ns += ":"
		}
		for _, task := range proj.Tasks {
			names = append(names, ns+task.Name)
			m[ns+task.Name] = task
		}
	}
	sort.Strings(names)
	longest := 0
	for _, name := range names {
		l := len(name)
		if l > longest {
			longest = l
		}
	}

	for _, name := range names {
		task := m[name]
		description := task.Description
		if description == "" {
			if len(task.Dependencies) > 0 {
				description = "Runs {" + strings.Join(task.Dependencies, ", ") + ", " + name + "} tasks"
			} else {
				description = "Runs " + name + " task"
			}
		}
		fmt.Printf("  %-"+strconv.Itoa(longest)+"s  %s\n", name, description)
	}
}

// Use uses another project's task within a namespace.
func (project *Project) Use(namespace string, tasksFunc func(*Project)) {
	namespace = strings.Trim(namespace, ":")
	proj := NewProject(tasksFunc)
	project.Namespace[namespace] = proj
}

// Task adds a task to the project.
func (project *Project) Task(name string, args ...interface{}) *Task {
	task := &Task{Name: name}

	for _, t := range args {
		switch t := t.(type) {
		default:
			util.Panic("project", "unexpected type %T", t) // %T prints whatever type t has
		case Watch:
			task.WatchGlobs = t
		case Pre:
			task.Dependencies = t
		case func():
			task.Handler = t
		case func(*Context):
			task.ContextHandler = t
		case string:
			task.Description = t
		}
	}
	project.Tasks[name] = task
	return task
}

func watchTask(root string, logName string, handler func(e *watcher.FileEvent)) {
	bufferSize := 2048
	watchr, err := watcher.NewWatcher(bufferSize)
	if err != nil {
		util.Panic("project", "%v\n", err)
	}
	waitTime := time.Duration(0.1 * float64(time.Second))
	watchr.WatchRecursive(root)
	watchr.ErrorHandler = func(err error) {
		util.Error("project", "%v\n", err)
	}

	// this function will block forever, Ctrl+C to quit app
	var lastHappenedTime int64
	firstTime := true
	for {
		if firstTime {
			util.Info(logName, "watching %s ...\n", root)
			firstTime = false
		}
		event := <-watchr.Event
		//util.Debug("DBG", "event %+v\n", event)
		isOlder := event.UnixNano < lastHappenedTime
		lastHappenedTime = event.UnixNano

		if isOlder {
			continue
		}
		handler(event)
		// prevent multiple restart in short time
		time.Sleep(waitTime)
	}
}

// Define defines tasks
func (project *Project) Define(fn func(*Project)) {
	fn(project)
}

func calculateWatchPaths(patterns []string) []string {
	paths := map[string]bool{}
	for _, glob := range patterns {
		if glob == "" {
			continue
		}
		pth := glob
		if strings.Contains(glob, "*") {
			pth = str.Between(glob, "", "*")
			if pth == "" {
				// this means watch current directy, no need to watch anything else
				return []string{"."}
			}
		}
		paths[pth] = true
	}

	var keys []string
	for key := range paths {
		keys = append(keys, key)
	}
	return keys
}

// Watch watches the Files of a task and reruns the task on a watch event. Any
// direct dependency is also watched.
func (project *Project) Watch(names []string) {
	funcs := []func(){}

	if len(names) == 0 {
		if project.Tasks["default"] != nil {
			names = append(names, "default")
		}
	}

	taskClosure := func(project *Project, task *Task, taskname string, logName string) func() {
		paths := calculateWatchPaths(task.WatchGlobs)
		return func() {
			if len(paths) == 0 {
				return
			}
			for _, pth := range paths {
				watchTask(pth, logName, func(e *watcher.FileEvent) {
					project.run(taskname, taskname, e)
				})
			}
		}
	}

	for _, taskname := range names {
		proj, task := project.mustTask(taskname)

		if len(task.WatchFiles) > 0 {
			funcs = append(funcs, taskClosure(proj, task, taskname, taskname))
		}

		// TODO should this be recursive? --mario
		if len(task.Dependencies) > 0 {
			for _, depname := range task.Dependencies {
				proj, task := project.mustTask(depname)
				if len(task.WatchFiles) > 0 {
					funcs = append(funcs, taskClosure(proj, task, taskname, taskname+">"+depname))
				}
			}
		}
	}
	if len(funcs) > 0 {
		done := all(funcs)
		<-done
	}
}

// all runs the functions in fns concurrently.
func all(fns []func()) (done <-chan bool) {
	var wg sync.WaitGroup
	wg.Add(len(fns))

	ch := make(chan bool, 1)
	for _, fn := range fns {
		go func(f func()) {
			f()
			wg.Done()
		}(fn)
	}
	go func() {
		wg.Wait()
		doneSig(ch, true)
	}()
	return ch
}

func doneSig(ch chan bool, val bool) {
	ch <- val
	close(ch)
}
