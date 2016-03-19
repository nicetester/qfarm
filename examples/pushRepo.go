package main

import (
	"flag"
	"fmt"
	"github.com/qfarm/qfarm/redis"
)

var srv *redis.Service
var (
	repo  = flag.String("repo", "github.com/qfarm/bad-go-code", "Repo to analysis")
	redis = flags.String("redis", "docker:6379", "Redis connection string")
)

func main() {
	flag.Parse()
	cfg := redis.NewConfig().WithConnection(redis).WithPassword("")
	var err error
	srv, err = redis.NewService(cfg)
	if err != nil {
		fmt.Printf("Can't create the redis service: %v\n", err)
		return
	}

	if err = pushRepo(*repo); err != nil {
		fmt.Printf("pushRepo test failed: %v\n", err)
		return
	}

	if err = publishRepo(*repo); err != nil {
		fmt.Printf("publishRepo test failed: %v\n", err)
		return
	}

	fmt.Printf("Worker notified about repo: %s\n", *repo)
}

func pushRepo(repo string) error {
	if err := srv.ListPush("test-q-list", repo); err != nil {
		return fmt.Errorf("Can't push value in redis: %v", err)
	}

	return nil
}

func publishRepo(repo string) error {
	if err := srv.Publish("test-q-channel", repo); err != nil {
		return fmt.Errorf("Can't push value in redis: %v", err)
	}

	return nil
}
