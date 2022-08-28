package api

import (
	"bitbucket.org/yasnikoff/shuffler"
	"net/http"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"bitbucket.org/yasnikoff/shuffler/types/id"
)

type projectData struct {
	*AppData
	Project *shuffler.Project
}

func (h *handler)getProjectData(w http.ResponseWriter, r *http.Request) *projectData {
	appData := h.getAppData(w, r)
	if appData == nil {
		return nil
	}

	ID, ok := appData.Vars["id"]
	if !ok {
		http.Error(w, fmt.Sprintf("no project id specified"), http.StatusBadRequest)
		return nil
	}

	pr, err := appData.App.GetProjectByIDString(ID)
	if err != nil {
		h.writeError(w, errors.Wrapf(err, "can't get project in repo %q", appData.App.Repo().Name()), http.StatusNotFound)
		return nil
	}

	return &projectData{AppData:appData, Project:pr}
}

type ProjectData struct {
	Id     id.ID        `json:"id"`
	Name   string        `json:"name"`
	Path   string        `json:"path"`
	Origin id.ID        `json:"origin"`
}

// TODO: tests

func (h *handler)handleProject(w http.ResponseWriter, r *http.Request) {
	// TODO: do not return project, it holds too much data.
	// return only some basic data instead.
	p := h.getProjectData(w, r)
	if p == nil {
		return
	}

	if addToRecent := r.FormValue("addToRecent"); addToRecent != "" {
		if addToRecent == "true" {
			p.App.Repo().AddToRecent(p.Project.ID());
		}
	}

	libPath, err := p.App.Repo().ToLibPath(p.Project.Path())
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}

	data := ProjectData{
		Id:p.Project.Id,
		Name:p.Project.NameData,
		Path:libPath,
		Origin:p.Project.Origin,
	}

	h.write(w, data)
}

func (h *handler)handleProjectSave(w http.ResponseWriter, r *http.Request) {
	p := h.getProjectData(w, r)
	if p == nil {
		return
	}

	prData := struct {
		Name   string
		Origin string
	}{}
	err := h.readJsonData(w, r, &prData)
	if err != nil {
		return
	}

	// Name
	p.Project.SetName(prData.Name)

	// Origin
	origin, err := id.Parse(prData.Origin)
	if err != nil {
		h.writeError(w, errors.Wrap(err, "can't parse project origin"), http.StatusBadRequest)
	}
	p.Project.Origin = origin

	//Save
	if err := p.Project.Save(); err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}

	h.write(w, nil)

}

func (h *handler)handleProjectPath(w http.ResponseWriter, r *http.Request) {
	d := h.getProjectData(w, r)
	if d == nil {
		return
	}
	path, err := d.App.Repo().ToLibPath(d.Project.Path())
	if err != nil {
		h.writeError(w, errors.Wrap(err, "can't get project lib path"), http.StatusBadRequest)
	}

	h.write(w, path)
}
func (h *handler)handleProjectPathSave(w http.ResponseWriter, r *http.Request) {

	d := h.getProjectData(w, r)
	if d == nil {
		return
	}

	var data struct {
		Path string
	}

	if err := h.readJsonData(w, r, &data); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
	}

	path, err := d.App.Repo().FromLibPath(data.Path)
	if err != nil {
		h.writeError(w, errors.Wrap(err, "can't convert new path to lib path"), http.StatusBadRequest)
	}

	if err != nil {
		h.writeError(w, errors.Wrapf(err, "can't save project %s", d.Project), http.StatusBadRequest)
		return
	}
	if path != d.Project.Path() {

		if _, err := d.App.MoveProject(d.Project.ID(), path); err != nil {
			h.writeError(w, err, http.StatusInternalServerError)
		}
	}

	h.write(w, nil)
}
func (h *handler)handleProjectDuplicate(w http.ResponseWriter, r *http.Request) {

	d := h.getProjectData(w, r)
	if d == nil {
		return
	}

	var data struct {
		// new name and path
		Path string `json:"path"`
		Name string `json:"name"`
	}

	if err := h.readJsonData(w, r, &data); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
	}

	path, err := d.App.Repo().FromLibPath(data.Path)
	if err != nil {
		h.writeError(w, errors.Wrapf(err, "can't duplicate project %s", d.Project.ID()),
			http.StatusBadRequest)
	}

	newPr, err := d.App.DuplicateProject(d.Project.ID().String(), data.Name, path)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
	}

	h.write(w,newPr.ID())
}

var isLocalHostRegEx = regexp.MustCompile("^localhost[:].")

func (h *handler)handleProjectOpen(w http.ResponseWriter, r *http.Request) {

	if !isLocalHostRegEx.MatchString(r.Host) {
		h.writeError(w, errors.New("not localhost"), http.StatusBadRequest)
		return
	}

	d := h.getProjectData(w, r)
	if d == nil {
		return
	}

	if err := d.App.OpenProjectFolder(d.Project.ID().String()); err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}

	h.write(w, nil)
}
