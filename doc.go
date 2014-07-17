// Package gosu is a project build toolkit for Go in the spirit of Rake, Grunt and
// others. Gosu supports watching, file globs, tasks and modular projects.
//
// Gosu requires a project configuration function, in which task are
// registered and other projects imported.
//
// To install
//
//      go get -u github.com/mgutz/gosu/gosu
//
// As an example, create a file 'tasks/Gosufile.go'
//
//      package tasks
//
//      import (
//			. "github.com/mgutz/gosu"
//		)
//
//      func Tasks(p *Project) {
//          p.Task("default", Pre{"hello, "views"})
//
//          p.Task("hello", func() {
//              util.Exec(`bash -c "echo Hello $USER!"`)
//          })
//
//          p.Task("views", Watch{"**/*"}, func(c *Context) {
//              if c.FileEvent == nil {
//                  for _, f := range c.Task.WatchFiles {
//                      // f.FileInfo and f.Path
//                      fmt.Printf("File: %s\n", f.Path)
//                  }
//              } else {
//                  // change event when watching
//                  fmt.Printf("%v\n", c.FileEvent)
//              }
//          })
//      }
//
// To run "views"
//
//      gosu hello
//
// To run the "default" task which runs the dependencies "hello", "views"
//
//      gosu
//
// Note the "views" task specifies "**/*" which is a glob pattern
// for watching everything. To rerun "views" whenever any file changes, run gosu in watch mode
//
//      gosu --watch
//
package gosu
