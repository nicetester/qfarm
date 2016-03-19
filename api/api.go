package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"fmt"
	"strconv"

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

	repo := strings.TrimRight(build.Repo, "/")

	if err := s.r.ListPush("test-q-list", repo); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}
	if err := s.r.Publish("test-q-channel", repo); err != nil {
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

	for i, j := 0, len(lastBuilds)-1; i < j; i, j = i+1, j-1 {
		lastBuilds[i], lastBuilds[j] = lastBuilds[j], lastBuilds[i]
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

	for i, j := 0, len(lastBuilds)-1; i < j; i, j = i+1, j-1 {
		lastBuilds[i], lastBuilds[j] = lastBuilds[j], lastBuilds[i]
	}

	if err := writeJSON(w, lastBuilds); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}
}

func (s *Service) UserRepos(w http.ResponseWriter, req *http.Request) {
	user := req.URL.Query().Get("user")
	if user == "" {
		writeErrJSON(w, errors.New("User should be set!"), http.StatusBadRequest)
		return
	}

	repos, err := s.r.SortedSetGetAllRev("users:" + user + ":repos")
	if err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}

	userRepos := make([]string, 0)
	for _, r := range repos {
		userRepos = append(userRepos, string(r))
	}

	if err := writeJSON(w, userRepos); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}
}

func (s *Service) Report(w http.ResponseWriter, req *http.Request) {
	repo := req.URL.Query().Get("repo")
	if repo == "" {
		writeErrJSON(w, errors.New("Repo should be set!"), http.StatusBadRequest)
		return
	}

	var err error
	buildNoInt := 0
	buildNo := req.URL.Query().Get("no")
	if buildNo == "" {
		buildNoInt, err = s.getLastBuildNo(repo)
		if err != nil {
			writeErrJSON(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		buildNoInt, err = strconv.Atoi(buildNo)
		if err != nil {
			writeErrJSON(w, err, http.StatusBadRequest)
			return
		}
	}

	reportJson, err := s.r.Get(fmt.Sprintf("reports:%s:%d", repo, buildNoInt))
	if err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(reportJson); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
	}
}

// RepoIssues returns list of specified repo issues.
func (s *Service) RepoIssues(w http.ResponseWriter, req *http.Request) {
	repo := req.URL.Query().Get("repo")
	if repo == "" {
		writeErrJSON(w, errors.New("Repo should be set!"), http.StatusBadRequest)
		return
	}

	var err error
	buildNoInt, sizeInt, skipInt := 0, 50, 0
	buildNo := req.URL.Query().Get("no")
	if buildNo == "" {
		buildNoInt, err = s.getLastBuildNo(repo)
		if err != nil {
			writeErrJSON(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		buildNoInt, err = strconv.Atoi(buildNo)
		if err != nil {
			writeErrJSON(w, err, http.StatusBadRequest)
			return
		}
	}

	size := req.URL.Query().Get("size")
	if size == "" {

	} else {
		sizeInt, err = strconv.Atoi(size)
		if err != nil {
			writeErrJSON(w, err, http.StatusBadRequest)
			return
		}
	}

	skip := req.URL.Query().Get("skip")
	if skip != "" {
		skipInt, err = strconv.Atoi(skip)
		if err != nil {
			writeErrJSON(w, err, http.StatusBadRequest)
			return
		}
	}

	data := make([][]byte, 0)
	filter := req.URL.Query().Get("filter")
	if filter != "" {
		data, err = s.r.SortedSetGet(fmt.Sprintf("issues:%s:%d:%s", repo, buildNoInt, filter), sizeInt, skipInt)
		if err != nil {
			writeErrJSON(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		data, err = s.r.SortedSetGet(fmt.Sprintf("issues:%s:%d", repo, buildNoInt), sizeInt, skipInt)
		if err != nil {
			writeErrJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	issues := make([]qfarm.Issue, 0)
	for _, b := range data {
		var single qfarm.Issue
		single.Linter = new(qfarm.Linter)
		if err := json.Unmarshal(b, &single); err != nil {
			writeErrJSON(w, err, http.StatusInternalServerError)
			return
		}

		issues = append(issues, single)
	}

	if err := writeJSON(w, issues); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}
}

// RepoFiles returns list of specified repo issues.
func (s *Service) RepoFiles(w http.ResponseWriter, req *http.Request) {
	repo := req.URL.Query().Get("repo")
	if repo == "" {
		writeErrJSON(w, errors.New("Repo should be set!"), http.StatusBadRequest)
		return
	}

	var err error
	var buildNoInt int
	buildNo := req.URL.Query().Get("no")
	if buildNo == "" {
		buildNoInt, err = s.getLastBuildNo(repo)
		if err != nil {
			writeErrJSON(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		buildNoInt, err = strconv.Atoi(buildNo)
		if err != nil {
			writeErrJSON(w, err, http.StatusBadRequest)
			return
		}
	}

	keys, err := s.r.Keys(fmt.Sprintf("files:%s:%d:*", repo, buildNoInt))
	if err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}

	nodes := make([]qfarm.Node, 0)
	for _, k := range keys {
		data, err := s.r.Get(k)
		if err != nil {
			writeErrJSON(w, err, http.StatusInternalServerError)
			return
		}
		var node qfarm.Node
		if err := json.Unmarshal(data, &node); err != nil {
			writeErrJSON(w, err, http.StatusInternalServerError)
			return
		}

		nodes = append(nodes, node)
	}

	if err := writeJSON(w, nodes); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}
}

func (s *Service) Badge(w http.ResponseWriter, req *http.Request) {
	repo := req.URL.Query().Get("repo")
	if repo == "" {
		writeErrJSON(w, errors.New("Repo should be set!"), http.StatusBadRequest)
		return
	}

	// get last report for spefied repo
	buildNoInt, err := s.getLastBuildNo(repo)
	if err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}

	reportJson, err := s.r.Get(fmt.Sprintf("reports:%s:%d", repo, buildNoInt))
	if err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}

	var r qfarm.Report
	if err = json.Unmarshal(reportJson, &r); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
		return
	}

	color := ""
	if r.Score > 80 {
		color = `#4CAF50`
	} else if r.Score > 60 {
		color = `#FFC107`
	} else {
		color = `#F44336`
	}

	badge := fmt.Sprintf(`
		<svg xmlns="http://www.w3.org/2000/svg" width="105" height="20">
		 <rect fill="#555" height="20" width="75"/>
		 <rect fill="%s" x="75" height="20" width="30"/>
		 <path fill-opacity=".1" d="M0 0h179v20h-179z"/>
		 <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
		  <text x="40" y="14">Quality</text>
		  <text x="90" y="14">%d</text>
		 </g>
		</svg>`, color, r.Score)

	w.Header().Set("Content-Type", "image/svg+xml")
	if _, err := w.Write([]byte(badge)); err != nil {
		writeErrJSON(w, err, http.StatusInternalServerError)
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

func (s *Service) getLastBuildNo(repo string) (int, error) {
	var build qfarm.Build
	data, err := s.r.ListGetLast("builds:" + repo)
	if err != nil {
		return -1, err
	}

	if err := json.Unmarshal(data.([]byte), &build); err != nil {
		return -1, err
	}

	return build.No, nil
}
