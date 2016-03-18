package worker

import (
	"fmt"

	"github.com/qfarm/qfarm/redis"
)

type Worker struct {
	redis    *redis.Service
	config   *Cfg
}

func NewWorker(config *Cfg) (*Worker, error) {
	w := &Worker{config: config}
	cfg := redis.NewConfig().WithConnection(config.RedisConn).WithPassword(config.RedisPass)
	var err error
	w.redis, err = redis.NewService(cfg)
	if err != nil {
		return nil, fmt.Errorf("Can't create the redis service: %v\n", err)
	}

	return w, nil
}

func (w *Worker) Run() error {
	// Run Redis pub-sub listener here
	return nil
}
