package godo

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/mgutz/str"
)

var isWindows = runtime.GOOS == "windows"

func TestEnvironment(t *testing.T) {
	var user string
	if isWindows {
		user = os.Getenv("USERNAME")
		os.Setenv("USER", user)
	} else {
		user = os.Getenv("USER")
	}

	SetEnviron("USER=$USER:godo", true)
	env := effectiveEnv(nil)
	if !sliceContains(env, "USER="+user+":godo") {
		t.Error("Environment interpolation failed", env)
	}

	SetEnviron("USER=$USER:godo", false)
	env = effectiveEnv(nil)
	if len(env) != 1 {
		t.Error("Disabling parent inheritance failed")
	}
	if !sliceContains(env, "USER="+user+":godo") {
		t.Error("Should have read parent var even if not inheriting")
	}

	// set back to defaults
	SetEnviron("", true)
	l := len(os.Environ())
	env = effectiveEnv([]string{"USER=$USER:$USER:func"})
	if !sliceContains(env, "USER="+user+":"+user+":func") {
		t.Error("Should have been overriden by func environmnt")
	}
	if len(env) != l {
		t.Error("Effective environment length changed")
	}

	env = effectiveEnv([]string{"GOSU_NEW_VAR=foo"})
	if !sliceContains(env, "GOSU_NEW_VAR=foo") {
		t.Error("Should have new var")
	}
	if len(env) != l+1 {
		t.Error("Effective environment length should have increased by 1")
	}

	SetEnviron(`
		USER1=$USER
		USER2=$USER1
	`, true)
	env = effectiveEnv([]string{"USER3=$USER2"})
	if !sliceContains(env, "USER1="+user) {
		t.Error("Should have interpolated from parent env")
	}
	if !sliceContains(env, "USER3="+user) {
		t.Error("Should have interpolated from effective env")
	}

	env = effectiveEnv([]string{"PATH=foo::bar::bah"})
	if !sliceContains(env, "PATH=foo"+string(os.PathListSeparator)+"bar"+string(os.PathListSeparator)+"bah") {
		t.Error("Should have replaced PathSeparator, got", env)
	}

	// set back to defaults
	SetEnviron("", true)
}

func TestQuotedVar(t *testing.T) {
	// set back to defaults
	defer SetEnviron("", true)
	env := effectiveEnv([]string{`FOO="a=bar b=bah c=baz"`})
	v := getEnv(env, "FOO", false)
	if v != `"a=bar b=bah c=baz"` {
		t.Errorf("Quoted var failed %q", v)
	}
}

func TestExpansion(t *testing.T) {
	SetEnviron(`
		FOO=foo
		FAIL=$FOObar:godo
		OK=${FOO}bar:godo
	`, true)

	env := effectiveEnv([]string{})
	if !sliceContains(env, "FAIL=:godo") {
		t.Error("$FOObar should not have interpolated")
	}
	if !sliceContains(env, "OK=foobar:godo") {
		t.Error("${FOO}bar should have expanded", env)
	}
}

func TestInheritedRunEnv(t *testing.T) {
	os.Setenv("TEST_RUN_ENV", "fubar")
	SetEnviron("", true)

	var output string

	if isWindows {	
		output, _ = RunOutput(`FOO=bar BAH=baz cmd /C "echo %TEST_RUN_ENV% %FOO%"`)	
	} else {
		output, _ = RunOutput(`FOO=bar BAH=baz bash -c "echo -n $TEST_RUN_ENV $FOO"`)	
	}

	
	if str.Clean(output) != "fubar bar" {
		t.Error("Environment was not inherited! Got", fmt.Sprintf("%q", output))
	}
}
