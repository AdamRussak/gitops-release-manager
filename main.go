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

	// List all tag references, both lightweight tags and annotated tags

	tagrefs, _ := r.Tags()
	_ = tagrefs.ForEach(func(t *plumbing.Reference) error {
		fmt.Println(t)
		return nil
	})
	tags, _ := r.TagObjects()
	_ = tags.ForEach(func(t *object.Tag) error {
		fmt.Println(t.Name)
		fmt.Println(t.Tagger.When)
		fmt.Println(bumpVersion(t.Name))
		getCommits(r, t.Target)
		return nil
	})

}

func getCommits(r *git.Repository, tagHash plumbing.Hash) {
	until := time.Now()
	fromCommit, _ := r.CommitObject(tagHash)
	cIter, _ := r.Log(&git.LogOptions{Since: &fromCommit.Author.When, Until: &until})
	// ... just iterates over the commits, printing it
	_ = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		splitCommitMessage(c.Message)
		return nil
	})
}

func bumpVersion(tag string) string {
	test, _ := version.NewSemver(tag)
	log.Info("before: " + tag)
	segments := test.Segments()
	return fmt.Sprint(segments[0]) + "." + fmt.Sprint(segments[1]) + "." + fmt.Sprint(segments[2]+1)
}

func splitCommitMessage(comment string) {
	// var output []string
	splited := strings.SplitAfter(comment, "]")
	// regCommit := regexp.MustCompilePOSIX(`\[(.*?)\]`)
	// stringsArray := regCommit.FindAllStringSubmatch(comment, 3)
	// for _, s := range stringsArray {
	// 	s = strings.TrimSuffix(s, "]")
	// 	s = strings.TrimPrefix(s, "[")
	// 	output = append(output, s)
	// }
	fmt.Println(splited)
	os.Exit(2)
}
