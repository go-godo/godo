package godo

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"

	"github.com/mgutz/str"
	"github.com/stretchr/testify/assert"
)

var cat = "cat"

func init() {
	if runtime.GOOS == "windows" {
		cat = "type"
	}
}

func TestRunMultiline(t *testing.T) {
	output, _ := RunOutput(`
		{{.cat}} test/foo.txt
		{{.cat}} test/bar.txt
		`, M{"cat": cat})
	assert.Equal(t, "foo\nbar\n", output)
}

func TestRunError(t *testing.T) {
	output, err := RunOutput(`
		{{.cat}} test/doesnotexist.txt
		{{.cat}} test/bar.txt
		`, M{"cat": cat})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "line=1")
	assert.Contains(t, output, "doesnotexist")
}

func TestInside(t *testing.T) {
	Inside("test", func() {
		var out string
		if isWindows {
			out, _ = RunOutput("foo.cmd")
		} else {
			out, _ = RunOutput("bash foo.sh")
		}

		if str.Clean(out) != "FOOBAR" {
			t.Error("Inside failed. Got", fmt.Sprintf("%q", out))
		}
	})

	version, _ := ioutil.ReadFile("./VERSION.go")
	if !strings.Contains(string(version), "var Version") {
		t.Error("Inside failed to reset work directory")
	}
}

func TestBash(t *testing.T) {
	if isWindows {
		return
	}
	out, _ := BashOutput(`echo -n foobar`)
	if out != "foobar" {
		t.Error("Simple bash failed. Got", out)
	}

	out, _ = BashOutput(`
		echo -n foobar
		echo -n bahbaz
	`)
	if out != "foobarbahbaz" {
		t.Error("Multiline bash failed. Got", out)
	}

	out, _ = BashOutput(`
		echo -n \
		foobar
	`)
	if out != "foobar" {
		t.Error("Bash line continuation failed. Got", out)
	}

	out, _ = BashOutput(`
		echo -n "foobar"
	`)
	if out != "foobar" {
		t.Error("Bash quotes failed. Got", out)
	}

	out, _ = BashOutput(`
		echo -n "fo\"obar"
	`)
	if out != "fo\"obar" {
		t.Error("Bash quoted string failed. Got", out)
	}
}

func TestTemplatedCommands(t *testing.T) {
	echo := "echo"
	if isWindows {
		echo = "cmd /c echo"

	}
	// in V2 BashOutput accepts an options map
	out, err := RunOutput(echo+" {{.name}}", M{"name": "oy"})
	assert.NoError(t, err)
	assert.Equal(t, "oy", str.Clean(out))

	if isWindows {
		return
	}

	// in V2 BashOutput accepts an options map
	out, err = BashOutput("echo {{.name}}", M{"name": "oy"})
	assert.NoError(t, err)
	assert.Equal(t, "oy", str.Clean(out))
}
