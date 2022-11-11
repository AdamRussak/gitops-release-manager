package main

import (
	"fmt"
	"os"
	"regexp"
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
	org := os.Args[2]
	project := os.Args[3]
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
		var commentsArray []workItem
		for _, commit := range commits {
			if isCommitConvention(commit.Comment) {
				split := splitCommitMessage(commit.Comment)
				commentsArray = append(commentsArray, workItem{ServiceName: split[0], Name: split[2], Hash: split[1]})
			} else {
				commentsArray = append(commentsArray, workItem{ServiceName: "untracked", Name: commit.Comment, Hash: ""})
			}
		}
		sortingForMD := sortCommitsForMD(commentsArray, org, project)
		writeToMD(sortingForMD, t.Name, bumbedVersion)
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

func sortCommitsForMD(commits []workItem, org, project string) []string {
	var returnedString []string
	for c := range commits {
		testString, itemInArray := StringContains(returnedString, commits[c].ServiceName)
		workItem := getWorkItem(commits[c].Name)
		if commits[c].ServiceName == "untracked" {
			if testString {
				returnedString[itemInArray] = returnedString[itemInArray] + "| " + commits[c].Name + " | " + commits[c].Hash + " |\n"
			} else {
				returnedString = append(returnedString, "## "+commits[c].ServiceName+"\n"+KmdTable+"| "+commits[c].Name+" | "+commits[c].Hash+" |\n")

			}
		} else {
			if testString {
				returnedString[itemInArray] = returnedString[itemInArray] + "| " + "[" + commits[c].Name + "](" + KadoUrl + org + "/" + project + "/_workitems//edit/" + workItem + ")" + " | " + commits[c].Hash + " |\n"
			} else {
				returnedString = append(returnedString, "## "+commits[c].ServiceName+"\n"+KmdTable+"| "+"["+commits[c].Name+"]("+KadoUrl+org+"/"+project+"/_workitems//edit/"+workItem+")"+" | "+commits[c].Hash+" |\n")
			}
		}

	}
	return returnedString
}

func getWorkItem(s string) string {
	workItemRegex := regexp.MustCompile(`[0-9]+`)
	return workItemRegex.FindString(s)
}

func isCommitConvention(commit string) bool {
	isCommit := regexp.MustCompile(`\[([a-zA-Z]+(-[a-zA-Z]+)+)]\[[A-Za-z0-9]+]\[[^\]]*]`)
	return isCommit.MatchString(commit)
}

func StringContains(s []string, e string) (bool, int) {
	for a := range s {
		if strings.Contains(s[a], e) {
			return true, a
		}
	}
	return false, 0
}

func writeToMD(commentsArray []string, oldVersion, header string) {
	var writingOutput string
	for _, array := range commentsArray {
		writingOutput = writingOutput + array
	}
	header = "# " + oldVersion + "...." + header + "\n"
	fmt.Println(header + writingOutput)
	writeToFile([]byte(header + writingOutput))
}

func writeToFile(row []byte) {
	err := os.WriteFile("/home/coder/project/gitops-release-manager/tests/test.md", row, 0644)
	if err != nil {
		panic(err)
	}
}
