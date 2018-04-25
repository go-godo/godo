package godo

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
	NOTE: Watching tests

	// touch initial modtime of files
	touch("tmp/foo.txt", 0)

	// start project
	func tasks(p *Project) {
		// ...
	}
	execClI(tasks, ...)

	// give the project time to initialize
	<- time.After(testProjectDelay)

	// use touch to change files. Modtime has to be after the initial
	// modtime, so the watcher can pick up the change. To make one file
	// newer than the other:
	touch("tmp/older.txt", 1*time.Second)
	touch("tmp/newer.txt", 2*time.Second)

	// wait for at least the watch delay, which is the interval
	// the watchers polls the file system
	<- time.After(testWatchDelay)

	// finally, do assertions
*/

// default [templates:*go.html compile]
//
// if any go.html changes then run "templates", "compile", "default"
func TestWatchTasksWithoutSrcShouldAlwaysRun(t *testing.T) {
	trace := ""
	pass := 1
	tasks := func(p *Project) {
		p.Task("A", nil, func(*Context) {
			trace += "A"
		})
		p.Task("B", nil, func(*Context) {
			trace += "B"
		})
		p.Task("C", nil, func(*Context) {
			trace += "C"
		}).Src("test/sub/foo.txt")

		p.Task("default", S{"A", "B", "C"}, func(*Context) {
			// on watch, the task is run once before watching
			if pass == 2 {
				p.Exit(0)
			}
			pass++
		})
	}

	go func() {
		execCLI(tasks, []string{"-w"}, nil)
	}()

	<-time.After(testProjectDelay)

	touch("test/sub/foo.txt", 0)

	<-time.After(testWatchDelay)

	assert.Equal(t, "ABCABC", trace)
}

// default [templates:*go.html styles:*.scss compile ]
//
// if any go.html changes then run "templates", "compile", "default".
// styles is not run since no changes to SCSS files occurred.
func TestWatchWithSrc(t *testing.T) {
	trace := ""
	pass := 1
	tasks := func(p *Project) {
		p.Task("compile", nil, func(*Context) {
			trace += "C"
		})

		p.Task("styles", nil, func(*Context) {
			trace += "S"
		}).Src("test/styles/*scss")

		p.Task("templates", nil, func(*Context) {
			trace += "T"
		}).Src("test/templates/*go.html")

		p.Task("default", S{"templates", "styles", "compile"}, func(*Context) {
			// on watch, the task is run once before watching
			if pass == 2 {
				p.Exit(0)
			}
			pass++
		})
	}

	go func() {
		execCLI(tasks, []string{"-w"}, nil)
	}()

	<-time.After(testProjectDelay)

	touch("test/templates/1.go.html", 100*time.Millisecond)

	<-time.After(testWatchDelay)

	assert.Equal(t, "TSCTC", trace)
}

func TestWatchShouldWatchNamespaceTasks(t *testing.T) {
	done := make(chan bool)
	trace := ""
	pass := 0

	dbTasks := func(p *Project) {
		p.Task("default", S{"models"}, func(*Context) {
			trace += "D"
			pass++
			if pass == 2 {
				p.Exit(0)
			}
		})
		p.Task("models", nil, func(*Context) {
			trace += "M"
		}).Src("test/sub/*.txt")
	}

	tasks := func(p *Project) {
		p.Use("db", dbTasks)
	}

	go func() {
		execCLI(tasks, []string{"db:default", "-w"}, func(code int) {
			done <- true
		})
	}()

	touchTil("test/sub/sub1.txt", 200*time.Millisecond, done)
	assert.Equal(t, "MDMD", trace)

	// test non-watch
	trace = ""
	runTask(tasks, "db:default")
	assert.Equal(t, "MD", trace)
}

func TestWatch(t *testing.T) {
	done := make(chan bool)
	ran := 0
	tasks := func(p *Project) {
		// this should run twice, watch always runs all tasks first then
		// the touch below
		p.Task("txt", nil, func(*Context) {
			ran++
			if ran == 2 {
				p.Exit(0)
			}
		}).Src("test/*.txt")
	}

	status := -1
	go func() {
		argv := []string{"txt", "-w"}
		execCLI(tasks, argv, func(code int) {
			status = code
			done <- true
		})
	}()

	touchTil("test/bar.txt", 100*time.Millisecond, done)
	assert.Equal(t, 2, ran)
}

func TestOutdatedNoDest(t *testing.T) {
	done := make(chan bool)
	ran := ""
	// each task will run once before watch is called
	tasks := func(p *Project) {
		p.Task("txt", nil, func(*Context) {
			ran += "T"
		}).
			Src("test/*.txt").
			Dest("tmp/*.foo2")

		p.Task("parent", S{"txt"}, func(*Context) {
			ran += "P"
			if strings.Count(ran, "P") == 2 {
				p.Exit(0)
			}
		}).Src("test/sub/*.txt")
	}

	status := -1
	go func() {
		argv := []string{"parent", "-w"}
		execCLI(tasks, argv, func(code int) {
			status = code
			done <- true
		})
	}()

	touchTil("test/sub/sub1.txt", 100*time.Millisecond, done)
	assert.Equal(t, "TPP", ran)
}

func TestOutdated(t *testing.T) {
	// force txt to be newer than foo, which should run txt
	touch("tmp/sub/1.foo", -1*time.Second)
	touch("tmp/sub/foo.txt", 0)

	ran := ""
	var project *Project
	tasks := func(p *Project) {
		project = p

		p.Task("txt", nil, func(*Context) {
			ran += "T"
		}).
			Src("tmp/sub/*.txt").
			Dest("tmp/sub/*.foo")

		p.Task("parent", S{"txt"}, func(*Context) {
			ran += "P"
		}).Src("tmp/*.txt")
	}

	go func() {
		argv := []string{"parent", "-w"}
		execCLI(tasks, argv, func(code int) {
			assert.Fail(t, "should not have exited")
		})
	}()

	// give the task enough time to setup watches
	<-time.After(testProjectDelay)

	// force txt to have older modtime than foo which should not run "text but run "parent"
	touch("tmp/1.foo", 3*time.Second)   // is not watched
	touch("tmp/foo.txt", 1*time.Second) // txt is watched

	// wait at least watchDelay which is the interval for updating watches
	<-time.After(testWatchDelay)

	assert.Equal(t, "TPP", ran)
}
