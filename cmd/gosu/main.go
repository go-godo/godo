package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mgutz/gosu/util"
	"github.com/mgutz/str"
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
	gosuFiles := []string{"Gosufile.go", "tasks/Gosufile.go"}
	src := ""
	rel := ""
	for _, filename := range gosuFiles {
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

	mainFile := buildMain(rel)
	if mainFile != "" {
		src = mainFile
		defer os.RemoveAll(filepath.Dir(mainFile))
	}
	cmd := "go run " + src + " " + strings.Join(os.Args[1:], " ")
	//log.Printf("DBG %s\n", cmd)
	// errors are displayed by tasks
	util.ExecError(cmd)
}

func buildMain(src string) string {
	tempFile := ""
	data, err := ioutil.ReadFile(src)
	checkError(err, "%s not found\n", src)

	if !hasMain(data) {
		if isPackageMain(data) {
			msg := `%s is not runnable. Rename package OR make it runnable by adding

    func main() {
        gosu.Run(Tasks)
    }
`
			fmt.Printf(msg, src)
			os.Exit(1)
		}

		template := `
	        package main
	        import (
	            "github.com/mgutz/gosu"
	            pkg "{{package}}"
	        )
	        func main() {
	            gosu.Run(pkg.Tasks)
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
		tempDir, err := ioutil.TempDir("", "gosu")
		if err != nil {
			panic("Could not create temp directory")
		}
		//log.Printf("code\n %s\n", code)
		tempFile = filepath.Join(tempDir, "Gosufile_main.go")
		err = ioutil.WriteFile(tempFile, []byte(code), 0644)
		if err != nil {
			log.Panicf("Could not write temp file %s\n", tempFile)
		}

		src = tempFile
		return src
	}
	return ""
}
