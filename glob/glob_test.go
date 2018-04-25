package glob

import (
	"regexp"
	"strings"
	"testing"
)

func TestMatching(t *testing.T) {
	var re *regexp.Regexp

	re = Globexp("a")
	if !re.MatchString("a") {
		t.Error("should match exactly")
	}

	re = Globexp("src/**/*.html")
	if !re.MatchString("src/test.html") {
		t.Error("/**/ should match zero directories")
	}

	re = Globexp("src/**/*.html")
	if !re.MatchString("src/foo/bar/test.html") {
		t.Error("/**/ should match intermediate directories")
	}

	re = Globexp("**/*.html")
	if !re.MatchString("test.html") {
		t.Error("**/ should match zero leading directories")
	}

	re = Globexp("**/*.html")
	if !re.MatchString("src/foo/bar/test.html") {
		t.Error("**/ should match leading directories")
	}

	re = Globexp("src/**/test.html")
	if !re.MatchString("src/foo/bar/test.html") {
		t.Error("** should match exact directories")
	}

	re = Globexp("*.js")
	if !re.MatchString(".config.js") {
		t.Error("* should match dot")
	}

	re = Globexp("*.js")
	if re.MatchString("foo/.config.js") {
		t.Error("* should not match directories")
	}

	re = Globexp("**/*.js")
	if !re.MatchString(".config.js") {
		t.Error("**/ slash should be optional")
	}

	re = Globexp("**/test.{html,js}")
	if !re.MatchString("src/test.html") || !re.MatchString("src/test.js") {
		t.Error("{} should match options")
	}

	re = Globexp("**/{{{{VERSION}}/*.foo")
	if !re.MatchString("src/{{VERSION}}/1.foo") {
		t.Error("{} should be escapable")
	}

	re = Globexp("public/**/*.uml")
	if !re.MatchString("public/{{VERSION}}/123/.4-5/a b/main-diagram.uml") {
		t.Error("should handle special chars")
	}

	re = Globexp("example/views/**/*.go.html")
	if !re.MatchString("example/views/admin/layout.go.html") {
		t.Error("should handle multiple subdirs admin")
	}
	if !re.MatchString("example/views/front/indexl.go.html") {
		t.Error("should handle multiple subdirs admin")
	}
}

func TestGlob(t *testing.T) {
	files, regexps, _ := Glob([]string{"./test/foo.txt"})
	if len(files) != 1 {
		t.Log("files", files)
		t.Error("should return file with no patterns")
	}
	if len(files) != len(regexps) {
		t.Error("Unequal amount of files and regexps")
	}

	files, _, _ = Glob([]string{"test/**/*.txt"})
	if len(files) != 5 {
		t.Log("files", files)
		t.Error("should return all txt files")
	}

	files, _, _ = Glob([]string{"test/**/*.txt", "!**/*sub1.txt"})
	if len(files) != 3 {
		t.Error("should return all go files but those suffixed with sub1.txt")
	}
	for _, file := range files {
		if strings.HasSuffix(file.Path, "sub1.txt") {
			t.Error("should have excluded a negated pattern")
		}
	}
}

func TestPatternRoot(t *testing.T) {
	s := PatternRoot("example/views/**/*.go.html")
	if s != "example/views" {
		t.Error("did not calculate root dir from pattern")
	}
}
