package glob

import (
	"strings"
	"testing"
)

func TestEffectiveCriteria(t *testing.T) {
	result, _ := EffectiveCriteria("xtest/*.txt", "xtest2/**/*.html", "xtest/*.js", "!xtest/*.html")

	if len(result.Roots()) != 2 {
		t.Error("expected 2 items in result set")
	}

	success := 0
	for _, c := range result.Items {
		if strings.HasSuffix(c.Root, "/xtest") {
			success++
		}
		if strings.HasSuffix(c.Root, "/xtest2") {
			success++
		}
	}

	if success != 2 {
		t.Error("should calc effective criteria")
	}

}
