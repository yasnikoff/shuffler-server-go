package api

import (
	"bitbucket.org/yasnikoff/shuffler"
	"net/http"
	"net/url"
	"bitbucket.org/yasnikoff/shuffler/repo"
	"strings"
	"bitbucket.org/yasnikoff/shuffler/types/id"
)

func (h *handler)query(w http.ResponseWriter, r *http.Request) {

	var q shuffler.ProjectQuery
	var result []id.ID

	appData := h.getAppData(w, r)
	if appData == nil {
		return
	}

	rawQuery := r.URL.Query().Get("query")
	rawQuery, err := url.QueryUnescape(rawQuery)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(rawQuery)==""{
		result=appData.App.Repo().Recent()
	}else{
		q = &repo.TagQuery{}
		if err := q.Parse(rawQuery); err != nil {
			h.writeError(w, err, http.StatusBadRequest)
			return
		}

		result, err = q.Execute(appData.App.Repo())
		if err != nil {
			h.writeError(w, err, http.StatusBadRequest)
		}
	}
	h.write(w, result)
}
