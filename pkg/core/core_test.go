package core

import (
	"fmt"
	"testing"
)

var testVersion = "1.1.1"

func TestEvaluateVersion(t *testing.T) {
	versionTestList := []string{"1.0.0", "1.1.1", "7.9.8", "8.2.2", "2.1.3"}
	latestVersion := EvaluateVersion(versionTestList)
	if latestVersion != "8.2.2" {
		t.Fatalf(`EvaluateVersion(%s) = %s,should have been %s`, versionTestList, latestVersion, "8.2.2")
	}
}
func TestBumpVersion(t *testing.T) {
	checkBump := BumpVersion(testVersion)
	if checkBump != "1.1.2" {
		t.Fatalf(`BumpVersion(%s) = %s,should have been %s`, testVersion, checkBump, "1.1.2")
	}
}
func TestSplitCommitMessage(t *testing.T) {
	commmit := "[test-service][da9e223d0b96bbb6sdf0ab5ddfb53318ab97][ Related work items: #204223, #204224 ]"
	splitOutput := SplitCommitMessage(commmit)
	if len(splitOutput) != 3 {
		for a := range splitOutput {
			t.Errorf(`item %x = %s`, a, splitOutput[a])
		}
		t.Fatalf(`SplitCommitMessage(%s) = %s,should have been %s`, commmit, fmt.Sprint(len(splitOutput)), "3")
	}
}
func TestIsSemVer(t *testing.T) {
	semVar := IsSemVer(testVersion)
	if !semVar {
		t.Fatalf(`IsSemVer(%s) = %s,should have been %s`, testVersion, fmt.Sprint(semVar), "true")
	}
}
