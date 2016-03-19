package worker

import (
	"fmt"

	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/qfarm/qfarm"
	"github.com/qfarm/qfarm/redis"
)

type Worker struct {
	linter   *Metalinter
	redis    *redis.Service
	notifier *Notifier
	coverage *CoverageChecker
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

	w.notifier = NewNotifier(w.redis)
	w.linter = NewMetalinter(config, w.redis, w.notifier)
	w.coverage = NewCoverageChecker(config, w.notifier)

	return w, nil
}

func (w *Worker) Run() error {
	if err := w.redis.Subscribe("test-q-channel", w.fetchAndAnalyze); err != nil {
		return err
	}

	return nil
}

func (w *Worker) fetchAndAnalyze(data interface{}) error {
	elem, err := w.redis.ListPop("test-q-list") // TODO: drain list to the bottom
	if err != nil {
		// do nothing other worker might got the value from list before
		return nil
	}

	if err := w.analyze(string(elem.([]byte))); err != nil {
		log.Printf("Error during worker analysis! Err: %v \n", err)
	}

	return nil
}

func (w *Worker) analyze(repo string) error {
	// download repo
	if err := w.download(repo); err != nil {
		return err
	}

	if err := w.markAsUserRepo(repo); err != nil {
		return err
	}

	lastCommitHash, err := lastCommitHash(repo)
	if err != nil {
		return err
	}

	log.Printf("Hash of last commit %s", lastCommitHash)

	// get last build number
	firstTimeBuild := false
	buildInfo, err := w.getLastBuildInfo(repo)
	if err != nil {
		if err == redis.ErrNotFound {
			firstTimeBuild = true
		} else {
			return err
		}
	}

	if !firstTimeBuild && w.config.CheckLastCommitHash {
		// someone wants to analyze the same repo twice
		if buildInfo.CommitHash == lastCommitHash {
			w.notifier.SendEventWithPayload(repo, fmt.Sprintf("Repo %s already analyzed!", repo), EventTypeAlreadyAnalyzed, fmt.Sprintf("%s", buildInfo.No))
			return fmt.Errorf("repo %s already analyzed!", repo)
		}
	}

	// generate new build no
	newBuild := qfarm.Build{Repo: repo, CommitHash: lastCommitHash, Time: time.Now().UTC()}
	if firstTimeBuild {
		newBuild.No = 1
	} else {
		newBuild.No = buildInfo.No + 1
	}

	// create repo config
	buildCfg, err := LoadRepoCfg(repo, path.Join(os.Getenv("GOPATH"), "src", repo))
	if err != nil {
		return err
	}

	// marshal build info
	newBuild.Config = *buildCfg
	data, err := json.Marshal(newBuild)
	if err != nil {
		return err
	}

	// add new build to global list of all builds
	if err := w.redis.ListPush("all-builds", data); err != nil {
		return err
	}

	// add new build to list of builds per repo
	if err := w.redis.ListPush("builds:"+repo, data); err != nil {
		return err
	}

	// generate directory structure
	ft, err := BuildTree(buildCfg.Path)
	if err != nil {
		return err
	}

	// run all linters
	if err := w.linter.Start(*buildCfg, newBuild.No, ft); err != nil {
		return err
	}

	// run coverage
	if err := w.coverage.Start(*buildCfg, ft); err != nil {
		return err
	}

	// generate report

	w.notifier.SendEvent(repo, "All tasks done!", EventTypeAllDone)

	fmt.Printf("All done\n")
	return nil
}

func (w *Worker) getLastBuildInfo(repo string) (qfarm.Build, error) {
	var build qfarm.Build
	data, err := w.redis.ListGetLast("builds:" + repo)
	if err != nil {
		return build, err
	}

	if err := json.Unmarshal(data.([]byte), &build); err != nil {
		return build, err
	}

	return build, nil
}

func (w *Worker) download(repo string) error {
	fmt.Printf("Downloading %s...\n", repo)
	if err := exec.Command("go", "get", "-u", "-t", path.Join(repo, "...")).Run(); err != nil {
		return err
	}

	fmt.Printf("Repo %s downloaded!\n", repo)

	w.notifier.SendEvent(repo, fmt.Sprintf("Repo %s downloaded", repo), EventTypeDownloadDone)

	return nil
}

func (w *Worker) markAsUserRepo(repo string) error {
	userName := strings.Split(repo, "/")[1]
	_, err := w.redis.SortedSetRank("users:"+userName+":repos", repo)

	return err
}

func lastCommitHash(repo string) (string, error) {
	repoPath := path.Join(os.Getenv("GOPATH"), "src", repo)

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
