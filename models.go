package qfarm

import "time"

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
