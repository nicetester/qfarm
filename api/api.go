package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/qfarm/qfarm"
	"github.com/qfarm/qfarm/redis"
)

// Service is an API service with Redis connection.
type Service struct {
	r *redis.Service
}

// NewService creates new API service.
func NewService(r *redis.Service) *Service {
	return &Service{r: r}
}

// TriggerBuild adds build request to Redis pubsub.
func (s *Service) TriggerBuild(w http.ResponseWriter, req *http.Request) {
	dec := json.NewDecoder(req.Body)
	build := new(qfarm.Build)
	if err := dec.Decode(build); err != nil {
		writeErrJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := s.r.ListPush("test-q-list", build.Repo); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}
	if err := s.r.Publish("test-q-channel", build.Repo); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}
}

// LastBuilds returns most recent builds among all repositories.
func (s *Service) LastBuilds(w http.ResponseWriter, req *http.Request) {
	builds, err := s.r.ListGetLastElements("all-builds", 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastBuilds := make([]qfarm.Build, 0)
	for _, b := range builds {
		var single qfarm.Build
		if err := json.Unmarshal(b, &single); err != nil {
			writeErrJSON(w, err, http.StatusInternalServerError)
			return
		}

		lastBuilds = append(lastBuilds, single)
	}

	if err := writeJSON(w, lastBuilds); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}
}

// LastRepoBuilds returns most recent builds among specified repository.
func (s *Service) LastRepoBuilds(w http.ResponseWriter, req *http.Request) {
	repo := req.URL.Query().Get("repo")
	if repo == "" {
		writeErrJSON(w, errors.New("Repo should be set!"), http.StatusBadRequest)
		return
	}

	builds, err := s.r.ListGetLastElements("builds:"+repo, 10)
	if err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}

	lastBuilds := make([]qfarm.Build, 0)
	for _, b := range builds {
		var single qfarm.Build
		if err := json.Unmarshal(b, &single); err != nil {
			writeErrJSON(w, err, http.StatusInternalServerError)
			return
		}

		lastBuilds = append(lastBuilds, single)
	}

	if err := writeJSON(w, lastBuilds); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}
}

// writeErrJSON wraps error in JSON structure.
func writeErrJSON(w http.ResponseWriter, err error, status int) {
	log.Print(err.Error())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var errMap = map[string]interface{}{
		"error": err.Error(),
	}

	body, _ := json.Marshal(errMap)
	http.Error(w, string(body), status)
}

// writeJSON write response to client, response is a struct defining JSON reply.
func writeJSON(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json, err := json.Marshal(response)
	if err != nil {
		return err
	}

	if _, err := w.Write(json); err != nil {
		return err
	}

	return nil
}
