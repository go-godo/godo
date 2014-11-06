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

func checkError(err error, format string, args ...interface{}) {
	if err != nil {
		util.Error("ERR", format, args...)
		os.Exit(1)
	}
}

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

func main() {
	// legacy version used tasks/
	godoFiles := []string{"Gododir/Godofile.go", "tasks/Godofile.go"}
	src := ""
	rel := ""
	for _, filename := range godoFiles {
		rel = filename
		filename, err := filepath.Abs(filename)
		if err != nil {
			panic("Could not get absolute path " + filename)
		}
		if !util.FileExists(filename) {
			continue
		}
		src = filename
		break
	}

	if src == "" {
		godo.Usage("")
		fmt.Printf("\n\n%s not found\n", src)
		os.Exit(1)
	}

	mainFile := buildMain(rel)
	if mainFile != "" {
		src = mainFile
		defer os.RemoveAll(filepath.Dir(mainFile))
	}
	cmd := "go run " + src + " " + strings.Join(os.Args[1:], " ")
	// errors are displayed by tasks
	godo.Run(cmd)
}

func buildMain(src string) string {
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

		template := `
	        package main
	        import (
	            "gopkg.in/godo.v1"
	            pkg "{{package}}"
	        )
	        func main() {
	            godo.Godo(pkg.Tasks)
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
