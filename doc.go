package gosu

// gosu is a build toolkit for Go in the spirit of Rake and othersl. gosu
// supports watching, file globs, tasks and importing other projects.
//
// gosu basic building block is a task. A task usually has a handler associated
// with it to perform some action such as compiling templates. A task
// also has dependencies. Dependencies are run before the task itself. A
// task must have either a handler or depencies to be valid.
