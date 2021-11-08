package shortcut_api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Shortcut is the primary client for this library.
// To initialize, using the NewShortcut constructor.
//
// Usage example:
//
//  api := shortcut_api.NewShortcut(
//	    shortcut_api.WithShortcutToken(os.Getenv("SHORTCUT_TOKEN")),
//	    shortcut_api.WithUrl("https://api.app.shortcut.com/api/v3"),
//	)
//  epics, err := api.ListEpics()
type Shortcut struct {
	url    string
	token  string
	client *http.Client
}

// NewShortcut is the constructor for Shortcut.
// NewShortcut uses the "Functional options" pattern
// to easy construction.
//
// Usage example:
//
//  api := shortcut_api.NewShortcut(
//	    shortcut_api.WithShortcutToken(os.Getenv("SHORTCUT_TOKEN")),
//	    shortcut_api.WithUrl("https://api.app.shortcut.com/api/v3"),
//	)
//  epics, err := api.ListEpics()
func NewShortcut(options ...ShortcutOption) *Shortcut {
	s := &Shortcut{url: "https://api.app.shortcut.com/api/v3", client: http.DefaultClient}
	for _, option := range options {
		option(s)
	}
	return s
}

type ShortcutOption func(shortcut *Shortcut)

// WithShortcutToken is meant to be used along with
// NewShortcut as a way to pass your "Shortcut-Token" header
// to shortcut api requests.
func WithShortcutToken(token string) ShortcutOption {
	return func(shortcut *Shortcut) {
		shortcut.token = token
	}
}

// WithUrl is meant to be used along with
// NewShortcut as a way to point to the shortcut api.
// By default, it points to "https://api.app.shortcut.com/api/v3".
func WithUrl(url string) ShortcutOption {
	return func(shortcut *Shortcut) {
		shortcut.url = url
	}
}

// WithHttpClient is meant to be used along with
// NewShortcut as a way to override the http client being used.
// By default, it uses http.DefaultClient.
func WithHttpClient(client *http.Client) ShortcutOption {
	return func(shortcut *Shortcut) {
		shortcut.client = client
	}
}

// makeQuery is a function which is the backbone of this library.
// All shortcut calls are made via this function.
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

// Entity is used to avoid redundancy in shortcut's api.
// That being said, their api is rather inconsistent.
// Look at how WorkflowState's id is an int. Whereas,
// Member's id is a string.
type Entity struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
