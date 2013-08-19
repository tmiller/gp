package main

import (
	"fmt"
	"github.com/tmiller/go-pivotal-tracker-api"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	workers int = 4
)

var pivotalIdPattern *regexp.Regexp = regexp.MustCompile(`\d{8,}`)

func getStoryIds(branches []string, storyIds chan<- *string) {
	for _, branch := range branches {
		if storyId := pivotalIdPattern.FindString(branch); storyId != "" {
			storyIds <- &storyId
		}

	}
	close(storyIds)
}

func getStories(storyIds chan *string, stories chan<- *pt.Story, finished chan<- bool) {
	for storyId := range storyIds {
		if story, err := pivotalTracker.FindStory(*storyId); err == nil {
			stories <- &story
		}
	}
	finished <- true
}

func monitorWorkers(stories chan<- *pt.Story, finished <-chan bool) {
	for i := 0; i < workers; i++ {
		<-finished
	}
	close(stories)
}

func printBranches() {
	output, err := exec.Command("git", "branch").Output()

	if err != nil {
		return
	}

	// Create a list of branches
	branches := strings.Split(strings.TrimRight(string(output), "\n"), "\n")

	storyIds := make(chan *string, 100)
	stories := make(chan *pt.Story, 100)
	finished := make(chan bool)

	go getStoryIds(branches, storyIds)
	go monitorWorkers(stories, finished)

	for i := 0; i < workers; i++ {
		go getStories(storyIds, stories, finished)
	}

	cachedStories := make(map[string]*pt.Story)
	for story := range stories {
		cachedStories[strconv.Itoa(story.Id)] = story
	}

	for _, branch := range branches {
		if storyId := pivotalIdPattern.FindString(branch); storyId != "" {
			if story, ok := cachedStories[storyId]; ok {
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
