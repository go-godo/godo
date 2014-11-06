package godo

import (
	"os"
	"regexp"
	"strings"

	"github.com/mgutz/str"
)

// Env is the default environment to use for all commands. That is,
// the effective environment for all commands is the merged set
// of (parent environment, Env, func specified environment). Whitespace
// or newline separate key value pairs. $VAR interpolation is allowed.
//
// Env = "GOOS=linux GOARCH=amd64"
// Env = `
//   GOOS=linux
//   GOPATH=./vendor:$GOPATH
// `
var Env string

// PathListSeparator is a cross-platform path list separator. On Windows, PathListSeparator
// is replacd by ";". On others, PathListSeparator is replaced by ":"
var PathListSeparator = "::"

// InheritParentEnv whether to inherit parent's environment
var InheritParentEnv bool

func init() {
	InheritParentEnv = true
}

var envvarRe = regexp.MustCompile(`\$(\w+)`)

func interpolateEnv(env []string, kv string) string {
	if strings.Contains(kv, PathListSeparator) {
		kv = strings.Replace(kv, PathListSeparator, string(os.PathListSeparator), -1)
	}

	// find all key=$EXISTING_VAR:foo and interpolate from os.Environ()
	matches := envvarRe.FindAllStringSubmatch(kv, -1)
	for _, match := range matches {
		existingVar := match[1]
		kv = strings.Replace(kv, "$"+existingVar, getEnv(env, existingVar, true), -1)
	}
	return kv
}

func getEnv(env []string, key string, checkParent bool) string {
	for _, kv := range env {
		pair := strings.Split(kv, "=")
		if pair[0] == key {
			return pair[1]
		}
	}

	if checkParent {
		return os.Getenv(key)
	}
	return ""
}

// upsertenv updates or inserts a key=value pair into an environment.
func upsertenv(env *[]string, kv string) {
	pair := strings.Split(kv, "=")
	if len(pair) != 2 {
		return
	}

	set := false
	for i, item := range *env {
		ipair := strings.Split(item, "=")
		if ipair[0] == pair[0] {
			(*env)[i] = interpolateEnv(*env, kv)
			set = true
			break
		}

	}

	if !set {
		*env = append(*env, interpolateEnv(*env, kv))
	}
}

// effectiveEnv is the effective environment for an exec function.
func effectiveEnv(funcEnv []string) []string {
	var env []string
	if InheritParentEnv {
		env = os.Environ()
	} else {
		env = []string{}
	}

	// merge in package Env
	for _, kv := range parseStringEnv(Env) {
		upsertenv(&env, kv)
	}

	// merge in func's env
	if funcEnv != nil {
		for _, kv := range funcEnv {
			upsertenv(&env, kv)
		}
	}
	return env
}

// parseStringEnv parse the package Env string and converts it into an
// environment slice.
func parseStringEnv(s string) []string {
	env := []string{}

	if s == "" {
		return env
	}

	s = str.Clean(s)
	pairs := strings.Split(s, " ")
	for _, kv := range pairs {
		if !strings.Contains(kv, "=") {
			continue
		}
		env = append(env, kv)
	}
	return env
}
