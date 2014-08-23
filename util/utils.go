package util

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/mgutz/str"
)

// ExecError is a simple way to execute a CLI utility.
func RunError(command string, options ...map[string]interface{}) error {
	argv := str.ToArgv(command)
	executable := argv[0]
	argv = argv[1:]
	// for _, arg := range args {
	// 	argv = append(argv, arg)
	// }
	cmd := exec.Command(executable, argv...)

	if len(options) == 1 {
		opts := options[0]
		if opts["Dir"] != nil {
			cmd.Dir = opts["Dir"].(string)
		}
		if opts["Env"] != nil {
			cmd.Env = opts["Env"].([]string)
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Exec is simple way to execute a CLI utility. `command` is parsed
// for arguments. args is optional and unparsed.
func Run(command string, options ...map[string]interface{}) {
	err := RunError(command, options...)
	if err != nil {
		Error("ERR", "%s\n%+v", command, err)
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

// Template reads a go template and writes it to dist given data.
func Template(src string, dest string, data map[string]interface{}) {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		Panic("template", "Could not read file %s\n", src)
	}

	tpl := template.New("vagrantFile")
	tpl, err = tpl.Parse(string(content))
	if err != nil {
		Panic("template", "Could not parse template %s\n", src)
	}

	f, err := os.Create(dest)
	if err != nil {
		Panic("template", "Could not create file for writing %s\n", dest)
	}
	defer f.Close()
	err = tpl.Execute(f, data)
	if err != nil {
		Panic("template", "Could not execute template %s\n", src)
	}
}

func StartError(command string, options ...map[string]interface{}) error {
	argv := str.ToArgv(command)
	executable := argv[0]
	argv = argv[1:]
	// for _, arg := range args {
	// 	argv = append(argv, arg)
	// }
	cmd := exec.Command(executable, argv...)
	if len(options) == 1 {
		opts := options[0]
		if opts["Dir"] != nil {
			cmd.Dir = opts["Dir"].(string)
		}
		if opts["Env"] != nil {
			cmd.Env = opts["Env"].([]string)
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	h := sha1.New()
	io.WriteString(h, command)
	pidfile := filepath.Join(os.TempDir(), fmt.Sprintf("%x.pid", h.Sum(nil)))

	killExisting := func() {
		if _, err := os.Stat(pidfile); !os.IsNotExist(err) {
			pidb, err := ioutil.ReadFile(pidfile)
			if err != nil {
				Error("Start", "Could not read pidfile %s\n", pidfile)
				return
			}
			pid := ToInt(string(pidb))
			if pid != 0 {
				existingProcess, err := os.FindProcess(pid)
				if err != nil {
					Error("Start", "Could not find process %d\n", pid)
					return
				}
				if existingProcess != nil {
					err := existingProcess.Kill()
					if err != nil {
						Error("Start", "Could not kill existing process %+v\n", existingProcess)
						return
					}
					_, err = existingProcess.Wait()
					if err != nil {
						Error("Start", "Error waiting %v\n", err)
						return
					}
				}
			}
		}
	}

	killExisting()

	err := cmd.Start()
	if err != nil {
		Error("Start", "Could not start process %s\n", command)
		return err
	}
	return ioutil.WriteFile(pidfile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
}

// Start is simple way to execute a CLI utility. `command` is parsed
// for arguments. args is optional and unparsed.
func Start(command string, options ...map[string]interface{}) {
	err := StartError(command, options...)
	if err != nil {
		Error("ERR", "%s\n%+v", command, err)
	}
}

func ToInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return result
}
