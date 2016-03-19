package main

import (
	"github.com/qfarm/qfarm/redis"
	"fmt"
	"flag"
)

var srv *redis.Service
var repo = flag.String("repo", "github.com/qfarm/bad-go-code/cover", "Repo to analysis")

func main() {
	flag.Parse()
	cfg := redis.NewConfig().WithConnection("127.0.0.1:6379").WithPassword("")
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