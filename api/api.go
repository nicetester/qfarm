package api

import (
	"encoding/json"
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
		WriteErrJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := s.r.ListPush("test-q-list", build.Repo); err != nil {
		WriteErrJSON(w, err, http.StatusInternalServerError)
		return
	}
	if err := s.r.Publish("test-q-channel", build.Repo); err != nil {
		WriteErrJSON(w, err, http.StatusInternalServerError)
		return
	}
}

// WriteErrJSON wraps error in JSON structure.
func WriteErrJSON(w http.ResponseWriter, err error, status int) {
	log.Print(err.Error())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var errMap = map[string]interface{}{
		"error": err.Error(),
	}

	body, _ := json.Marshal(errMap)
	http.Error(w, string(body), status)
}
