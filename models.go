package qfarm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Build represents single build.
type Build struct {
	User       string    `json:"user,omitempty"`
	Repo       string    `json:"repo,omitempty"`
	No         int       `json:"no"`
	Score      int       `json:"score"`
	Time       time.Time `json:"time,omitempty"`
	CommitHash string    `json:"commitHash,omitempty"`
	Config     BuildCfg  `json:"config,omitempty"`
}

// BuildCfg represents configuration of the build.
type BuildCfg struct {
	// Repo identifier eg. github.com/influxdata/influxdb
	Repo string `json:"repo"`

	// Project path
	Path string `json:"path"`

	// Skip directories
	SkipDirs []string `json:"skipDirs"`

	// Linters which should be used in analysis
	Linters []string `json:"linters"`

	// Enable vendoring support (skips 'vendor' directories and sets GO15VENDOREXPERIMENT=1).
	Vendor bool `json:"vendor"`

	// Go version
	Go string `json:"go"`

	// Include test files
	IncludeTests bool `json:"includeTests"`
}

// CoverageReport holds info about coverage analysis of entire repo.
type CoverageReport struct {
	Repo          string          `json:"repo"`
	TotalCoverage float64         `json:"totalCoverage"`
	TotalTestsNo  int             `json:"totalTestsNo"`
	TotalPassedNo int             `json:"totalPassedNo"`
	TotalFailedNo int             `json:"totalFailedNo"`
	TotalTime     time.Duration   `json:"totalTime"`
	Failed        bool            `json:"failed"`
	Packages      []PackageReport `json:"packages"`
}

// PackageReport holds info about coverage analysis of specified package.
type PackageReport struct {
	Name     string        `json:"name"`
	Coverage float64       `json:"coverage"`
	Failed   bool          `json:"failed"`
	TestsNo  int           `json:"testsNo"`
	PassedNo int           `json:"passedNo"`
	FailedNo int           `json:"failedNo"`
	Time     time.Duration `json:"time"`
	Files    map[string]CoverFileReport
}

// CoverFileReport holds coverage report for single file.
type CoverFileReport struct {
	Coverage float64
	Blocks   []CoverBlock
}

// CoverBlock holds part of coverage profile for single block in file.
type CoverBlock struct {
	Start   Cursor `json:"start"`
	End     Cursor `json:"end"`
	NumStmt int    `json:"numStmt"`
	Count   int    `json:"count"`
}

// Cursor points at specific location in file and is used by CoverBlock.
type Cursor struct {
	Line int `json:"line"`
	Col  int `json:"col"`
}

// Node represents a node in directory tree. It might be a file or a directory.
type Node struct {
	Path       string   `json:"path"`
	Nodes      []Node  `json:"nodes"`
	ParentPath string    `json:"parent"`
	Dir        bool     `json:"dir"`
	Coverage   float64  `json:"coverage"`
	TestsNo    int      `json:"testsNo"`
	FailedNo   int      `json:"failedNo"`
	PassedNo   int      `json:"passedNo"`
	IssuesNo   int      `json:"issuesNo"`
	ErrorsNo   int      `json:"errorsNo"`
	WarningsNo int      `json:"warningsNo"`
	Issues     []*Issue `json:"issues"`
	Content    []byte   `json:"content"`
}

// Linter represents linter details. It's used in metalinter.
type Linter struct {
	Name             string   `json:"name"`
	Command          string   `json:"command"`
	CompositeCommand string   `json:"composite_command,omitempty"`
	Pattern          string   `json:"pattern"`
	InstallFrom      string   `json:"install_from"`
	SeverityOverride Severity `json:"severity,omitempty"`
	MessageOverride  string   `json:"message_override,omitempty"`
	EventType        string

	Regex            *regexp.Regexp
}

// MarshalJSON marshals struct to JSON.
func (l *Linter) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Name)
}

// String returns formatted string.
func (l *Linter) String() string {
	return l.Name
}

// Severity level of the issue.
type Severity string

// Linter message severity levels.
const (
	Warning Severity = "warning"
	Error   Severity = "error"
)

// Issue represents issue in any of the go file.
type Issue struct {
	Linter   *Linter  `json:"linter"`
	Severity Severity `json:"severity"`
	Path     string   `json:"path"`
	Line     int      `json:"line"`
	Col      int      `json:"col"`
	Message  string   `json:"message"`
}

// String returns formatted string.
func (i *Issue) String() string {
	col := ""
	if i.Col != 0 {
		col = fmt.Sprintf("%d", i.Col)
	}
	return fmt.Sprintf("%s:%d:%s:%s: %s (%s)", strings.TrimSpace(i.Path), i.Line, col, i.Severity, strings.TrimSpace(i.Message), i.Linter)
}

