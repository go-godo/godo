package gosu

import (
	"sort"
	"testing"
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
