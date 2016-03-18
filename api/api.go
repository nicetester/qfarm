package api

import (
	"encoding/json"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.r.ListPush("test-q-list", build.Repo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := s.r.Publish("test-q-channel", build.Repo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
