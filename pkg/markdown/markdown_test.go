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
func TestGetWorkItem(t *testing.T) {
	t.Run("Work Items ID detected", func(t *testing.T) {
		commit := "Released Work Items #234 #324"
		testArray := getWorkItem(commit)
		if len(testArray) == 2 {
			t.Fatalf(`getWorkItem() = %s,should have been %s`, fmt.Sprint(testArray), "2")
		}
	})
	t.Run("Work Items ID not detected", func(t *testing.T) {
		commit := "Released Work Items Not Listed"
		testArray := getWorkItem(commit)
		if len(testArray) == 0 {
			t.Fatalf(`getWorkItem() = %s ,should have been %s `, fmt.Sprint(testArray), "0")
		}
	})
}

// func TestCreateMDStrings(t *testing.T) {
// 	createMDStrings()

// }
func TestGetReleventWI(t *testing.T) {
	t.Run("Good Info", func(t *testing.T) {
		commit := WorkItem{Name: "1", ServiceName: "service A", Hash: "serviceAHash"}
		WorkItem := []string{"3", "4"}
		testInt, testBool := getReleventWI(wiStruct, commit, WorkItem)
		if testBool && testInt == 0 {
			t.Fatalf(`getReleventWI() = %s and %s,should have been %s and %s`, fmt.Sprint(testInt), fmt.Sprint(testBool), "1", "false")
		}
	})
	t.Run("Bad Info", func(t *testing.T) {
		commit := WorkItem{Name: "3", ServiceName: "service A", Hash: "serviceAHash"}
		WorkItem := []string{"3", "4"}
		testInt, testBool := getReleventWI(wiStruct, commit, WorkItem)
		if !testBool && testInt == 1 {
			t.Fatalf(`getReleventWI() = %s and %s,should have been %s and %s`, fmt.Sprint(testInt), fmt.Sprint(testBool), "0", "true")
		}
	})
}

func TestCheckDuplicateItem(t *testing.T) {
	t.Run("New Work Item", func(t *testing.T) {
		WorkItemArray := []string{"1", "2"}
		newWI := "3"
		testBool := checkDuplicateItem(WorkItemArray, newWI)
		if testBool {
			t.Fatalf(`checkDuplicateItem() = %s ,should have been %s`, fmt.Sprint(testBool), "false")
		}
	})
	t.Run("Duplicate Work Item", func(t *testing.T) {
		WorkItemArray := []string{"1", "2"}
		newWI := "2"
		testBool := checkDuplicateItem(WorkItemArray, newWI)
		if !testBool {
			t.Fatalf(`checkDuplicateItem() = %s ,should have been %s`, fmt.Sprint(testBool), "true")
		}
	})

}
