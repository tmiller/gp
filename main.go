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

	pivotalIdPattern := regexp.MustCompile(`\d{8,}`)
	branches := strings.Split(strings.TrimRight(string(output), "\n"), "\n")

	for _, branch := range branches {
		var storySummary string
		if storyId := pivotalIdPattern.FindString(branch); storyId != "" {
			if story, err := pivotalTracker.FindStory(storyId); err == nil {
				storySummary = fmt.Sprintf("[%s] %s (%s)", story.State(), story.Name, story.Url)
			}
		}
		fmt.Println(branch, storySummary)
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
		fmt.Printf("[#%d] \n\n%s\n%s\n", story.Id, story.Name, story.Url)
	}
}
