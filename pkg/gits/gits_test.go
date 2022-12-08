package gits

import (
	"fmt"
	"testing"
)

// TODO: create mock package: https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/
func TestIsVersionTag(t *testing.T) {
	t.Run("Good Info", func(t *testing.T) {
		testTag := "1.0.0"
		testIsVersioning := isVersionTag(testTag)
		if !testIsVersioning {
			t.Fatalf(`isVersionTag(%s) = %s,should have been %s`, testTag, fmt.Sprint(testIsVersioning), "true")
		}
	})

	t.Run("Bad Info", func(t *testing.T) {
		testTag := "not semver"
		testIsVersioning := isVersionTag(testTag)
		if testIsVersioning {
			t.Fatalf(`isVersionTag(%s) = %s,should have been %s`, testTag, fmt.Sprint(testIsVersioning), "false")
		}
	})
}
