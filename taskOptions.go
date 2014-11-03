package godo

// Dependencies are tasks which must run before a task.
type Dependencies []string

// Watch are the glob file patterns which are watched and trigger rerunning
// a task on change.
type Watch []string

// Debounce is the number of milliseconds before a task can run again.
type Debounce int64

// D is short for Dependencies option.
type D Dependencies

// W is short for Watch option.
type W Watch
