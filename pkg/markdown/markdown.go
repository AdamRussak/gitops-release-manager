package markdown

import (
	"encoding/json"
	"fmt"
	"gitops-release-manager/pkg/core"
	"gitops-release-manager/pkg/provider"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// sort workitems after they have been splited in the gits pkg
func SortCommitsForMD(commits []WorkItem, org, project, pat string) ([]string, []string) {
	var workitemsID []string
	var wIOutput []WorkItemOutput
	for c := range commits {
		workItems := getWorkItem(commits[c].Name)
		for _, workItem := range workItems {
			if !checkDuplicateItem(workitemsID, workItem) {
				workitemsID = append(workitemsID, workItem)
				wIOutput = append(wIOutput, WorkItemOutput{itemID: c, workItem: workItem, Hash: commits[c].Hash})
			}
		}
	}
	return createMDStrings(commits, org, project, pat, wIOutput, workitemsID), workitemsID
}

func WriteToMD(commentsArray []string, oldVersion, header, path string) {
	var writingOutput string
	for _, array := range commentsArray {
		writingOutput = writingOutput + array
	}
	header = "# " + oldVersion + "..." + header + "\n"
	log.Debug(header + writingOutput)
	writeToFile([]byte(header+writingOutput), path)
}

func writeToFile(row []byte, path string) {
	err := os.WriteFile(path, row, 0644)
	core.OnErrorFail(err, "failed to save MD file")
}
func mdContains(s []string, e string) (bool, int) {
	for a := range s {
		if strings.Contains(s[a], "## "+e) {
			return true, a
		}
	}
	return false, 0
}
func getWorkItem(s string) []string {
	var ret []string
	if strings.Contains(s, "No related work items") {
		ret = append(ret, "No related work items")
	} else {
		WorkItems := regexp.MustCompile(`[0-9]+`)
		matches := WorkItems.FindAllString(s, 100)
		for i := range matches {
			log.Debugf("%s is a work Item", matches[i])
			ret = append(ret, matches[i])
		}

	}
	return ret
}
func checkDuplicateItem(workItemsArray []string, newWorkItem string) bool {
	for _, a := range workItemsArray {
		if a == newWorkItem {
			return true
		}
	}
	return false
}

func createMDStrings(commits []WorkItem, org, project, pat string, workItemOutput []WorkItemOutput, workitemsID []string) []string {
	var returnedString []string
	workItems := provider.GetWorkItemBatchStruct(org, project, pat, workitemsID)
	for c := range commits {
		// check if already exist in string for MD, returns bool and int (location of item in array if exist)
		testString, itemInArray := mdContains(returnedString, commits[c].ServiceName)
		relevantWI, wIExist := getReleventWI(workItems, commits[c], returnedString)
		// untracked is for items not in commit convention
		if commits[c].ServiceName == "untracked" || commits[c].ServiceName == "No related work items" {
			if testString {
				returnedString[itemInArray] = returnedString[itemInArray] + "| " + commits[c].Name + " | " + commits[c].Hash + " |\n"
			} else {
				returnedString = append(returnedString, "## "+commits[c].ServiceName+"\n"+KmdTable+"| "+commits[c].Name+" | "+commits[c].Hash+" |\n")
			}
		} else if wIExist {
			continue
		} else {
			// DebugPrintStruct(relevantWI)
			if relevantWI == 0 && wIExist {
				if testString {
					returnedString[itemInArray] = returnedString[itemInArray] + "| " + commits[c].Name + " | N/A | " + commits[c].Hash + " |\n"
				} else {
					returnedString = append(returnedString, "## "+commits[c].ServiceName+"\n"+KmdTable+"| N/A | "+commits[c].Name+" | "+commits[c].Hash+" |\n")
				}
			} else {
				if testString {
					returnedString[itemInArray] = returnedString[itemInArray] + " | " + fmt.Sprint(workItems.Value[relevantWI].ID) + "| " + "[" + workItems.Value[relevantWI].Fields.SystemTitle + "](" + KadoUrl + org + "/" + project + "/_workitems/edit/" + fmt.Sprint(workItems.Value[relevantWI].ID) + ")" + " | " + commits[c].Hash + " |\n"
				} else {
					returnedString = append(returnedString, "## "+commits[c].ServiceName+"\n"+KmdTable+" | " + fmt.Sprint(workItems.Value[relevantWI].ID) + " | "+"["+workItems.Value[relevantWI].Fields.SystemTitle+"]("+KadoUrl+org+"/"+project+"/_workitems/edit/"+fmt.Sprint(workItems.Value[relevantWI].ID)+")"+" | "+commits[c].Hash+" |\n")
				}
			}
		}
	}
	// returns an array of strings for MD and list of work Items for ADO
	return returnedString
}

func getReleventWI(wiStruct provider.BatchWorkItems, commit WorkItem, mdArray []string) (int, bool) {
	for v := range wiStruct.Value {
		if strings.Contains(commit.Name, fmt.Sprint(wiStruct.Value[v].ID)) {
			for _, md := range mdArray {
				if strings.Contains(md, fmt.Sprint(wiStruct.Value[v].ID)) {
					return 0, true
				}
			}
			return v, false
		}
	}
	return 0, true
}

func DebugPrintStruct(st interface{}) {
	out, _ := json.Marshal(st)
	log.Debug(string(out))
}
