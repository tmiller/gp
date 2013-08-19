package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func printMessage() {
	gitBranchCommand := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	output, err := gitBranchCommand.Output()

	if err != nil {
		return
	}

	storyId := strings.TrimSpace(string(output))

	if story, err := pivotalTracker.FindStory(storyId); err == nil {
		fmt.Printf("[#%v] \n\n%v\n%v\n", story.Id, story.Name, story.Url)
	}
}
