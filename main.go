package main

import (
	"flag"
	"fmt"
	"github.com/tmiller/go-pivotal-tracker-api"
	"io/ioutil"
	"os"
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
