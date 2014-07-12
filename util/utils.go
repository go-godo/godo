package util

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mgutz/str"
)

// Exec is sugary way to execute a command. Exec converts cmd into an
// an executable and an argv.
func Exec(command string) {
	argv := str.ToArgv(command)
	executable := argv[0]
	argv = argv[1:]
	cmd := exec.Command(executable, argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		Error("ERR", "running: %s\n", command)
	}
}

// FileExists determines if path exists
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// PackageName determines the package name from sourceFile if it is within $GOPATH
func PackageName(sourceFile string) (string, error) {
	if filepath.Ext(sourceFile) != ".go" {
		return "", errors.New("sourcefile must end with .go")
	}
	sourceFile, err := filepath.Abs(sourceFile)
	if err != nil {
		panic("Could not convert to absolute path: " + sourceFile)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", errors.New("Environment variable GOPATH is not set")
	}
	paths := strings.Split(gopath, string(os.PathListSeparator))
	for _, path := range paths {
		srcDir := filepath.Join(path, "src")
		srcDir, err := filepath.Abs(srcDir)
		if err != nil {
			continue
		}

		//log.Printf("srcDir %s sourceFile %s\n", srcDir, sourceFile)
		rel, err := filepath.Rel(srcDir, sourceFile)
		if err != nil {
			continue
		}
		return filepath.Dir(rel), nil
	}
	return "", errors.New("sourceFile not reachable from GOPATH")
}
