package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aliics/shortcut-api"
	"os"
)

var (
	mostRecent string
	branchName bool
)

func main() {
	flag.StringVar(&mostRecent, "most-recent", "", "Most recent story for a user.")
	flag.BoolVar(&branchName, "branch-name", false, "Show stories as branch names.")
	flag.Parse()

	api := shortcut_api.NewShortcut(
		shortcut_api.WithShortcutToken(os.Getenv("SHORTCUT_TOKEN")),
		shortcut_api.WithUrl("https://api.app.shortcut.com/api/v3"),
	)

	if mostRecent != "" {
		handleMostRecent(api)
	}
}

func handleMostRecent(api *shortcut_api.Shortcut) {
	stories, err := api.SearchStories("owner:"+mostRecent, 1)
	if err != nil {
		panic(err)
	}

	story := stories.Data[0]

	if branchName {
		fmt.Println(story.BranchName(mostRecent))
	} else {
		output, err := json.Marshal(story)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(output))
	}
}
