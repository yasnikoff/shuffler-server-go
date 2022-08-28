package api

import (
	"net/http"
	"io/ioutil"
	"log"
	"path/filepath"
)

func (h *handler)handlePreview(w http.ResponseWriter, r *http.Request) {

	p := h.getProjectData(w, r)
	if p == nil {
		return
	}

	previewPath := filepath.Join(p.Project.Path(), "preview.jpg")

	imgData, err := ioutil.ReadFile(previewPath)
	if err != nil {
		return
		//h.writeError(w, errors.Wrapf(err,
		//	"can't read preview file for project %q", p.Project.ID()),
		//	http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if _, err := w.Write(imgData); err != nil {
		log.Printf("error writing image data for preview of %q: %q", p.Project.ID(), err)
	}
}

