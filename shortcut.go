package shortcut_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Shortcut struct {
	url    string
	token  string
	client *http.Client
}

func NewShortcut(options ...ShortcutOption) *Shortcut {
	s := &Shortcut{client: http.DefaultClient}
	for _, option := range options {
		option(s)
	}
	return s
}

func (s Shortcut) SearchStories(query string, pageSize int) (*StorySearchResult, error) {
	req, err := http.NewRequest(
		"GET",
		s.url+"/search/stories",
		bytes.NewBufferString("{\"query\":\""+query+"\",\"page_size\":"+strconv.Itoa(pageSize)+"}"),
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

func (s Shortcut) makeQuery(req *http.Request, t interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Shortcut-Token", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer func() { err = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return errors.New("request failed [" + resp.Status + "]: " + string(rb))
	}

	return json.Unmarshal(rb, t)
}

type ShortcutOption func(shortcut *Shortcut)

func WithShortcutToken(token string) ShortcutOption {
	return func(shortcut *Shortcut) {
		shortcut.token = token
	}
}

func WithUrl(url string) ShortcutOption {
	return func(shortcut *Shortcut) {
		shortcut.url = url
	}
}

func WithHttpClient(client *http.Client) ShortcutOption {
	return func(shortcut *Shortcut) {
		shortcut.client = client
	}
}
