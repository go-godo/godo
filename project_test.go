package gosu

import (
	//"log"
	"testing"
)

func TestSimpleTask(t *testing.T) {
	project := NewProject()
	result := ""
	tasks := func(task TaskFunc) {
		task("foo", func(c *Context) {
			result = "A"
		})
	}
	project.Define(tasks)
	project.Run("foo")
	if result != "A" {
		t.Error("should have run simple task")
	}
}

func TestDependency(t *testing.T) {
	project := NewProject()
	result := ""
	tasks := func(task TaskFunc) {
		task("foo", func(c *Context) {
			result = "A"
		})

		task("bar", Pre{"foo"})
	}
	project.Define(tasks)
	project.Run("bar")
	if result != "A" {
		t.Error("should have run task's dependency")
	}
}

func TestMultiProject(t *testing.T) {
	result := ""

	otherTasks := func(task TaskFunc) {
		task("foo", Pre{"bar"}, func(c *Context) {
			result += "B"
		})

		task("bar", func(c *Context) {
			result += "C"
		})
	}

	project := NewProject()
	tasks := func(task TaskFunc, use UseFunc) {
		use("other", otherTasks)

		task("foo", func(c *Context) {
			result += "A"
		})

		task("bar", Pre{"foo", "other:foo"})
	}
	project.Define(tasks)

	project.Run("bar")
	if result != "ACB" {
		t.Error("should have run dependent project")
	}
}

func TestShouldExpandGlobs(t *testing.T) {
	project := NewProject()
	result := ""
	tasks := func(task TaskFunc) {
		task("foo", Watch{"test/**/*.txt"}, func(c *Context) {
			result = "A"
		})

		task("bar", Watch{"test/**/*.html"}, Pre{"foo"})
	}
	project.Define(tasks)
	project.Run("bar")
	if len(project.Tasks["bar"].WatchFiles) != 1 {
		t.Error("bar should have 1 HTML file")
	}
	if len(project.Tasks["foo"].WatchFiles) != 5 {
		t.Error("foo should have 5 txt files, one is hidden")
	}
}
