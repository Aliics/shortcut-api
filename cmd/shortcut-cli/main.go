package main

import (
	"encoding/json"
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
		handleWorkflowStateFlag(api)
	}
}

func handleWorkflowStateFlag(api *shortcut_api.Shortcut) {
	workflows, err := api.ListWorkflows()
	if err != nil {
		panic(err)
	}

	var devWorkflowState shortcut_api.WorkflowState
	for _, wf := range workflows[0].States {
		if wf.Name == workflowState {
			devWorkflowState = wf
			break
		}
	}

	ms, err := api.ListMembers()
	if err != nil {
		panic(err)
	}

	var members []shortcut_api.Member
	for _, m := range ms {
		if m.Profile.MentionName == owner || owner == "" {
			members = append(members, m)
		}
	}

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

	var stories []shortcut_api.Story
	for _, story := range ss {
		var owned bool
		for _, ownerId := range story.OwnerIds {
			for _, member := range members {
				if ownerId == member.Id {
					owned = true
					break
				}
			}
		}

		if owned && story.WorkflowStateId == devWorkflowState.Id {
			stories = append(stories, story)
		}
	}

	for _, member := range members {
		showStories(member, stories)
	}
}

func showStories(member shortcut_api.Member, stories []shortcut_api.Story) {
	showStory := func(member shortcut_api.Member, story shortcut_api.Story) {
		if branchName {
			fmt.Println(story.GetBranchName(member))
		} else {
			output, err := json.Marshal(story)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(output))
		}
	}

	if mostRecent {
		showStory(member, stories[0])
	} else {
		for _, story := range stories {
			showStory(member, story)
		}
	}
}
