package main

import (
	"flag"
	"fmt"
	"github.com/aliics/shortcut-cli"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	mostRecentFlag = flag.String("most-recent", "", "Most recent story for a user.")
)

func main() {
	cli := shortcut_cli.NewShortcut(
		shortcut_cli.WithShortcutToken(os.Getenv("SHORTCUT_TOKEN")),
		shortcut_cli.WithUrl("https://api.app.shortcut.com/api/v3"),
	)

	flag.Parse()

	if *mostRecentFlag != "" {
		stories, err := cli.SearchStories("owner:"+*mostRecentFlag, 1)
		if err != nil {
			panic(err)
		}

		story := stories.Data[0]
		owner := strings.ReplaceAll(strings.ToLower(url.PathEscape(*mostRecentFlag)), "%20", "-")
		id := strconv.Itoa(story.Id)
		name := strings.ReplaceAll(strings.ToLower(url.PathEscape(story.Name)), "%20", "-")

		fmt.Println(owner + "/sc-" + id + "/" + name)
	}
}
