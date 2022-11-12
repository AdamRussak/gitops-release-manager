package main

import (
	"giops-reelase-manager/pkg/core"
	"giops-reelase-manager/pkg/gits"
	"giops-reelase-manager/pkg/markdown"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
)

func main() {
	directory, org, project := os.Args[1], os.Args[2], os.Args[3]
	r, err := git.PlainOpen(directory)
	core.OnErrorFail(err, "faild to get git repo")
	gits.CheckOutBranch(r, "main")
	tags, _ := r.TagObjects()
	var tagsArray []string
	err = tags.ForEach(func(t *object.Tag) error {
		log.Infof("found tag %s", t.Name)
		tagsArray = append(tagsArray, t.Name)
		return nil
	})
	core.OnErrorFail(err, "err of ForEach tags process")
	latestTag := core.EvaluateVersion(tagsArray)
	bumbedVersion := core.BumpVersion(latestTag)
	log.Infof("New Version is: %s", bumbedVersion)
	latestTagObject, err := r.Tag(latestTag)
	core.OnErrorFail(err, "failed to get Tag Object")
	log.Info(latestTagObject.Hash().String())
	commits := gits.GetCommits(r, latestTagObject.Hash())
	var commentsArray []workItem
	for _, commit := range commits {
		if gits.IsCommitConvention(commit.Comment) {
			split := core.SplitCommitMessage(commit.Comment)
			commentsArray = append(commentsArray, workItem{ServiceName: split[0], Name: split[2], Hash: split[1]})
		} else {
			commentsArray = append(commentsArray, workItem{ServiceName: "untracked", Name: commit.Comment, Hash: ""})
		}
	}
	sortingForMD := sortCommitsForMD(commentsArray, org, project)
	markdown.WriteToMD(sortingForMD, latestTag, bumbedVersion)
}

func sortCommitsForMD(commits []workItem, org, project string) []string {
	var returnedString []string
	for c := range commits {
		testString, itemInArray := gits.StringContains(returnedString, commits[c].ServiceName)
		workItem := gits.GetWorkItem(commits[c].Name)
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
