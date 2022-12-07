package gits

import (
	"fmt"
	"testing"
)

func TestIsVersionTag(t *testing.T) {
	testTag := "1.0.0"
	testIsVersioning := isVersionTag(testTag)
	if !testIsVersioning {
		t.Fatalf(`isVersionTag(%s) = %s,should have been %s`, testTag, fmt.Sprint(testIsVersioning), "true")
	}
}
