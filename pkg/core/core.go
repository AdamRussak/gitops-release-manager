package core

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
)

func OnErrorFail(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s\n", message, err)
	}
}

// evaluate latest version from addon version list
func EvaluateVersion(list []string) string {
	var latest string
	for _, v := range list {
		var lt *version.Version
		var err error
		v1, err := version.NewVersion(v)
		OnErrorFail(err, "Error Evaluating Version")
		if latest == "" {
			lt, err = version.NewVersion("0.0")
		} else {
			lt, err = version.NewVersion(latest)
		}
		OnErrorFail(err, "Error Evaluating Version")
		// Options availabe
		if v1.GreaterThan(lt) {
			latest = v
		} // GreaterThen
	}
	return latest
}

func BumpVersion(tag string) string {
	log.Debug("Starting BumpVersion() function")
	test, _ := version.NewSemver(tag)
	log.Info("before: " + tag)
	segments := test.Segments()
	return fmt.Sprint(segments[0]) + "." + fmt.Sprint(segments[1]) + "." + fmt.Sprint(segments[2]+1)
}

func SplitCommitMessage(comment string) []string {
	var output []string
	splited := strings.SplitAfter(comment, "]")
	for _, s := range splited {
		s = strings.TrimSuffix(s, "]")
		s = strings.TrimPrefix(s, "[")
		if s != "" {
			output = append(output, s)
		}
	}
	return output
}

func IsSemVer(tag string) bool {
	isCommit := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	boolinas := isCommit.MatchString(tag)
	return boolinas
}

func ValidateIsDIrectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		OnErrorFail(err, "Error checking path:")
	}
	if fileInfo.IsDir() {
		log.Infof("%s is a valid directory\n", path)
		return true
	} else {
		log.Warningf("%s is not a directory\n", path)
		return false
	}
}
