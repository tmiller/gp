package main

import (
	"fmt"
	"github.com/tmiller/go-pivotal-tracker-api"
	"os/exec"
	"regexp"
	"strings"
)

func printBranches() {
	output, err := exec.Command("git", "branch").Output()

	if err != nil {
		return
	}

	// Create a list of branches
	branches := strings.Split(strings.TrimRight(string(output), "\n"), "\n")

	// Build string[] of nothing but matching story ids
	pivotalIdPattern := regexp.MustCompile(`\d{8,}`)
	var storyIds []string
	for _, branch := range branches {
		if storyId := pivotalIdPattern.FindString(branch); storyId != "" {
			storyIds = append(storyIds, storyId)
		}
	}

	// Partion the story ids into 4 separate lists
	var storyIdPartitions [partitions][]string
	for i, storyId := range storyIds {
		storyIdPartitions[i%partitions] =
			append(storyIdPartitions[i%partitions], storyId)
	}

	pivotalResult := make(chan pt.Story, partitions)
	done := make(chan bool, partitions)

	for i := 0; i < len(storyIdPartitions); i++ {
		go func(part int) {
			for _, storyId := range storyIdPartitions[part] {
				if story, err := pivotalTracker.FindStory(storyId); err == nil {
					pivotalResult <- story
				}
			}
			done <- true
		}(i)
	}

	go func() {
		for i := 0; i < partitions; i++ {
			<-done
		}
		close(pivotalResult)
	}()

	stories := make(map[string]pt.Story)
	for story := range pivotalResult {
		stories[story.Id] = story
	}

	for _, branch := range branches {
		if storyId := pivotalIdPattern.FindString(branch); storyId != "" {
			if story, ok := stories[storyId]; ok {
				fmt.Printf(
					"%s [%s] %s (%s)\n",
					branch,
					story.State(),
					story.Name,
					story.Url)
			} else {
				fmt.Println(branch)
			}
		} else {
			fmt.Println(branch)
		}
	}
}
