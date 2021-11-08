package shortcut_api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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

type Entity struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
