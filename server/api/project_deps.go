package api

import (
	"net/http"
	"bitbucket.org/yasnikoff/shuffler"
	"encoding/json"
)

type depsData struct {
	*projectData
	deps *shuffler.Dependencies
}

func (h *handler)getDepsData(w http.ResponseWriter, r *http.Request) *depsData {
	result := &depsData{
		projectData:h.getProjectData(w, r),
	}
	if result.projectData == nil {
		return nil
	}
	result.deps = result.Project.Deps
	return result
}

func (h *handler)handleDeps(w http.ResponseWriter, r *http.Request) {
	dd := h.getDepsData(w, r)
	if dd == nil {
		return
	}
	h.write(w, dd.deps)
}

type DepAddRequest struct {
	Path    string `json:"path"`
	Type    string `json:"type"`
	DepData json.RawMessage `json:"dep"`
}

func (h *handler)handleDepsAdd(w http.ResponseWriter, r *http.Request) {
	dd := h.getDepsData(w, r)
	if dd == nil {
		return
	}

	req := new(DepAddRequest)

	if err := h.readJsonData(w, r, req); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	newDep, err := shuffler.NewDep(req.Type)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(req.DepData, newDep); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	if err := dd.App.AddDep(dd.Project.ID().String(), req.Path, newDep.Descriptor().ID().String()); err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}
	h.write(w, nil)
}

type depRemoveRequest struct {
	Path string
}

func (h *handler)handleDepRemove(w http.ResponseWriter, r *http.Request) {
	// TODO: tests
	dd := h.getDepsData(w, r)
	if dd == nil {
		return
	}

	req := new(depRemoveRequest)

	if err := h.readJsonData(w, r, req); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	if err := dd.App.RemoveDep(dd.Project.ID().String(), req.Path); err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}
	h.write(w, nil)
}

type depUpdateRequest struct {
	Path string
}

func (h *handler)handleDepUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO: tests
	dd := h.getDepsData(w, r)
	if dd == nil {
		return
	}

	req := new(depUpdateRequest)

	if err := h.readJsonData(w, r, req); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		if err := dd.App.UpdateDepAll(dd.Project.ID().String(), true); err != nil {
			h.writeError(w, err, http.StatusInternalServerError)
			return
		}
	} else { // TODO: test when the is no a dep with the path
		if err := dd.App.UpdateDep(dd.Project.ID().String(), req.Path, true); err != nil {
			h.writeError(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.write(w, nil)
}


type depReplaceRequest struct {
	ID   string
	Path string
}

func (h *handler)handleDepReplace(w http.ResponseWriter, r *http.Request) {
	// TODO: tests
	dd := h.getDepsData(w, r)
	if dd == nil {
		return
	}

	req := new(depReplaceRequest)

	if err := h.readJsonData(w, r, req); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	if err := dd.App.ReplaceDep(dd.Project.ID().String(),req.Path ,req.ID); err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}

	h.write(w, nil)
}

