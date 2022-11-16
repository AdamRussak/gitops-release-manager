package main

/*
Copyright © 2022 Adam Russak adam+russak@gmail.com
*/

import (
	"giops-reelase-manager/pkg/cmd"
)

// TODO: add creation of new tag : https://github.com/go-git/go-git/blob/master/_examples/tag-create-push/main.go
// TODO: add push to new tag after creation (should be after succesfull finish of process)

func main() {
	cmd.Execute()
}
