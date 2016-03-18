package main

import (
	"fmt"

	"github.com/qfarm/qfarm/redis"
)

var srv *redis.Service

func main() {
	cfg := redis.NewConfig().WithConnection("127.0.0.1:6379").WithPassword("")
	var err error
	srv, err = redis.NewService(cfg)
	if err != nil {
		fmt.Printf("Can't create the redis service: %v\n", err)
		return
	}

	if err = getSetTest(); err != nil {
		fmt.Printf("GetSet test failed: %v\n", err)
		return
	}

	if err = pushPopTest(); err != nil {
		fmt.Printf("PushPop test failed: %v\n", err)
		return
	}

	fmt.Printf("Test done with success\n")
}

func getSetTest() error {
	val := "test value goes here"
	if err := srv.Set("test", 100, val); err != nil {
		return fmt.Errorf("Can't set value in redis: %v", err)
	}

	bytes, err := srv.Get("test")
	if err != nil {
		return fmt.Errorf("Can't get the value from redis: %v", err)
	}

	if string(bytes) != val {
		return fmt.Errorf("Value from redis is different that expected! Expected: %s Got: %s", val, string(bytes))
	}

	return nil
}

func pushPopTest() error {
	val := "test value goes here"
	if err := srv.ListPush("test-list", val); err != nil {
		return fmt.Errorf("Can't push value in redis: %v", err)
	}

	elem, err := srv.ListPop("test-list")
	if err != nil {
		return fmt.Errorf("Can't pop the value from redis: %v", err)
	}

	if string(elem.([]byte)) != val {
		return fmt.Errorf("Value from redis is different that expected! Expected: %s Got: %s", val, string(elem.([]byte)))
	}

	return nil
}