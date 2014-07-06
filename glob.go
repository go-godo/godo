package gosu

import (
	"bytes"
	//"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	// any rune but '/'
	NotSlash = "[^/]"

	// '*' any number of characters
	AnyRune = NotSlash + "*"

	ZeroOrMoreDirectories = "((?:[\\w\\.\\-]+\\/)*)"

	TrailingStarStar = "/**"

	SlashStarStarSlash = "/**/"
)

// Two consecutive asterisks ("**") in patterns matched against full pathname may
// have special meaning:
//
// A leading "**" followed by a slash means match in all directories. For example,
// "**/foo" matches file or directory "foo" anywhere, the same as pattern "foo".
// "**/foo/bar" matches file or directory "bar" anywhere that is directly under
// directory "foo".
//
// A trailing "/**" matches everything inside. For example, "abc/**" matches all
// files inside directory "abc", relative to the location of the .gitignore file,
// with infinite depth.
//
// A slash followed by two consecutive asterisks then a slash matches zero or more
// directories. For example, "a/**/b" matches "a/b", "a/x/b", "a/x/y/b" and so on.

// Glob creates a Regexp from extended glob pattern.
func Globexp(glob string) *regexp.Regexp {
	var re bytes.Buffer

	re.WriteString("^")

	i, inGroup := 0, false
	for i < len(glob) {
		r, w := utf8.DecodeRuneInString(glob[i:])

		switch r {
		default:
			re.WriteRune(r)

		case '\\', '$', '^', '+', '.', '(', ')', '=', '!', '|':
			re.WriteRune('\\')
			re.WriteRune(r)

		case '/':
			// TODO optimize later, string could be long
			rest := glob[i:]
			re.WriteRune('/')
			if strings.HasPrefix(rest, "/**/") {
				re.WriteString(ZeroOrMoreDirectories)
				w *= 4
			} else if rest == "/**" {
				re.WriteString(".*")
				w *= 3
			}

		case '?':
			re.WriteRune('.')

		case '[', ']':
			re.WriteRune(r)

		case '{':
			inGroup = true
			re.WriteRune('(')

		case '}':
			inGroup = false
			re.WriteRune(')')

		case ',':
			if inGroup {
				re.WriteRune('|')
			} else {
				re.WriteRune('\\')
				re.WriteRune(r)
			}

		case '*':
			rest := glob[i:]
			if strings.HasPrefix(rest, "**/") {
				re.WriteString(ZeroOrMoreDirectories)
				w *= 3
			} else {
				re.WriteString(AnyRune)
			}
		}

		i += w
	}

	re.WriteString("$")
	//log.Printf("regex string %s", re.String())
	return regexp.MustCompile(re.String())
}

// Glob returns files that match patterns.
//
// Special chars.
//
//   /**/   - match zero or more directories
//   {a,b}  - match a or b, no spaces
//   *      - match any non-separator char
//   ?      - match a single non-separator char
//   **/    - match any directory, start of pattern only
//   /**    - match any this directory, end of pattern only
//   !      - removes files from resultset, start of pattern only
//
func Glob(patterns []string) ([]*FileAsset, error) {
	// TODO very inefficient and unintelligent, optimize later

	m := map[string]*FileAsset{}

	for _, pattern := range patterns {
		remove := strings.HasPrefix(pattern, "!")
		if remove {
			pattern = pattern[1:]
			re := Globexp(pattern)
			for path := range m {
				if re.MatchString(path) {
					m[path] = nil
				}
			}
		} else {
			re := Globexp(pattern)
			root := patternRoot(pattern)
			if root == "" {
				Panicf("glob", "Cannot get root from pattern: %s", pattern)
			}
			chann := walk(root)
			for file := range chann {
				if re.MatchString(file.Path) {
					// TODO closure problem assigning &file
					tmp := file
					m[file.Path] = &tmp
				}
			}
		}
	}

	//log.Printf("m %v", m)
	keys := []*FileAsset{}
	for _, it := range m {
		if it != nil {
			keys = append(keys, it)
		}
	}

	return keys, nil
}

type FileAsset struct {
	os.FileInfo
	Path string
}

// hasMeta determines if a path has special chars used to be bild a Regexp.
func hasMeta(path string) bool {
	return strings.IndexAny(path, "*?[{") >= 0
}

// patternRoot gets a real directory root from a pattern. The directory
// returned is used as the start location for globbing.
func patternRoot(s string) string {
	// A negation does not walk the file system
	if strings.HasPrefix(s, "!") {
		return ""
	}
	// No directory in pattern
	parts := strings.Split(s, "/")
	if len(parts) == 1 {
		return "./"
	}
	// Build path until a dirname has a char used to build regex
	root, i, l := "", 0, len(parts)
	for i < l-1 {
		part := parts[i]
		if hasMeta(part) {
			break
		}
		if root == "" {
			root = part
		} else {
			root += "/" + part
		}
		i++
	}
	// Default to cwd
	if root == "" {
		root = "."
	}
	return root
}

// walk walks a directory starting at root returning all directories and files
// include those found in subdirectories.
func walk(root string) chan FileAsset {
	chann := make(chan FileAsset)
	go func() {
		// TODO replace. The consensus is the built-in Walk is horribly slow.
		filepath.Walk(root, func(path string, fi os.FileInfo, _ error) (err error) {
			result := FileAsset{FileInfo: fi, Path: path}
			chann <- result
			return
		})
		defer close(chann)
	}()
	return chann
}
