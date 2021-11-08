package main

import (
	"flag"
	"fmt"
	"github.com/aliics/shortcut-api"
	"os"
)

var (
	workflowState string
	owner         string
	mostRecent    bool
	branchName    bool
)

func main() {
	flag.StringVar(&workflowState, "workflow-state", "", "Select stories in a workflow state.")
	flag.StringVar(&owner, "owner", "", "Owner query selection.")
	flag.BoolVar(&mostRecent, "most-recent", false, "Show the most recent story.")
	flag.BoolVar(&branchName, "branch-name", false, "Show stories as branch names.")
	flag.Parse()

	api := shortcut_api.NewShortcut(
		shortcut_api.WithShortcutToken(os.Getenv("SHORTCUT_TOKEN")),
		shortcut_api.WithUrl("https://api.app.shortcut.com/api/v3"),
	)

	if workflowState != "" {
		handleWorkflowStateArg(api)
	} else if mostRecent || owner != "" {
		handleArbitraryFetch(api)
	} else if branchName {
		fmt.Println("argument --branch-name should be used with --owner, --most-recent, or --workflow-state")
	}
}

func handleWorkflowStateArg(api *shortcut_api.Shortcut) {
	workflowStateCh := fetchWorkflowStates(api)
	membersCh := fetchMembers(api)

	// All non-completed epics and their stories.
	epics, err := api.QueryEpics(func(epic shortcut_api.Epic) bool {
		return !epic.Completed
	})
	if err != nil {
		panic(err)
	}

	ss, err := api.ListStoriesForEpics(epics)
	if err != nil {
		panic(err)
	}

	devWorkflowState := <-workflowStateCh
	members := <-membersCh

	// Stories which match the workflow state and owner predicates.
	var stories []shortcut_api.Story
	for _, story := range ss {
		var owned bool
		whenOwned(members, story, func(_ shortcut_api.Member) {
			owned = true
		})

		if owned && story.WorkflowStateId == devWorkflowState.Id {
			stories = append(stories, story)
		}
	}

	showStories(members, stories)
}

func handleArbitraryFetch(api *shortcut_api.Shortcut) {
	membersCh := fetchMembers(api)

	// All non-completed epics and their stories.
	epics, err := api.QueryEpics(func(epic shortcut_api.Epic) bool {
		return !epic.Completed
	})
	if err != nil {
		panic(err)
	}

	ss, err := api.ListStoriesForEpics(epics)
	if err != nil {
		panic(err)
	}

	members := <-membersCh

	// Stories which match the workflow state and owner predicates.
	var stories []shortcut_api.Story
	for _, story := range ss {
		var owned bool
		whenOwned(members, story, func(_ shortcut_api.Member) {
			owned = true
		})

		if owned {
			stories = append(stories, story)
		}
	}

	showStories(members, stories)
}
