package godo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiProject(t *testing.T) {
	result := ""

	otherTasks := func(p *Project) {
		p.Task("foo", S{"bar"}, func(c *Context) {
			result += "B"
		})

		p.Task("bar", nil, func(c *Context) {
			result += "C"
		})
	}

	tasks := func(p *Project) {
		p.Use("other", otherTasks)

		p.Task("foo", nil, func(c *Context) {
			result += "A"
		})

		p.Task("bar", S{"foo", "other:foo"}, nil)
	}
	runTask(tasks, "bar")
	if result != "ACB" {
		t.Error("should have run dependent project")
	}
}

func TestNestedNamespaces(t *testing.T) {
	levels := ""
	var subsubTasks = func(p *Project) {
		p.Task("A", S{"B"}, func(*Context) {
			levels += "2:"
		})
		p.Task("B", nil, func(*Context) {
			levels += "2B:"
		})
	}
	var subTasks = func(p *Project) {
		p.Use("sub", subsubTasks)
		p.Task("A", S{"sub:A"}, func(*Context) {
			levels += "1:"
		})
	}
	var tasks = func(p *Project) {
		p.Use("sub", subTasks)
		p.Task("A", S{"sub:A"}, func(*Context) {
			levels += "0:"
		})
	}

	runTask(tasks, "A")
	assert.Equal(t, levels, "2B:2:1:0:")
}

func TestNestedNamespaceDependency(t *testing.T) {
	levels := ""
	var subsubTasks = func(p *Project) {
		p.Task("A", S{"B"}, func(*Context) {
			levels += "2:"
		})
		p.Task1("B", func(*Context) {
			levels += "2B:"
		})
	}
	var subTasks = func(p *Project) {
		p.Use("sub", subsubTasks)
		p.Task("A", S{"sub:A"}, func(*Context) {
			levels += "1:"
		})
	}
	var tasks = func(p *Project) {
		p.Use("sub", subTasks)
		p.Task("A", S{"sub:sub:A"}, func(*Context) {
			levels += "0:"
		})
	}

	runTask(tasks, "A")
	assert.Equal(t, levels, "2B:2:0:")
}

func TestRelativeNamespace(t *testing.T) {
	levels := ""
	var subsubTasks = func(p *Project) {
		p.Task("A", S{"/sub:A"}, func(*Context) {
			levels += "2:"
		})
	}
	var subTasks = func(p *Project) {
		p.Use("sub", subsubTasks)
		p.Task("A", S{"sub:A"}, func(*Context) {
			levels += "1:"
		})
	}
	var tasks = func(p *Project) {
		p.Use("sub", subTasks)
		p.Task("A", S{"sub:sub:A"}, func(*Context) {
			levels += "0:"
		})
	}

	runTask(tasks, "A")
	assert.Equal(t, "1:2:0:", levels)
}
