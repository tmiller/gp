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

	branch := strings.TrimSpace(string(output))

  if storyId := pivotalIdPattern.FindString(branch); storyId != "" {
    if story, err := pivotalTracker.FindStory(storyId); err == nil {
      fmt.Printf("[#%v] \n\n%v\n%v\n", story.Id, story.Name, story.Url)
    }
  }
}
