package api

import (
	"net/http"
	"github.com/pkg/errors"
)

type TagsData struct {
	*projectData
	TagNames []string
}

func (h *handler)getTagsData(w http.ResponseWriter, r *http.Request) *TagsData {

	result := &TagsData{projectData:h.getProjectData(w, r)}
	if result.projectData == nil {
		return nil
	}
	if err := h.readJsonData(w, r, &result.TagNames); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return nil
	}

	return result
}

func (h *handler)handleTags(w http.ResponseWriter, r *http.Request) {
	p := h.getProjectData(w, r)
	if p == nil {
		return
	}
	h.write(w, p.Project.Tags.Get())
}

func (h *handler)handleTagsAdd(w http.ResponseWriter, r *http.Request) {

	data := h.getTagsData(w, r)
	if data == nil {
		return
	}

	data.Project.Tags.Add(data.TagNames...)
	if err := data.Project.Save(); err != nil {
		h.writeError(w, errors.Wrap(err, "can't save project"), http.StatusInternalServerError)
		return
	}
	h.write(w, data.Project.Tags.Get())
}

func (h *handler)handleTagsRemove(w http.ResponseWriter, r *http.Request) {

	data := h.getTagsData(w, r)
	if data == nil {
		return
	}

	data.Project.Tags.Remove(data.TagNames...)
	if err := data.Project.Save(); err != nil {
		h.writeError(w, errors.Wrap(err, "can't save project"), http.StatusInternalServerError)
		return
	}
	h.write(w, data.Project.Tags.Get())
}
