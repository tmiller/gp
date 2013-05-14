package main

import (
	"flag"
	"fmt"
	"github.com/tmiller/go-pivotal-tracker-api"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	partitions int = 4
)

var pivotalTracker pt.PivotalTracker
var branchesFlag bool
var messageFlag bool

func main() {

	parseFlags()
	initPivotalTracker()

	switch {
	case messageFlag:
		printMessage()
	case branchesFlag:
		printBranches()
	}
}

func parseFlags() {

	flag.BoolVar(&branchesFlag, "branches", true, "Print branches with story information")
	flag.BoolVar(&messageFlag, "message", false, "Generate a commit message from story information")

	flag.BoolVar(&branchesFlag, "b", true, "Print branches with story information")
	flag.BoolVar(&messageFlag, "m", false, "Generate a commit message from story information")

	flag.Parse()
}

func initPivotalTracker() {
	configFilePath := os.ExpandEnv("${HOME}/.pivotal_tracker_api_key")
	configFile, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		fmt.Printf("Please put your Pivotal Tracker API key in %s\n", configFilePath)
		os.Exit(1)
	}

	pivotalTrackerApiKey := strings.TrimSpace(string(configFile))
	pivotalTracker = pt.PivotalTracker{pivotalTrackerApiKey}
}

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
	go func() {
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

func printMessage() {
	gitBranchCommand := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	output, err := gitBranchCommand.Output()

	if err != nil {
		return
	}

	storyId := strings.TrimSpace(string(output))

	if story, err := pivotalTracker.FindStory(storyId); err == nil {
		fmt.Printf("[#%s] \n\n%s\n%s\n", story.Id, story.Name, story.Url)
	}
}
