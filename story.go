package shortcut_api

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type StorySearchResult struct {
	Data []Story `json:"data"`
}

type Story struct {
	Entity
	OwnerIds        []string `json:"owner_ids"`
	WorkflowStateId int      `json:"workflow_state_id"`
}

func (s Story) GetBranchName(owner Member) string {
	escape := func(s string) string { return strings.ReplaceAll(strings.ToLower(url.PathEscape(s)), "%20", "-") }

	return escape(owner.Profile.MentionName) + "/sc-" + strconv.Itoa(s.Id) + "/" + escape(s.Name)
}

func (s Shortcut) SearchStories(query string, pageSize int) (*StorySearchResult, error) {
	req, err := http.NewRequest(
		"GET",
		s.url+"/search/stories",
		bytes.NewBufferString(`{query":"`+query+`","page_size":`+strconv.Itoa(pageSize)+"}"),
	)
	if err != nil {
		return nil, err
	}

	result := &StorySearchResult{}
	err = s.makeQuery(req, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
