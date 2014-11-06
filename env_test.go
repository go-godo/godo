package godo

import (
	"os"
	"testing"
)

func TestEnvironment(t *testing.T) {
	Env = `
	USER=$USER:godo
	`
	user := os.Getenv("USER")
	env := effectiveEnv(nil)
	if !sliceContains(env, "USER="+user+":godo") {
		t.Error("Environment interpolation failed")
	}

	InheritParentEnv = false
	env = effectiveEnv(nil)
	if len(env) != 1 {
		t.Error("Disabling parent inheritance failed")
	}
	if !sliceContains(env, "USER="+user+":godo") {
		t.Error("Should have read parent var even if not inheriting")
	}

	// set back to defaults
	Env = ""
	InheritParentEnv = true

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

	Env = `
	USER1=$USER
	USER2=$USER1
	`
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
	Env = ""
	InheritParentEnv = true
}

func TestLongFormInterpolation(t *testing.T) {
	Env = `
	FOO=foo
	FAIL=$FOObar:godo
	OK=${FOO}bar:godo
	`

	env := effectiveEnv([]string{})
	if !sliceContains(env, "FAIL=:godo") {
		t.Error("$FOObar should not have interpolated")
	}
	if !sliceContains(env, "OK=foobar:godo") {
		t.Error("${FOO}bar should have expanded", env)
	}
}
