package api

import (
	"bitbucket.org/yasnikoff/shuffler/actions"
	"net/http"
)

type actionsData struct {
	*projectData
	actions *actions.Collection
}

func (h *handler)getActionsData(w http.ResponseWriter, r *http.Request) *actionsData {
	result := &actionsData{
		projectData:h.getProjectData(w, r),
	}
	if result.projectData == nil {
		return nil
	}
	result.actions = result.Project.Actions
	return result
}

