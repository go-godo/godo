package gosu

import (
	//"log"
	"testing"
)

func TestSimpleTask(t *testing.T) {
	project := NewProject()
	result := ""
	var proj = func(p *Project) {
		p.Task("foo", func(c *Context) {
			result = "A"
		})
	}
	proj(project)
	project.Run("foo")
	if result != "A" {
		t.Error("should have run simple task")
	}
}

func TestDependency(t *testing.T) {
	project := NewProject()
	result := ""
	var proj = func(p *Project) {
		p.Task("foo", func(c *Context) {
			result = "A"
		})

		p.Task("bar", []string{"foo"})
	}
	proj(project)
	project.Run("bar")
	if result != "A" {
		t.Error("should have run task's dependency")
	}
}

func TestMultiProject(t *testing.T) {
	result := ""

	var otherProj = func(p *Project) {
		p.Task("foo", []string{"bar"}, func(c *Context) {
			result += "B"
		})

		p.Task("bar", func(c *Context) {
			result += "C"
		})
	}

	project := NewProject()
	var proj = func(p *Project) {
		p.Use("other", otherProj)

		p.Task("foo", func(c *Context) {
			result += "A"
		})

		p.Task("bar", []string{"foo", "other:foo"})
	}
	proj(project)

	project.Run("bar")
	if result != "ACB" {
		t.Error("should have run dependent project")
	}
}

func TestShouldExpandGlobs(t *testing.T) {
	project := NewProject()
	result := ""
	var proj = func(p *Project) {
		p.Task("foo", Files{"test/**/*.txt"}, func(c *Context) {
			result = "A"
		})

		p.Task("bar", Files{"test/**/*.html"}, []string{"foo"})
	}
	proj(project)
	project.Run("bar")
	if len(project.Tasks["bar"].WatchFiles) != 1 {
		t.Error("bar should have 1 HTML file")
	}
	if len(project.Tasks["foo"].WatchFiles) != 5 {
		t.Error("foo should have 5 txt files, one is hidden")
	}
}
