package markdown

import (
	"fmt"
	"gitops-release-manager/pkg/provider"
	"testing"
)

var wiStruct = provider.BatchWorkItems{Count: 3, Value: []provider.WorkItem{{ID: 1, Rev: 1, Fields: provider.WiFields{SystemID: 1, SystemTags: "1.0.0", SystemTitle: "work item 1", SystemWorkItemType: "user story"}}, {ID: 2, Rev: 2, Fields: provider.WiFields{SystemID: 1, SystemTags: "1.1.0", SystemTitle: "work item 2", SystemWorkItemType: "Bug"}}}}

// func TestSortCommitsForMD(t *testing.T) {
// 	SortCommitsForMD()
// }
// func TestWriteToMD(t *testing.T) {
// 	WriteToMD()
// }

// func TestWriteToFile(t *testing.T) {
// 	writeToFile()
// }
// func TestMdContains(t *testing.T) {
// 	mdContains()

// }
// func TestGetWorkItem(t *testing.T) {
// 	getWorkItem()

// }
// func TestCheckDuplicateItem(t *testing.T) {
// 	checkDuplicateItem()

// }
// func TestCreateMDStrings(t *testing.T) {
// 	createMDStrings()

// }
func TestGetReleventWI(t *testing.T) {
	commit := WorkItem{Name: "1", ServiceName: "service A", Hash: "serviceAHash"}
	WorkItem := []string{"3", "4"}
	testInt, testBool := getReleventWI(wiStruct, commit, WorkItem)
	if testBool && testInt == 0 {
		t.Fatalf(`getReleventWI() = %s and %s,should have been %s and %s`, fmt.Sprint(testInt), fmt.Sprint(testBool), "1", "false")
	}
}

// func TestMdAddLine(t *testing.T) {
// 	mdAddLine()
// }
