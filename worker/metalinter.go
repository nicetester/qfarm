package worker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/shlex"
	"github.com/qfarm/qfarm"
	"log"
	"github.com/qfarm/qfarm/redis"
	"encoding/json"
)

const (
	Aligncheck  = "aligncheck"
	Deadcode    = "deadcode"
	Dupl        = "dupl"
	Errcheck    = "errcheck"
	Goconst     = "goconst"
	Gocyclo     = "gocyclo"
	Gofmt       = "gofmt"
	Goimports   = "goimports"
	Golint      = "golint"
	Gotype      = "gotype"
	Ineffassign = "ineffassign"
	Interfacer  = "interfacer"
	Lll         = "lll"
	Structcheck = "structcheck"
	Test        = "test"
	Testify     = "testify"
	Varcheck    = "varcheck"
	Vet         = "vet"
	Vetshadow   = "vetshadow"
	Unconvert   = "unconvert"
)

var (
	linters = map[string]string{
		Aligncheck:  `aligncheck .:^(?:[^:]+: )?(?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.+)$`, // slow
		Deadcode:    `deadcode .:^deadcode: (?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.*)$`,
		Dupl:        `dupl -plumbing -threshold {duplthreshold} ./*.go:^(?P<path>[^\s][^:]+?\.go):(?P<line>\d+)-\d+:\s*(?P<message>.*)$`,
		Errcheck:    `errcheck -abspath .:^(?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+)\t(?P<message>.*)$`, // slow
		Goconst:     `goconst -min-occurrences {min_occurrences} .:PATH:LINE:COL:MESSAGE`,
		Gocyclo:     `gocyclo -over {mincyclo} .:^(?P<cyclo>\d+)\s+\S+\s(?P<function>\S+)\s+(?P<path>[^:]+):(?P<line>\d+):(\d+)$`,
		Gofmt:       `gofmt -l -s ./*.go:^(?P<path>[^\n]+)$`,
		Goimports:   `goimports -l ./*.go:^(?P<path>[^\n]+)$`,
		Golint:      "golint -min_confidence {min_confidence} .:PATH:LINE:COL:MESSAGE",
		Gotype:      "gotype -e {tests=-a} .:PATH:LINE:COL:MESSAGE",
		Ineffassign: `ineffassign -n .:PATH:LINE:COL:MESSAGE`,
		Interfacer:  `interfacer ./:PATH:LINE:COL:MESSAGE`,
		Lll:         `lll -g -l {maxlinelength} ./*.go:PATH:LINE:MESSAGE`,
		Structcheck: `structcheck {tests=-t} .:^(?:[^:]+: )?(?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.+)$`, // slow
		Test:        `go test:^--- FAIL: .*$\s+(?P<path>[^:]+):(?P<line>\d+): (?P<message>.*)$`,                             // slow
		Testify:     `go test:Location:\s+(?P<path>[^:]+):(?P<line>\d+)$\s+Error:\s+(?P<message>[^\n]+)`,                    // slow
		Varcheck:    `varcheck .:^(?:[^:]+: )?(?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>\w+)$`,              // slow
		Vet:         "go tool vet ./*.go:PATH:LINE:MESSAGE",
		Vetshadow:   "go tool vet --shadow ./*.go:PATH:LINE:MESSAGE",
		Unconvert:   "unconvert .:PATH:LINE:COL:MESSAGE", // slow
	}
	linterMessageOverrideFlag = map[string]string{
		Errcheck:    "error return value not checked ({message})",
		Varcheck:    "unused global variable {message}",
		Structcheck: "unused struct field {message}",
		Gocyclo:     "cyclomatic complexity {cyclo} of function {function}() is high (> {mincyclo})",
		Gofmt:       "file is not gofmted",
		Goimports:   "file is not goimported",
		Unconvert:   "redundant type conversion",
	}
	linterSeverityFlag = map[string]string{
		Gotype:  "error",
		Test:    "error",
		Testify: "error",
		Vet:     "error",
	}
	predefinedPatterns = map[string]string{
		"PATH:LINE:COL:MESSAGE": `^(?P<path>[^\s][^\r\n:]+?\.go):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.*)$`,
		"PATH:LINE:MESSAGE":     `^(?P<path>[^\s][^\r\n:]+?\.go):(?P<line>\d+):\s*(?P<message>.*)$`,
	}
	installMap = map[string]string{
		Golint:      "github.com/golang/lint/golint",
		Gotype:      "golang.org/x/tools/cmd/gotype",
		Goimports:   "golang.org/x/tools/cmd/goimports",
		Errcheck:    "github.com/kisielk/errcheck",
		Varcheck:    "github.com/opennota/check/cmd/varcheck",
		Structcheck: "github.com/opennota/check/cmd/structcheck",
		Aligncheck:  "github.com/opennota/check/cmd/aligncheck",
		Deadcode:    "github.com/tsenart/deadcode",
		Gocyclo:     "github.com/alecthomas/gocyclo",
		Ineffassign: "github.com/gordonklaus/ineffassign",
		Dupl:        "github.com/mibk/dupl",
		Interfacer:  "github.com/mvdan/interfacer/cmd/interfacer",
		Lll:         "github.com/walle/lll/cmd/lll",
		Unconvert:   "github.com/mdempsky/unconvert",
		Goconst:     "github.com/jgautheron/goconst/cmd/goconst",
	}
	defaultLinters = []string{
		Deadcode,
		Dupl,
		Goconst,
		Gocyclo,
		Gofmt,
		Goimports,
		Golint,
		Gotype,
		Ineffassign,
		Interfacer,
		Lll,
		Vet,
		Vetshadow,
	}
)

type Metalinter struct {
	cfg      *Cfg
	notifier *Notifier
	redis *redis.Service
}

func NewMetalinter(cfg *Cfg, redis *redis.Service, notifier *Notifier) *Metalinter {
	return &Metalinter{cfg: cfg, redis: redis, notifier: notifier}
}

type Severity string

// Linter message severity levels.
const (
	Warning Severity = "warning"
	Error   Severity = "error"
)

func LinterFromName(name string) (*qfarm.Linter, error) {
	s := linters[name]
	parts := strings.SplitN(s, ":", 2)
	pattern := parts[1]
	if p, ok := predefinedPatterns[pattern]; ok {
		pattern = p
	}
	re, err := regexp.Compile("(?m:" + pattern + ")")
	if err != nil {
		return nil, err
	}

	return &qfarm.Linter{
		Name:             name,
		Command:          s[0:strings.Index(s, ":")],
		Pattern:          pattern,
		InstallFrom:      installMap[name],
		SeverityOverride: qfarm.Severity(linterSeverityFlag[name]),
		MessageOverride:  linterMessageOverrideFlag[name],
		Regex:            re,
		EventType:        linterEventsMapping[name],
	}, nil
}

func (m *Metalinter) Start(cfg qfarm.BuildCfg, buildNo int, ft *FilesMap) error {
	start := time.Now()
	paths := m.expandPaths([]string{cfg.Path + "/..."}, cfg.SkipDirs)

	m.debug("Analyzing following paths: %v", paths)

	linters := m.linters(cfg.Linters)
	issues, errch := m.runLinters(linters, cfg.Repo, paths, m.cfg.Concurrency, cfg.IncludeTests)

	for issue := range issues {
		if strings.HasSuffix(issue.Path, ".gen.go") || strings.HasSuffix(issue.Path, ".pb.go") {
			continue
		}

		// apply issue to all parents in file tree
		if err := ft.ApplyIssue(issue); err != nil {
			return err
		}

		// trim path in json
		issue.Path = strings.Replace(issue.Path, cfg.Path, "", -1)

		// marshal issue to json
		data, err := json.Marshal(issue)
		if err != nil {
			return err
		}

		// store issue in global list of issues
		_, err = m.redis.SortedSetAdd(fmt.Sprintf("issues:%s:%d", cfg.Repo, buildNo), data, issue.Severity.Rank())
		if err != nil {
			return err
		}
	}

	for err := range errch {
		warning("%s", err)
	}

	elapsed := time.Now().Sub(start)
	m.debug("total elapsed time %s", elapsed)

	return nil
}

func (m *Metalinter) debug(format string, args ...interface{}) {
	if m.cfg.Debug {
		log.Printf("DEBUG: "+format+"\n", args...)
	}
}

func warning(format string, args ...interface{}) {
	log.Printf("WARNING: "+format, args...)
}

type Vars map[string]string

func (v Vars) Copy() Vars {
	out := Vars{}
	for k, v := range v {
		out[k] = v
	}
	return out
}

func (v Vars) Replace(s string) string {
	for k, v := range v {
		prefix := regexp.MustCompile(fmt.Sprintf("{%s=([^}]*)}", k))
		if v != "" {
			s = prefix.ReplaceAllString(s, "$1")
		} else {
			s = prefix.ReplaceAllString(s, "")
		}
		s = strings.Replace(s, fmt.Sprintf("{%s}", k), v, -1)
	}
	return s
}

func (m *Metalinter) runLinters(linters map[string]*qfarm.Linter, repo string, paths []string, concurrency int, includeTests bool) (chan *qfarm.Issue, chan error) {
	errch := make(chan error, len(linters)*len(paths))
	concurrencych := make(chan bool, concurrency)
	incomingIssues := make(chan *qfarm.Issue, 1000000)
	wg := &sync.WaitGroup{}
	for _, linter := range linters {
		// Recreated in each loop because it is mutated by executeLinter().
		vars := Vars{
			"duplthreshold":   fmt.Sprintf("%d", m.cfg.DuplThreshold),
			"mincyclo":        fmt.Sprintf("%d", m.cfg.Cyclo),
			"maxlinelength":   fmt.Sprintf("%d", m.cfg.LLLineLength),
			"min_confidence":  fmt.Sprintf("%f", m.cfg.GolintMinConfidence),
			"min_occurrences": fmt.Sprintf("%d", m.cfg.GoconstMinOccurrences),
			"tests":           "",
		}
		if includeTests {
			vars["tests"] = "-t"
		}

		// wait group for single linter.
		wgl := sync.WaitGroup{}
		for _, path := range paths {
			wg.Add(1)
			wgl.Add(1)
			state := &linterState{
				Linter: linter,
				issues: incomingIssues,
				path:   path,
				vars:   vars.Copy(),
				repo:   repo,
			}
			go func() {
				concurrencych <- true
				err := m.executeLinter(state)
				if err != nil {
					errch <- err
				}
				<-concurrencych
				wg.Done()
				wgl.Done()
			}()
		}
		go func(repo, linterName, eventType string) {
			wgl.Wait()
			m.notifier.SendEvent(repo, fmt.Sprintf("Linter %s finished!", linterName), eventType)
		}(repo, linter.Name, linter.EventType)
	}

	go func() {
		wg.Wait()
		close(incomingIssues)
		close(errch)
	}()
	return incomingIssues, errch
}

func (m *Metalinter) expandPaths(paths, skip []string) []string {
	if len(paths) == 0 {
		paths = []string{"."}
	}
	skipMap := map[string]bool{}
	for _, name := range skip {
		skipMap[name] = true
	}
	dirs := map[string]bool{}
	for _, path := range paths {
		if strings.HasSuffix(path, "/...") {
			root := filepath.Dir(path)
			_ = filepath.Walk(root, func(p string, i os.FileInfo, err error) error {
				if err != nil {
					warning("invalid path %q: %s", p, err)
					return err
				}

				base := filepath.Base(p)
				skip := skipMap[base] || skipMap[p] || (strings.ContainsAny(base[0:1], "_.") && base != "." && base != "..")
				if i.IsDir() {
					if skip {
						return filepath.SkipDir
					}
				} else if !skip && !strings.HasSuffix(p, ".gen.go") && !strings.HasSuffix(p, ".pb.go") && strings.HasSuffix(p, ".go") {
					dirs[filepath.Clean(filepath.Dir(p))] = true
				}
				return nil
			})
		} else {
			dirs[filepath.Clean(path)] = true
		}
	}
	out := make([]string, 0, len(dirs))
	for d := range dirs {
		out = append(out, d)
	}
	sort.Strings(out)
	for _, d := range out {
		m.debug("linting path %s", d)
	}
	return out
}

func (m *Metalinter) makeInstallCommand(linters ...string) []string {
	cmd := []string{"get"}
	cmd = append(cmd, "-v")
	cmd = append(cmd, "-u")
	cmd = append(cmd, "-f")
	cmd = append(cmd, linters...)
	return cmd
}

func (m *Metalinter) InstallAllLinters() error {
	var linters []string
	for _, v := range installMap {
		linters = append(linters, v)
	}

	cmd := m.makeInstallCommand(linters...)
	c := exec.Command("go", cmd...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

type linterState struct {
	*qfarm.Linter
	path   string
	issues chan *qfarm.Issue
	vars   Vars
	repo   string
}

func (l *linterState) InterpolatedCommand() string {
	l.vars["path"] = l.path
	return l.vars.Replace(l.Command)
}

func parseCommand(dir, command string) (string, []string, error) {
	args, err := shlex.Split(command)
	if err != nil {
		return "", nil, err
	}
	if len(args) == 0 {
		return "", nil, fmt.Errorf("invalid command %q", command)
	}
	exe, err := exec.LookPath(args[0])
	if err != nil {
		return "", nil, err
	}
	out := []string{}
	for _, arg := range args[1:] {
		if strings.Contains(arg, "*") {
			pattern := filepath.Join(dir, arg)
			globbed, err := filepath.Glob(pattern)
			if err != nil {
				return "", nil, err
			}
			for i, g := range globbed {
				if strings.HasPrefix(g, dir+"/") {
					globbed[i] = g[len(dir)+1:]
				}
			}
			out = append(out, globbed...)
		} else {
			out = append(out, arg)
		}
	}
	return exe, out, nil
}

func (m *Metalinter) executeLinter(state *linterState) error {
	m.debug("linting with %s: %s (on %s)", state.Name, state.Command, state.path)

	start := time.Now()
	command := state.InterpolatedCommand()
	exe, args, err := parseCommand(state.path, command)
	if err != nil {
		return err
	}
	m.debug("executing %s %q", exe, args)
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command(exe, args...)
	cmd.Dir = state.path
	cmd.Stdout = buf
	cmd.Stderr = buf
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to execute linter %s: %s", command, err)
	}

	done := make(chan bool)
	go func() {
		err = cmd.Wait()
		done <- true
	}()

	// Wait for process to complete or deadline to expire.
	<-done

	if err != nil {
		m.debug("warning: %s returned %s", command, err)
	}

	err = m.processOutput(state, buf.Bytes())
	if err != nil {
		return err
	}

	elapsed := time.Now().Sub(start)
	m.debug("%s linter took %s", state.Name, elapsed)

	return nil
}

func (l *linterState) fixPath(path string) string {
	abspath, err := filepath.Abs(l.path)
	if filepath.IsAbs(path) {
		if err == nil && strings.HasPrefix(path, abspath) {
			normalised := filepath.Join(abspath, filepath.Base(path))
			if _, err := os.Stat(normalised); err == nil {
				path := filepath.Join(l.path, filepath.Base(path))
				return path
			}
		}
	} else {
		return filepath.Join(l.path, path)
	}
	return path
}

func (m *Metalinter) linters(cfgLinters []string) map[string]*qfarm.Linter {
	out := map[string]*qfarm.Linter{}
	for _, name := range cfgLinters {
		_, ok := linters[name]
		if ok {
			out[name], _ = LinterFromName(name)
		} else {
			warning("Linter %s doesn't exist!", name)
		}
	}

	return out
}

func (m *Metalinter) processOutput(state *linterState, out []byte) error {
	re := state.Regex
	all := re.FindAllSubmatchIndex(out, -1)
	m.debug("%s hits %d: %s", state.Name, len(all), state.Pattern)
	for _, indices := range all {
		group := [][]byte{}
		for i := 0; i < len(indices); i += 2 {
			fragment := out[indices[i]:indices[i+1]]
			group = append(group, fragment)
		}

		issue := &qfarm.Issue{Line: 1}
		issue.Linter, _ = LinterFromName(state.Name)
		for i, name := range re.SubexpNames() {
			part := string(group[i])
			if name != "" {
				state.vars[name] = part
			}
			switch name {
			case "path":
				issue.Path = state.fixPath(part)

			case "line":
				n, err := strconv.ParseInt(part, 10, 32)
				if err != nil {
					return fmt.Errorf("line matched invalid integer: %v", err)
				}
				issue.Line = int(n)

			case "col":
				n, err := strconv.ParseInt(part, 10, 32)
				if err != nil {
					return fmt.Errorf("col matched invalid integer: %v", err)
				}
				issue.Col = int(n)
			case "message":
				issue.Message = part

			case "":
			}
		}
		if m, ok := linterMessageOverrideFlag[state.Name]; ok {
			issue.Message = state.vars.Replace(m)
		}
		if sev, ok := linterSeverityFlag[state.Name]; ok {
			issue.Severity = qfarm.Severity(sev)
		} else {
			issue.Severity = "warning"
		}
		state.issues <- issue
	}
	return nil
}
