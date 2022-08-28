package api

import (
	"bitbucket.org/yasnikoff/shuffler/app"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"strings"
	"bitbucket.org/yasnikoff/shuffler"
)

type AppData struct {
	App  app.App
	Repo shuffler.Repo
	Vars map[string]string
}

func (h *handler)getAppData(w http.ResponseWriter, r *http.Request) *AppData {

	vars := mux.Vars(r)
	name, ok := vars["repo"]

	if !ok {
		http.Error(w, fmt.Sprintf("repo name is not specified"), http.StatusBadRequest)
		return nil
	}
	app, ok := h.repos[name]
	if !ok {
		h.writeError(w, fmt.Errorf("no repo with the name %s", name), http.StatusBadRequest)
		return nil
	}

	return &AppData{App:app, Repo:app.Repo(), Vars:vars}
}

func (h *handler)handleRoot(w http.ResponseWriter, r *http.Request) {
	appLinks := make(map[string]string, len(h.repos))
	for name := range h.repos {
		appLinks[name] = fmt.Sprintf("%s/%s", strings.TrimSuffix(r.URL.String(), "/"), name)
	}
	h.write(w, appLinks)
}

func (h *handler)handleRepo(w http.ResponseWriter, r *http.Request) {
	a := h.getAppData(w, r)
	if a == nil {
		return
	}
	repo := a.App.Repo()

	resp := struct {
		Name string `json:"name"`
	}{
		Name:repo.Name(),
		//Recent:repo.Recent(),
	}
	h.write(w, resp)
}

func (h *handler)handleRepoIndexUpdate(w http.ResponseWriter, r *http.Request) {
	appData := h.getAppData(w, r)
	if appData == nil {
		return
	}

	data := struct {
		Path string `json:"path"`
	}{}
	if err := h.readJsonData(w, r, &data); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
	}

	repo := appData.App.Repo()

	switch data.Path{
	case "", ".":
		data.Path = repo.Root()
	}

	absPath, err:= repo.FromLibPath(data.Path)
	if err!=nil{
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := repo.UpdateIndexForPath(absPath); err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
	}
	h.write(w, nil)
}

func (h *handler)handleRepoNewProject(w http.ResponseWriter, r *http.Request) {
	// TODO: handleRepoNewProject tests
	appData := h.getAppData(w, r)
	if appData == nil {
		return
	}

	data := struct {
		Path string
		Name string
	}{}
	if err := h.readJsonData(w, r, &data); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	app := appData.App
	absPath, err:=appData.Repo.FromLibPath(data.Path)
	if err!=nil{
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	pr, err:= app.NewProject(absPath)
	if err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}

	pr.SetName(data.Name)
	if err:=pr.Save(); err != nil {
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}

	h.write(w, pr)

}
