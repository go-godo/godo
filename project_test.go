package godo

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"github.com/mgutz/str"
	"github.com/stretchr/testify/assert"
)

func TestSimpleTask(t *testing.T) {
	result := ""
	tasks := func(p *Project) {
		p.Task("foo", func(c *Context) {
			result = "A"
		})
	}

	project := NewProject(tasks)
	project.Run("foo")
	if result != "A" {
		t.Error("should have run simple task")
	}
}

func TestErrorReturn(t *testing.T) {
	result := ""
	tasks := func(p *Project) {
		p.Task("err", func() error {
			return fmt.Errorf("error caught")
		})

		p.Task("foo", D{"err"}, func() {
			result = "A"
		})
	}

	project := NewProject(tasks)
	err := project.Run("foo")
	if result == "A" {
		t.Error("parent task should not run on error")
	}
	if err.Error() != `"foo>err": error caught` {
		t.Error("dependency errors should stop parent")
	}

	project.reset()
	err = project.Run("err")
	if err.Error() != `"err": error caught` {
		t.Error("error was not handle properly")
	}
}

func TestTaskArgs(t *testing.T) {
	assert := assert.New(t)
	result := ""
	tasks := func(p *Project) {
		p.Task("foo", func(c *Context) {
			name := c.Args.MustString("name")
			result = name
		})
	}

	godo(tasks, []string{"foo", "--", "--name=gopher"})
	assert.Equal("gopher", result)

	assert.Panics(func() {
		godo(tasks, []string{"foo"})
	})
}

func TestDependency(t *testing.T) {
	result := ""
	tasks := func(p *Project) {
		p.Task("foo", func(c *Context) {
			result = "A"
		})

		p.Task("bar", D{"foo"})
	}
	project := NewProject(tasks)
	project.Run("bar")
	if result != "A" {
		t.Error("should have run task's dependency")
	}
}

func TestMultiProject(t *testing.T) {
	result := ""

	otherTasks := func(p *Project) {
		p.Task("foo", D{"bar"}, func(c *Context) {
			result += "B"
		})

		p.Task("bar", func(c *Context) {
			result += "C"
		})
	}

	tasks := func(p *Project) {
		p.Use("other", otherTasks)

		p.Task("foo", func(c *Context) {
			result += "A"
		})

		p.Task("bar", D{"foo", "other:foo"})
	}
	project := NewProject(tasks)
	project.Run("bar")
	if result != "ACB" {
		t.Error("should have run dependent project")
	}
}

func TestShouldExpandGlobs(t *testing.T) {
	result := ""
	tasks := func(p *Project) {
		p.Task("foo", Watch{"test/**/*.txt"}, func(c *Context) {
			result = "A"
		})

		p.Task("bar", Watch{"test/**/*.html"}, D{"foo"})
	}
	project := NewProject(tasks)
	project.Run("bar")
	if len(project.Tasks["bar"].WatchFiles) != 1 {
		t.Error("bar should have 1 HTML file")
	}
	if len(project.Tasks["foo"].WatchFiles) != 5 {
		t.Error("foo should have 5 txt files, one is hidden")
	}
}

func TestCalculateWatchPaths(t *testing.T) {
	// test wildcards
	paths := []string{
		"example/views/**/*.go.html",
		"example.html",
	}
	paths = calculateWatchPaths(paths)
	if len(paths) != 2 {
		t.Error("Expected exact elements")
	}
	sort.Strings(paths)
	if paths[0] != "example.html" {
		t.Error("Expected exact file paths got", paths[0])
	}
	if paths[1] != "example/views/" {
		t.Error("Expected example/views/ got", paths[1])
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

func TestInside(t *testing.T) {
	Inside("test", func() {
		var out string
		if isWindows {
			out, _ = RunOutput("foo.cmd")
		} else {
			out, _ = RunOutput("bash foo.sh")
		}

		if str.Clean(out) != "FOOBAR" {
			t.Error("Inside failed. Got", fmt.Sprintf("%q", out))
		}
	})

	version, _ := ioutil.ReadFile("./VERSION.go")
	if !strings.Contains(string(version), "var Version") {
		t.Error("Inside failed to reset work directory")
	}
}

func TestBash(t *testing.T) {
	if isWindows {
		return
	}
	out, _ := BashOutput(`echo -n foobar`)
	if out != "foobar" {
		t.Error("Simple bash failed. Got", out)
	}

	out, _ = BashOutput(`
		echo -n foobar
		echo -n bahbaz
	`)
	if out != "foobarbahbaz" {
		t.Error("Multiline bash failed. Got", out)
	}

	out, _ = BashOutput(`
		echo -n \
		foobar
	`)
	if out != "foobar" {
		t.Error("Bash line continuation failed. Got", out)
	}

	out, _ = BashOutput(`
		echo -n "foobar"
	`)
	if out != "foobar" {
		t.Error("Bash quotes failed. Got", out)
	}

	out, _ = BashOutput(`
		echo -n "fo\"obar"
	`)
	if out != "fo\"obar" {
		t.Error("Bash quoted string failed. Got", out)
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
	assert.Equal(t, "asdf", str.Clean(out))

	// need to support V1 though
	out, err = RunOutput(cat+" foo.txt", In{"test"})
	assert.NoError(t, err)
	assert.Equal(t, "asdf", str.Clean(out))

	if isWindows {
		return
	}

	//// Bash

	// in V2 BashOutput accepts an options map
	out, err = BashOutput("cat foo.txt", M{"$in": "test"})
	assert.NoError(t, err)
	assert.Equal(t, "asdf", str.Clean(out))

	// need to support V1 though
	out, err = BashOutput("cat foo.txt", In{"test"})
	assert.NoError(t, err)
	assert.Equal(t, "asdf", str.Clean(out))
}

func TestTemplatedCommands(t *testing.T) {
	echo := "echo"
	if isWindows {
		echo = "cmd /c echo"

	}
	// in V2 BashOutput accepts an options map
	out, err := RunOutput(echo+" {{.name}}", M{"name": "oy"})
	assert.NoError(t, err)
	assert.Equal(t, "oy", str.Clean(out))

	if isWindows {
		return
	}

	// in V2 BashOutput accepts an options map
	out, err = BashOutput("echo {{.name}}", M{"name": "oy"})
	assert.NoError(t, err)
	assert.Equal(t, "oy", str.Clean(out))
}

func sliceContains(slice []string, val string) bool {
	for _, it := range slice {
		if it == val {
			return true
		}
	}
	return false
}
