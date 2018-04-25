package godo

import (
	"sort"
	"testing"
	"time"

	"github.com/mgutz/str"
	"github.com/stretchr/testify/assert"
)

func TestSimpleTask(t *testing.T) {
	result := ""
	tasks := func(p *Project) {
		p.Task1("foo", func(c *Context) {
			result = "A"
		})
	}

	runTask(tasks, "foo")
	if result != "A" {
		t.Error("should have run simple task")
	}
}

func TestErrorReturn(t *testing.T) {
	result := ""
	tasks := func(p *Project) {
		p.Task1("err", func(*Context) {
			Halt("error caught")
			// should not get here
			result += "ERR"
		})

		p.Task("foo", S{"err"}, func(*Context) {
			result = "A"
		})
	}

	_, err := runTask(tasks, "foo")
	if result == "A" {
		t.Error("parent task should not run on error")
	}
	if err.Error() != `"foo>err": error caught` {
		t.Error("dependency errors should stop parent")
	}

	_, err = runTask(tasks, "err")
	if err.Error() != `"err": error caught` {
		t.Error("error was not handle properly")
	}
}

func TestTaskArgs(t *testing.T) {
	assert := assert.New(t)
	result := ""
	tasks := func(p *Project) {
		p.Task1("foo", func(c *Context) {
			name := c.Args.MustString("name")
			result = name
		})
	}

	execCLI(tasks, []string{"foo", "--", "--name=gopher"}, nil)
	assert.Equal("gopher", result)
	assert.Panics(func() {
		runTask(tasks, "foo")
	})
}

func TestDependency(t *testing.T) {
	result := ""
	tasks := func(p *Project) {
		p.Task1("foo", func(c *Context) {
			result = "A"
		})

		p.Task("bar", S{"foo"}, nil)
	}
	runTask(tasks, "bar")
	if result != "A" {
		t.Error("should have run task's dependency")
	}
}

func TestShouldExpandGlobs(t *testing.T) {
	result := ""
	tasks := func(p *Project) {
		p.Task("foo", nil, func(c *Context) {
			result = "A"
		}).Src("test/**/*.txt")

		p.Task("bar", S{"foo"}, nil).Src("test/**/*.html")
	}
	proj, err := runTask(tasks, "bar")
	assert.NoError(t, err)
	if len(proj.Tasks["bar"].SrcFiles) != 2 {
		t.Error("bar should have 2 HTML file")
	}
	if len(proj.Tasks["foo"].SrcFiles) != 7 {
		t.Error("foo should have 7 txt files, one is hidden, got",
			len(proj.Tasks["foo"].SrcFiles))
	}
}

func TestCalculateWatchPaths(t *testing.T) {
	// test wildcards, should watch current directory
	paths := []string{
		"example/views/**/*.go.html",
		"example.html",
	}
	paths = calculateWatchPaths(paths)
	if len(paths) != 1 {
		t.Error("Expected exact elements")
	}
	sort.Strings(paths)
	if paths[0] != "." {
		t.Error("Expected exact file paths got", paths[0])
	}

	// should only watch current directory
	paths = []string{
		"**/*.go.html",
		"example.html",
	}
	paths = calculateWatchPaths(paths)

	if len(paths) != 1 {
		t.Error("Expected exact elements")
	}
	if paths[0] != "." {
		t.Error("Expected . got", paths[0])
	}
}

func TestLegacyIn(t *testing.T) {

	var cat = "cat"
	if isWindows {
		cat = "cmd /c type"
	}
	//// Run

	// in V2 BashOutput accepts an options map

	out, err := RunOutput(cat+" foo.txt", M{"$in": "test"})
	assert.NoError(t, err)
	assert.Equal(t, "foo", str.Clean(out))

	if isWindows {
		return
	}

	//// Bash

	// in V2 BashOutput accepts an options map
	out, err = BashOutput("cat foo.txt", M{"$in": "test"})
	assert.NoError(t, err)
	assert.Equal(t, "foo", str.Clean(out))
}

func TestInvalidTask(t *testing.T) {
	tasks := func(p *Project) {
	}

	assert.Panics(t, func() {
		runTask(tasks, "dummy")
	})
}

func TestParallel(t *testing.T) {
	var result string
	tasks := func(p *Project) {
		p.Task1("A", func(*Context) {
			result += "A"
		})
		p.Task1("B", func(*Context) {
			time.Sleep(10 * time.Millisecond)
			result += "B"
		})
		p.Task("C", nil, func(*Context) {
			result += "C"
		})
		p.Task("D", nil, func(*Context) {
			result += "D"
		})
		p.Task("default", P{"A", "B", "C", "D"}, nil)
	}

	argv := []string{}
	ch := make(chan int)
	go func() {
		execCLI(tasks, argv, func(code int) {
			assert.Equal(t, code, 0)
			assert.True(t, len(result) == 4)
			ch <- code
			close(ch)
		})
	}()
	<-ch
}

func TestTaskD(t *testing.T) {
	trace := ""
	tasks := func(p *Project) {
		p.Task1("A", func(*Context) {
			trace += "A"
		})

		p.TaskD("default", S{"A"})
	}
	runTask(tasks, "default")
	assert.Equal(t, "A", trace)
}

func TestRunOnce(t *testing.T) {
	trace := ""
	tasks := func(p *Project) {
		p.Task("once?", nil, func(*Context) {
			trace += "1"
		})

		p.Task("A", S{"once"}, func(*Context) {
			trace += "A"
		})

		p.Task("B", S{"once"}, func(*Context) {
			trace += "B"
		})
	}

	execCLI(tasks, []string{"A", "B"}, nil)
	assert.Equal(t, "1AB", trace)
}
