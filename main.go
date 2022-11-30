package main

/*
Copyright © 2022 Adam Russak adam+russak@gmail.com
*/

import (
	"gitops-release-manager/pkg/cmd"
)

// AGREED FLOW:
// TODO: add creation of new tag : https://github.com/go-git/go-git/blob/master/_examples/tag-create-push/main.go

func main() {
	cmd.Execute()
}
