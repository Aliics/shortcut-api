package shortcut_api

import "net/http"

type Member struct {
	Id      string        `json:"id"`
	Profile MemberProfile `json:"profile"`
}

type MemberProfile struct {
	Entity
	Id          string `json:"id"`
	MentionName string `json:"mention_name"`
}

// ListMembers returns a list of the visible Members.
func (s Shortcut) ListMembers() ([]Member, error) {
	req, err := http.NewRequest(
		"GET",
		s.url+"/members",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var result []Member
	err = s.makeQuery(req, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
