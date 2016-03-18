package main

import (
	"flag"
	"log"

	"github.com/qfarm/qfarm/worker"
)

var configPath = flag.String("config-path", "config/worker.toml", "Path to configuration file.")

func main() {
	flag.Parse()
	log.Printf("Creating new worker...")

	cfg, err := worker.Load(*configPath)
	if err != nil {
		log.Fatalf("Can't load config file: %v", err)
	}

	w, err := worker.NewWorker(cfg)
	if err != nil {
		log.Fatalf("Can't initialize worker: %v", err)
	}

	log.Printf("Worker created! Starting!")
	if err := w.Run(); err != nil {
		log.Fatalf("Can't initialize worker: %v", err)
	}
}
