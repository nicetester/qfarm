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
	Repo          string
	TotalCoverage float64
	TotalTestsNo  int
	TotalPassedNo int
	TotalFailedNo int
	TotalTime     time.Duration
	Failed        bool
	Packages      []PackageReport
}

// PackageReport holds info about coverage analysis of specified package.
type PackageReport struct {
	Name     string
	Coverage float64
	Failed   bool
	TestsNo  int
	PassedNo int
	FailedNo int
	Time     time.Duration
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
	NumStmt int64  `json:"numStmt"`
	Count   int64  `json:"count"`
}

// Cursor points at specific location in file and is used by CoverBlock.
type Cursor struct {
	Line int `json:"line"`
	Col  int `json:"col"`
}

// Node represents a node in directory tree. It might be a file or a directory.
type Node struct {
	Path       string       `json:"path"`
	Nodes      []Node       `json:"nodes"`
	ParentPath string       `json:"parent"`
	Dir        bool         `json:"dir"`
	Coverage   float64      `json:"coverage"`
	Blocks     []CoverBlock `json:"coverBlocks"`
	TestsNo    int          `json:"testsNo"`
	FailedNo   int          `json:"failedNo"`
	PassedNo   int          `json:"passedNo"`
	IssuesNo   int          `json:"issuesNo"`
	ErrorsNo   int          `json:"errorsNo"`
	WarningsNo int          `json:"warningsNo"`
	Issues     []*Issue     `json:"issues"`
	Content    []byte       `json:"content"`
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

	Regex *regexp.Regexp
}

// MarshalJSON marshals struct to JSON.
func (l Linter) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Name)
}

// UnmarshalJSON marshals struct to JSON.
func (l *Linter) UnmarshalJSON(data []byte) error {
	if data == nil {
		return fmt.Errorf("Empty byte slice")
	}

	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	*l = Linter{Name: str}
	return nil
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

func (s *Severity) Rank() int {
	if *s == Warning {
		return 1
	}

	if *s == Error {
		return 2
	}

	return -1
}

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

// Reports stores information about whole analysis.
type Report struct {
	Repo       string    `json:"repo"`
	No         int       `json:"no"`
	Score      int       `json:"score"`
	Time       JSONTime  `json:"time"`
	Took       string    `json:"took"`
	CommitHash string    `json:"commitHash"`
	Config     BuildCfg  `json:"config"`

	Coverage          float64 `json:"coverage"`
	TestsNo           int     `json:"testsNo"`
	FailedNo          int     `json:"failedNo"`
	PassedNo          int     `json:"passedNo"`
	IssuesNo          int     `json:"issuesNo"`
	ErrorsNo          int     `json:"errorsNo"`
	WarningsNo        int     `json:"warningsNo"`
	TechnicalDeptCost int     `json:"technicalDeptCost"`
	TechnicalDeptTime string  `json:"technicalDeptTime"`
}

const defaultDateFormat = "2006-01-02 15:04:05"

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(`"` + defaultDateFormat + `"`)), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	val, err := time.Parse(`"`+defaultDateFormat+`"`, string(data))
	if err != nil {
		return err
	}

	*t = JSONTime(val)
	return nil
}

func (t *JSONTime) Time() time.Time {
	return time.Time(*t)
}

func (t *JSONTime) String() string {
	if t == nil {
		return ""
	}

	return time.Time(*t).Format(defaultDateFormat)
}