// Package gosu is a project build toolkit for Go in the spirit of Rake, Grunt and
// others. Gosu supports watching, file globs, tasks and modular projects.
//
// Gosu does not provide an executable. Instead, use gosu to build your own
// project build tool. Gosu requires a project configuration function,
// in which task are registered and other projects imported.
//
// Here is an example
//
//      import "github.com/mgutz/gosu"
//
//      func Project(p *gosu.Project) {
//          p.Task("default", []string{"views", "styles"})
//
//          p.Task("views", gosu.Files{"views/**/*.go.html"}, func() {
//              exec.Command("razor", "views", "views").Run()
//          })
//
//          p.Task("styles", gosu.Files{"css/**/*.less"}, func() {
//              exec.Command("lessc", "css/styles.less", "css/styles.css").Run()
//          })
//      }
//
//      func main() {
//          gosu.Run(Project)
//      }
//
// To run "views"
//      go build -o gosu gosu.go
//      ./gosu views
// Or
//      go run gosu.go views
//
// To run default, which in turn runs "views", "styles"
//
//      ./gosu
//
// Note the "views" task specifies "views/**/*.go.html" which is a glob pattern
// used when watching. To rerun "views" whenever any file ending with "go.html"
// changes, run gosu in watch mode
//
//      ./gosu --watch
package gosu
