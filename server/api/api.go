package api

import (
	"golang.org/x/net/context"

	"net/http"
	"bitbucket.org/yasnikoff/shuffler/app"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
)

//type handlerFunc func(s *handler, w http.ResponseWriter, r *http.Request)

type handler struct {
	repos   map[string]app.App
	uiPath  string
	Root    *mux.Router
	apiRoot *mux.Router
	Context context.Context
}

func NewHandler() *handler {
	router:=mux.NewRouter()
	handler := &handler{
		repos:make(map[string]app.App, 1),
		Root:router,
		uiPath:`./ui/web`,
		apiRoot:router.PathPrefix("/api/").Subrouter(),
		Context:context.Background(),
	}
	// ui
	handler.Root.PathPrefix("/").Handler(http.FileServer(http.Dir(handler.uiPath)))
	// api
	handler.apiRoot.HandleFunc("/", handler.handleRoot)
	// repo
	repo := handler.apiRoot.PathPrefix("/{repo}/").Subrouter()
	repo.HandleFunc("/", handler.handleRepo)
	// repo.HandleFunc("/query/", handler.query)
	repo.HandleFunc("/index/", handler.handleEmpty)
	repo.HandleFunc("/index/update/", handler.handleRepoIndexUpdate)
	// projects
	projects:=repo.PathPrefix("/projects/").Subrouter()
	projects.Methods("GET").Subrouter().HandleFunc("/", handler.query)
	projects.Methods("POST").Subrouter().HandleFunc("/", handler.handleRepoNewProject)
	// project
	project := projects.PathPrefix("/{id}/").Subrouter()
	// GET
	projectGet:=project.Methods("GET").Subrouter()
	projectGet.HandleFunc("/", handler.handleProject)
	projectGet.HandleFunc("/preview/", handler.handlePreview)
	projectGet.HandleFunc("/open/", handler.handleProjectOpen)
	// POST
	projectPost:=project.Methods("POST").Subrouter()
	projectPost.HandleFunc("/", handler.handleProjectSave)
	projectPost.HandleFunc("/duplicate/", handler.handleProjectDuplicate)
	// path
	path := project.PathPrefix("/path/").Subrouter()
	path.HandleFunc("/", handler.handleProjectPath)
	path.HandleFunc("/save/", handler.handleProjectPathSave)
	// tags
	tags := project.PathPrefix("/tags/").Subrouter()
	tags.HandleFunc("/", handler.handleTags)
	tags.HandleFunc("/add/", handler.handleTagsAdd)
	tags.HandleFunc("/remove/", handler.handleTagsRemove)
	// deps
	deps := project.PathPrefix("/dependencies/").Subrouter()
	deps.HandleFunc("/", handler.handleDeps)
	deps.HandleFunc("/add/", handler.handleDepsAdd)
	deps.HandleFunc("/remove/", handler.handleDepRemove)
	deps.HandleFunc("/update/", handler.handleDepUpdate)
	deps.HandleFunc("/replace/", handler.handleDepReplace)

	return handler
}

func (h *handler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Root.ServeHTTP(w, r)
}

// registerApp makes the application available under specified name.
// The name will be used as the first input of the URL path.
func (s *handler)registerApp(a app.App, name string) error {
	if _, ok := s.repos[name]; ok {
		return fmt.Errorf("app with the name %q is already registerred.", name)
	}
	//s.appLock.Lock()
	s.repos[name] = a
	//s.appLock.Unlock()
	return nil
}

// ListenAndServe serves given repositories on the given address.
// "repos" map is a mapping from repository names to their absolute paths.
// Repository names will be used as the first input in url paths.
func ListenAndServe(addr string, repos map[string]string) error {

	h := NewHandler()
	for name, repoPath := range repos {
		a, err := app.Bootstrap(repoPath)
		if err != nil {
			return err
		}
		err = h.registerApp(a, name)
		if err != nil {
			return err
		}
	}

	return http.ListenAndServe(addr, h)
}

func (h *handler)write(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if resp != nil {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error sending json of %#v: %s", resp, err)
		}
	}
}

func (h *handler)writeError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Printf("Error: %s", err)
	http.Error(w, err.Error(), status)
}

func (h *handler)readBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		err = errors.Wrap(err, "can't read request body")
		h.writeError(w, err, http.StatusBadRequest)
		return nil, err
	}
	return data, nil
}

func (h *handler)readJsonData(w http.ResponseWriter, r *http.Request, target interface{}) error {
	var m map[string]json.RawMessage

	data, err := h.readBody(w, r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &m); err != nil {
		err := errors.Wrap(err, "can't unmarshal json body")
		h.writeError(w, err, http.StatusBadRequest)
		return err
	}

	msg, ok := m["data"]
	if !ok {
		err := errors.New("no \"data\" in json body")
		h.writeError(w, errors.New("no \"data\" in json body"), http.StatusBadRequest)
		return err
	}

	if err := json.Unmarshal(msg, target); err != nil {
		err = errors.Wrapf(err, "can't unmarshal json data %q into %T", msg, target)
	}

	return nil
}

func (h *handler)handleEmpty(w http.ResponseWriter, r *http.Request) {
	h.write(w, nil)
}
