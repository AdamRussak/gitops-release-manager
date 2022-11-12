package markdown

import (
	"fmt"
	"os"
)

func WriteToMD(commentsArray []string, oldVersion, header string) {
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
