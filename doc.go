// Package godo is a task runner, file watcher in the spirit of Rake, Gulp ...
//
// To install
//
//      go get -u gopkg.in/godo.v1/cmd/godo
//
// As an example, create a file 'tasks/Godofile.go' with this content
//
//    package main
//
//    import (
//        . "gokpkg.in/godo.v1"
//    )
//
//    func Tasks(p *Project) {
//        Env = "GOPATH=.vendor:$GOPATH OTHER_VAR=val"
//
//        p.Task("default", D{"hello", "build"})
//
//        p.Task("hello", func() {
//            Bash("echo Hello $USER!")
//        })
//
//        p.Task("build", W{"**/*.go"}, func() {
//            Run("GOOS=linux GOARCH=amd64 go build", In{"cmd/godo"})
//        })
//
//        p.Task("server", D{"views"}, W{"**/*.go"}, Debounce(3000), func() {
//            // Start recompiles and restarts on changes when watching
//            Start("main.go", In{"cmd/server"})
//        })
//    }
//
//    func main() {
//        Godo(Tasks)
//    }
//
// To run "views" task from terminal
//
//      godo views
//
// To rerun "views" whenever any `*.go.html` file changes
//
//      godo views --watch
//
// To run the "default" task which runs "hello" and "views"
//
//      godo
//
// Task names may add a "?" suffix to indicate run only once
//
//      // run once regardless of number of dependees
//      p.Task("build?", func() {})
//
// Task options
//
//      D{} or Dependencies{} - dependent tasks which run before task
//      Debounce              - minimum milliseconds before task can run again
//      W{} or Watches{}      - array of glob file patterns to watch
//
//          /**/   - match zero or more directories
//          {a,b}  - match a or b, no spaces
//          *      - match any non-separator char
//          ?      - match a single non-separator char
//          **/    - match any directory, start of pattern only
//          /**    - match any this directory, end of pattern only
//          !      - removes files from resultset, start of pattern only
//
// Task handlers
//
//      func() {}           - Simple function handler
//      func(c *Context) {} - Handler which accepts the current context
package godo
