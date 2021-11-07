package shortcut_cli

type StorySearchResult struct {
	Data []Story `json:"data"`
}

type Story struct {
	Id int `json:"id"`
	Name string `json:"name"`
}
