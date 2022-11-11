package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
)

func main() {
	directory := os.Args[1]
	r, _ := git.PlainOpen(directory)

	tagrefs, _ := r.Tags()
	_ = tagrefs.ForEach(func(t *plumbing.Reference) error {
		fmt.Println(t)
		return nil
	})
	tags, _ := r.TagObjects()
	_ = tags.ForEach(func(t *object.Tag) error {
		fmt.Println(t.Name)
		bumbedVersion := bumpVersion(t.Name)
		commits := getCommits(r, t.Target)
		var commentsArray [][]string
		for _, commit := range commits {
			commentsArray = append(commentsArray, splitCommitMessage(commit.Comment))
		}
		writeToMD(commentsArray, bumbedVersion)
		return nil
	})

}

func getCommits(r *git.Repository, tagHash plumbing.Hash) []commit {
	var comments []commit
	until := time.Now()
	fromCommit, _ := r.CommitObject(tagHash)
	cIter, _ := r.Log(&git.LogOptions{Since: &fromCommit.Author.When, Until: &until})
	// ... just iterates over the commits, printing it
	_ = cIter.ForEach(func(c *object.Commit) error {
		comments = append(comments, commit{Hash: c.Hash.String(), Comment: c.Message})
		return nil
	})
	return comments
}

func bumpVersion(tag string) string {
	test, _ := version.NewSemver(tag)
	log.Info("before: " + tag)
	segments := test.Segments()
	return fmt.Sprint(segments[0]) + "." + fmt.Sprint(segments[1]) + "." + fmt.Sprint(segments[2]+1)
}

func splitCommitMessage(comment string) []string {
	var output []string
	splited := strings.SplitAfter(comment, "]")
	for _, s := range splited {
		s = strings.TrimSuffix(s, "]")
		s = strings.TrimPrefix(s, "[")
		output = append(output, s)
	}
	return output
}

// TODO: create struct for each service
// TODO: create structerd comments: product version,service, work-item + work item description, commit id
func writeToMD(commentsArray [][]string, header string) {
	var writingOutput string
	for _, array := range commentsArray {
		var commentMD = "- [ ] "
		for _, comment := range array {
			commentMD = commentMD + comment + " "
		}
		writingOutput = writingOutput + commentMD + "\n"
	}
	header = "# " + header + "\n"
	fmt.Println(header + writingOutput)
	writeToFile([]byte(header + writingOutput))
}

func writeToFile(row []byte) {
	err := os.WriteFile("/home/coder/project/gitops-release-manager/tests/test.md", row, 0644)
	if err != nil {
		panic(err)
	}
}
