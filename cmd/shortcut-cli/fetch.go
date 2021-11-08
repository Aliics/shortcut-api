package main

import "github.com/aliics/shortcut-api"

// Functions in this file return receiving channels
// which wrap goroutine calls to the Shortcut Api.
// This is to avoid long blocking calls on the main
// goroutine.

func fetchMembers(api *shortcut_api.Shortcut) <-chan []shortcut_api.Member {
	membersCh := make(chan []shortcut_api.Member)
	go func() {
		defer close(membersCh)

		ms, err := api.ListMembers()
		if err != nil {
			panic(err)
		}

		var members []shortcut_api.Member
		for _, m := range ms {
			if m.Profile.MentionName == owner || owner == "" {
				members = append(members, m)
			}
		}

		membersCh <- members
	}()

	return membersCh
}

func fetchWorkflowStates(api *shortcut_api.Shortcut) <-chan shortcut_api.WorkflowState {
	workflowStateCh := make(chan shortcut_api.WorkflowState, 1)
	go func() {
		defer close(workflowStateCh)

		workflows, err := api.ListWorkflows()
		if err != nil {
			panic(err)
		}

		var devWorkflowState shortcut_api.WorkflowState
		for _, wf := range workflows[0].States {
			if wf.Name == workflowState {
				devWorkflowState = wf
				break
			}
		}

		workflowStateCh <- devWorkflowState
	}()

	return workflowStateCh
}
