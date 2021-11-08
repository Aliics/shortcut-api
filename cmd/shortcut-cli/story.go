package main

import (
	"encoding/json"
	"fmt"
	shortcut_api "github.com/aliics/shortcut-api"
)

func showStories(members []shortcut_api.Member, stories []shortcut_api.Story) {
	showStory := func(members []shortcut_api.Member, story shortcut_api.Story) {
		whenOwned(members, story, func(member shortcut_api.Member) {
			if branchName {
				fmt.Println(story.GetBranchName(member))
			} else {
				output, err := json.Marshal(story)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(output))
			}
		})
	}

	if mostRecent {
		mostRecentStory := stories[0]
		for _, story := range stories {
			if story.UpdatedAt.After(mostRecentStory.UpdatedAt) {
				mostRecentStory = story
			}
		}

		showStory(members, mostRecentStory)
	} else {
		for _, story := range stories {
			showStory(members, story)
		}
	}
}

func whenOwned(members []shortcut_api.Member, story shortcut_api.Story, f func(shortcut_api.Member)) {
	for _, member := range members {
		for _, ownerId := range story.OwnerIds {
			if member.Id == ownerId {
				f(member)
			}
		}
	}
}

