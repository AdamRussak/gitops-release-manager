package markdown

import (
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
	if strings.Contains(s, KlogResp) {
		ret = append(ret, KlogResp)
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
		log.Debugf("the commit string: %s", commits[c])
		// check if already exist in string for MD, returns bool and int (location of item in array if exist)
		log.Debugf("The returnString: %s", returnedString)
		log.Debugf("The ServiceName: %s", commits[c].ServiceName)
		testString, itemInArray := mdContains(returnedString, commits[c].ServiceName)
		log.Debugf("The testString: %t", testString)
		log.Debugf("The itemInArray: %x", itemInArray)
		relevantWI, wIExist := getReleventWI(workItems, commits[c], returnedString)
		// untracked is for items not in commit convention
		log.Debugf("The relevantWI: %x", relevantWI)
		log.Debugf("The wIExist: %t", wIExist)
		adPath := org + "/" + project
		log.Debugf("The adPath: %s", adPath)
		log.Debugf("Added String: %s", commits[c])
		// log.Debugf("workItems ID: %x", string(workItems.Value))
		returnedString = mdAddLine(itemInArray, commits[c], wIExist, testString, returnedString, workItems.Value[relevantWI], adPath)
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

func mdAddLine(itemInArray int, commit WorkItem, wIExist, testString bool, returnedString []string, workItem provider.WorkItem, adPath string) []string {
	log.Debugf("Starting the mdAddLine() func")
	if commit.ServiceName == "untracked" || commit.ServiceName == KlogResp {
		comentName := strings.ReplaceAll(commit.Name, "\n", " ")
		log.Debug(commit.Name)
		if testString {
			returnedString[itemInArray] = returnedString[itemInArray] + "| NA | NA |" + comentName + " | NA |\n"
		} else {
			returnedString = append(returnedString, "## "+commit.ServiceName+"\n"+KmdTable+"| NA | NA |"+comentName+"| NA |\n")
		}
	} else if !wIExist {
		if testString {
			returnedString[itemInArray] = returnedString[itemInArray] + " | " + fmt.Sprint(workItem.ID) + " | " + workItem.Fields.SystemWorkItemType + " | " + "[" + workItem.Fields.SystemTitle + "](" + KadoUrl + adPath + "/_workitems/edit/" + fmt.Sprint(workItem.ID) + ")" + " | " + commit.Hash + " |\n"
		} else {
			returnedString = append(returnedString, "## "+commit.ServiceName+"\n"+KmdTable+" | "+fmt.Sprint(workItem.ID)+" | "+workItem.Fields.SystemWorkItemType+" | "+"["+workItem.Fields.SystemTitle+"]("+KadoUrl+adPath+"/_workitems/edit/"+fmt.Sprint(workItem.ID)+")"+" | "+commit.Hash+" |\n")
		}
	}
	log.Debugf("The output from the mdAddLine(): %s", returnedString)
	return returnedString
}
func HasMDSuffix(file, extension string) string {
	if !strings.HasSuffix(file, extension) {
		file += "." + extension
	}
	return file
}
