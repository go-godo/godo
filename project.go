package gosu

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mgutz/gosu/fsnotify"
)

var watching = flag.Bool("watch", false, "Watch task and dependencies")
var help = flag.Bool("help", false, "View this usage screen")

// Project is a container for tasks.
type Project struct {
	Tasks     map[string]*Task
	Namespace map[string]*Project
}

// NewProject creates am empty project ready for tasks.
func NewProject() *Project {
	project := &Project{Tasks: map[string]*Task{}}
	project.Namespace = map[string]*Project{}
	project.Namespace[""] = project
	return project
}

func (project *Project) mustTask(name string) (*Project, *Task) {
	namespace, taskName := project.namespaceTaskName(name)

	proj := project.Namespace[namespace]
	if proj == nil {
		Panicf("project", "Could not find project having namespace \"%s\"\n", namespace)
	}

	task := proj.Tasks[taskName]
	if task == nil {
		Errorf("project", "Task is not defined \"%s\"\n", name)
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
func (project *Project) RunWithEvent(name string, logName string, e *fsnotify.FileEvent) {
	project.run(name, logName, e)
}

// run runs the project, executing any tasks named on the command line.
func (project *Project) run(name string, logName string, e *fsnotify.FileEvent) error {
	_, task := project.mustTask(name)

	// Run each task including their dependencies.
	for _, depName := range task.Dependencies {
		namespace, taskName := project.namespaceTaskName(depName)
		proj := project.Namespace[namespace]
		if proj == nil {
			fmt.Errorf("Project was not loaded for \"%s\" task", name)
		}
		project.Namespace[namespace].RunWithEvent(taskName, name+">"+depName, e)
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
		fmt.Printf("  %-24s %s\n", cyan(name), description)
	}
}

// Use uses another project's task within a namespace.
func (project *Project) Use(namespace string, projectFunc func(*Project)) {
	namespace = strings.Trim(namespace, ":")
	proj := NewProject()
	projectFunc(proj)
	project.Namespace[namespace] = proj
}

// Task adds a task to the project.
func (project *Project) Task(name string, args ...interface{}) *Task {
	task := &Task{Name: name}

	for _, t := range args {
		switch t := t.(type) {
		default:
			Panicf("project", "unexpected type %T", t) // %T prints whatever type t has
		case Files:
			task.WatchGlobs = t
		case []string:
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

func shortestDir(files []*FileAsset) string {
	dirs := []string{}
	for _, fa := range files {
		dirs = append(dirs, fa.Path)
	}
	sort.Strings(dirs)
	return path.Dir(dirs[0])
}

func watchTask(root string, taskName string, handler func(e *fsnotify.FileEvent)) {
	bufferSize := 2048
	watcher, err := fsnotify.NewWatcher(bufferSize)
	if err != nil {
		Panicf("project", "%v\n", err)
	}
	waitTime := time.Duration(0.2 * float64(time.Second))
	watcher.WatchRecursive(root)
	watcher.ErrorHandler = func(err error) {
		Errorf("project", "%v\n", err)
	}

	// this function will block forever, Ctrl+C to quit app
	lastHappendTime := time.Now()
	firstTime := true
	lastRename := ""
	for {
		if firstTime {
			Infof(taskName, "watching %s ...\n", magenta(root))
			firstTime = false
		}
		event := <-watcher.Event
		//changing a file sends rename and create as two separate events
		if event.IsRename() {
			lastRename = event.Name
		}
		if event.IsCreate() {
			if lastRename == event.Name {
				continue
			}
			lastRename = ""
		}
		if event.Time.Before(lastHappendTime) {
			continue
		}
		handler(event)
		// prevent multiple restart in short time
		time.Sleep(waitTime)
		lastHappendTime = time.Now()
	}
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

	taskClosure := func(project *Project, task *Task, taskname string) func() {
		root := shortestDir(task.WatchFiles)
		return func() {
			watchTask(root, taskname, func(e *fsnotify.FileEvent) {
				project.run(taskname, taskname, e)
			})
		}
	}

	for _, taskname := range names {
		proj, task := project.mustTask(taskname)

		if len(task.WatchFiles) > 0 {
			funcs = append(funcs, taskClosure(proj, task, taskname))
		}

		// TODO should this be recursive? --mario
		if len(task.Dependencies) > 0 {
			for _, depname := range task.Dependencies {
				proj, task := project.mustTask(depname)
				if len(task.WatchFiles) > 0 {
					funcs = append(funcs, taskClosure(proj, task, taskname))
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
