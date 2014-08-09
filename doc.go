// Package gosu is a project build toolkit for Go in the spirit of Rake, Grunt and
// others. Gosu supports watching, file globs, tasks and modular projects.
//
// Gosu requires a tasks configuration function, in which task are
// registered and other tasks imported.
//
// To install
//
//      go get -u github.com/mgutz/gosu
//      go get -u github.com/mgutz/gosu/cmd/gosu
//
// As an example, create a file 'tasks/Gosufile.go'
//
//      package tasks
//
//      import (
//          . "github.com/mgutz/gosu"
//      )
//
//      func Tasks(p *Project) {
//          p.Task("default", Pre{"hello, "views"})
//
//          p.Task("hello", func() {
//              util.Exec(`bash -c "echo Hello $USER!"`)
//          })
//
//          p.Task("views", Watch{"**/*.go.html"}, func(c *Context) {
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
//      gosu views
//
// To run the "default" task which runs the dependencies "hello", "views"
//
//      gosu
//
// Note the "views" task specifies "**/*.go.html", which is a glob pattern
// to watch any file with .go.html extension. To rerun "views" whenever any file changes, run gosu in watch mode
//
//      gosu --watch
//
package gosu
