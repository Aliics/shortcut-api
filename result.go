package shortcut_api

import (
	"net/url"
	"strconv"
	"strings"
)

type StorySearchResult struct {
	Data []Story `json:"data"`
}

type Story struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (s Story) BranchName(owner string) string {
	escape := func(s string) string { return strings.ReplaceAll(strings.ToLower(url.PathEscape(s)), "%20", "-") }

	return escape(owner) + "/sc-" + strconv.Itoa(s.Id) + "/" + escape(s.Name)
}
