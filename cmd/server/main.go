package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/qfarm/qfarm/api"
	"github.com/qfarm/qfarm/redis"
)

var listen = flag.String("listen", ":8080", "HTTP listen on")
var redisConn = flag.String("redis-conn", "redis:6379", "Redis connection string")

func main() {
	flag.Parse()

	cfg := redis.NewConfig().WithConnection(*redisConn)
	r, err := redis.NewService(cfg)
	if err != nil {
		log.Fatalf("Can't create redis service: %v", err)
	}

	as := api.NewService(r)
	router := mux.NewRouter()
	router.HandleFunc("/build/", as.TriggerBuild).Methods("POST")
	router.HandleFunc("/last_builds/", as.LastBuilds).Methods("GET")
	router.HandleFunc("/last_repo_builds/", as.LastRepoBuilds).Methods("GET")
	router.HandleFunc("/user_repos/", as.UserRepos).Methods("GET")
	router.HandleFunc("/issues/", as.RepoIssues).Methods("GET")
	router.HandleFunc("/files/", as.RepoFiles).Methods("GET")
	router.HandleFunc("/reports/", as.Report).Methods("GET")

	http.Handle("/", handlers.CORS()(router))
	log.Printf("Starting to serve on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
