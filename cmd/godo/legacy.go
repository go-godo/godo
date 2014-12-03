package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mgutz/str"
	"gopkg.in/godo.v1"
	"gopkg.in/godo.v1/util"
)

// This file contains support for the legacy `tasks/Godofile.go` tasks definitions file.
// I (GeertJohan) propose to drop legacy support in godo.v2

func hasMain(data []byte) bool {
	hasMainRe := regexp.MustCompile(`\nfunc main\(`)
	matches := hasMainRe.Find(data)
	return len(matches) > 0
}

func isPackageMain(data []byte) bool {
	isMainRe := regexp.MustCompile(`(\n|^)?package main\b`)
	matches := isMainRe.Find(data)
	return len(matches) > 0
}

func runLegacy() (bool, error) {
	rel := "tasks/Godofile.go"
	filename, err := filepath.Abs(rel)
	if err != nil {
		panic("Could not get absolute path " + filename)
	}
	if !util.FileExists(filename) {
		return false, nil
	}

	mainFile := legacyBuildMain(rel)
	if mainFile != "" {
		filename = mainFile
		defer os.RemoveAll(filepath.Dir(mainFile))
	}
	cmd := "go run " + filename + " " + strings.Join(os.Args[1:], " ")
	// errors are displayed by tasks
	godo.Run(cmd)

	return true, nil
}

func legacyBuildMain(src string) string {
	tempFile := ""
	data, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !hasMain(data) {
		if isPackageMain(data) {
			msg := `%s is not runnable. Rename package OR make it runnable by adding

	func main() {
		godo.Godo(Tasks)
	}
`
			fmt.Printf(msg, src)
			os.Exit(1)
		}

		template := `package main

import (
	"gopkg.in/godo.v1"
	tasks "{{package}}"
)

func main() {
	godo.Godo(tasks.Tasks)
}
`
		packageName, err := util.PackageName(src)
		if err != nil {
			panic(err)
		}
		code := str.Template(template, map[string]interface{}{
			"package": filepath.ToSlash(packageName),
		})
		//log.Println("DBG template", code)
		tempDir, err := ioutil.TempDir("", "godo")
		if err != nil {
			panic("Could not create temp directory")
		}
		//log.Printf("code\n %s\n", code)
		tempFile = filepath.Join(tempDir, "Godofile_main.go")
		err = ioutil.WriteFile(tempFile, []byte(code), 0644)
		if err != nil {
			log.Panicf("Could not write temp file %s\n", tempFile)
		}

		src = tempFile
		return src
	}
	return ""
}
