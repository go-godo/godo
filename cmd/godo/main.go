package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/mgutz/minimist"
	"github.com/mgutz/str"
	"gopkg.in/godo.v1"
	"gopkg.in/godo.v1/util"
)

var isWindows = runtime.GOOS == "windows"

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

	exe := buildMain(rel)
	cmd := exe + " " + strings.Join(os.Args[1:], " ")
	cmd = str.Clean(cmd)
	// errors are displayed by tasks

	godo.Run(cmd)
}

type godorc struct {
	ModTime time.Time
	Size    int64
}

func mustBeMain(src string) {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !hasMain(data) {
		msg := `%s is not runnable. Rename package OR make it runnable by adding

	func main() {
		godo.Godo(tasks)
	}
	`
		fmt.Printf(msg, src)
		os.Exit(1)
	}

	if !isPackageMain(data) {
		msg := `%s is not runnable. It must be package main`
		fmt.Printf(msg, src)
		os.Exit(1)
	}
}

func buildMain(src string) string {
	mustBeMain(src)

	exeFile := "godobin-" + godo.Version
	if isWindows {
		exeFile += ".exe"
	}

	dir := filepath.Dir(src)
	exe := filepath.Join(dir, exeFile)
	reasonFormat := "Godofile changed. Rebuilding %s...\n"

	argm := minimist.Parse()
	rebuild := argm.ZeroBool("rebuild")
	if rebuild {
		os.Remove(exe)
	}

	// see if last run exists .godoinfo
	fiExe, err := os.Stat(exe)
	build := os.IsNotExist(err)
	if build {
		reasonFormat = "Building %s...\n"
	}

	fiGodofile, err := os.Stat(src)
	if os.IsNotExist(err) {
		log.Fatalln(err)
		os.Exit(1)
	}
	build = build || fiExe.ModTime().Before(fiGodofile.ModTime())

	if build {
		util.Debug("godo", reasonFormat, exe)

		err := godo.Run("go build -a -o "+exeFile, godo.In{dir})
		if err != nil {
			panic(err)
		}
	}

	if rebuild {
		util.Info("godo", "ok")
		os.Exit(0)
	}

	return exe
}
