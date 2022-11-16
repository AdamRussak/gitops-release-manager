package markdown

import (
	"giops-reelase-manager/pkg/core"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func SortCommitsForMD(commits []WorkItem, org, project, pat, newVersion string) ([]string, []string) {
	var returnedString []string
	var workitemsID []string
	for c := range commits {
		testString, itemInArray := stringContains(returnedString, commits[c].ServiceName)
		workItem := getWorkItem(commits[c].Name)
		if commits[c].ServiceName == "untracked" {
			if testString {
				returnedString[itemInArray] = returnedString[itemInArray] + "| " + commits[c].Name + " | " + commits[c].Hash + " |\n"
			} else {
				returnedString = append(returnedString, "## "+commits[c].ServiceName+"\n"+KmdTable+"| "+commits[c].Name+" | "+commits[c].Hash+" |\n")

			}
		} else {
			workitemsID = append(workitemsID, workItem)
			if testString {
				returnedString[itemInArray] = returnedString[itemInArray] + "| " + "[" + commits[c].Name + "](" + KadoUrl + org + "/" + project + "/_workitems//edit/" + workItem + ")" + " | " + commits[c].Hash + " |\n"
			} else {
				returnedString = append(returnedString, "## "+commits[c].ServiceName+"\n"+KmdTable+"| "+"["+commits[c].Name+"]("+KadoUrl+org+"/"+project+"/_workitems//edit/"+workItem+")"+" | "+commits[c].Hash+" |\n")
			}
		}

	}

	return returnedString, workitemsID
}

func WriteToMD(commentsArray []string, oldVersion, header string) {
	var writingOutput string
	for _, array := range commentsArray {
		writingOutput = writingOutput + array
	}
	header = "# " + oldVersion + "...." + header + "\n"
	log.Debug(header + writingOutput)
	writeToFile([]byte(header + writingOutput))
}

func writeToFile(row []byte) {
	err := os.WriteFile("/home/coder/project/gitops-release-manager/tests/test.md", row, 0644)
	core.OnErrorFail(err, "failed to save MD file")
}
func stringContains(s []string, e string) (bool, int) {
	for a := range s {
		if strings.Contains(s[a], e) {
			return true, a
		}
	}
	return false, 0
}
func getWorkItem(s string) string {
	workItemRegex := regexp.MustCompile(`[0-9]+`)
	return workItemRegex.FindString(s)
}
