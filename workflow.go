package shortcut_api

import (
	"net/http"
)

type Workflows struct {
	Entity
	States []WorkflowState `json:"states"`
}

type WorkflowState struct {
	Entity
}

// ListWorkflows returns a list of the visible Workflows.
func (s Shortcut) ListWorkflows() ([]Workflows, error) {
	req, err := http.NewRequest(
		"GET",
		s.url+"/workflows",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var result []Workflows
	err = s.makeQuery(req, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
