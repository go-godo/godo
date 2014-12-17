package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/godo.v1"
	"gopkg.in/godo.v1/util"
)

const (
	filenameGododir = `Gododir`
	tmplMainText    = `package main

import (
	"gopkg.in/godo.v1"
	tasks "{{.packageImportPath}}"
)

func main() {
	godo.Godo(tasks.Tasks)
}
`
)

func checkError(err error, format string, args ...interface{}) {
	if err != nil {
		util.Error("ERR", format, args...)
		os.Exit(1)
	}
}

func main() {
	// check if `Gododir` exists
	if util.DirExists(filenameGododir) {
		// get absolute path
		gododirAbs, err := filepath.Abs(filenameGododir)
		if err != nil {
			fmt.Printf("error getting absolute path for Gododir: %v\n", err)
			os.Exit(1)
		}
		// use go/build to parse the files in the directory and get package information
		pkg, err := build.ImportDir(gododirAbs, 0)
		if err != nil {
			fmt.Printf("error in Gododir: %v\n", err)
			return
		}
		// error when package is a command
		if pkg.IsCommand() {
			fmt.Println("error: the source files in Gododir should not be a command (package main).")
			os.Exit(1)
		}

		// create temp go file that imports the Gododir package and executes the tasks
		tmplMain := template.Must(template.New("main").Parse(tmplMainText))
		tempDir, err := ioutil.TempDir("", "godo")
		if err != nil {
			fmt.Printf("could not create temp dir: %v\n", err)
			os.Exit(1)
		}
		tempFileName := filepath.Join(tempDir, "Gododir_main.go")
		tempFile, err := os.Create(tempFileName)
		if err != nil {
			fmt.Printf("error creating temp main file: %v\n", err)
			os.Exit(1)
		}
		defer tempFile.Close()
		err = tmplMain.Execute(tempFile, map[string]string{
			"packageImportPath": pkg.ImportPath,
		})
		if err != nil {
			fmt.Printf("could not write/execute main template to temp file: %v\n", err)
			os.Exit(1)
		}
		tempFile.Close()

		// run package
		cmd := "go run " + tempFileName + " " + strings.Join(os.Args[1:], " ")
		godo.Run(cmd)

		// all done
		return
	}

	// fallback to legacy method
	ok, err := runLegacy()
	if err != nil {
		fmt.Printf("error running tasks/Godofile.go: %v\n", err)
		os.Exit(1)
	}
	if !ok {
		godo.Usage("")
		fmt.Println("")
		os.Exit(1)
	}
}
