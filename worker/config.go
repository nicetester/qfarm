package worker

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"reflect"
)

type Cfg struct {
	// CheckLastCommitHash flag whether analysis should be done for the same commit hash - default false
	CheckLastCommitHash bool

	// RedisConn - Connection string for Redis - default: 127.0.0.1:6379
	RedisConn string

	// RedisPass - Password for Redis - default ""
	RedisPass string

	// Debug - Display messages for failed linters, etc. - default false
	Debug bool

	// Concurrency - Number of concurrent linters to run - default 16
	Concurrency int

	// Cyclo - Report functions with cyclomatic complexity over N (using gocyclo - default 10
	Cyclo int

	// Report lines longer than N (using lll) - default 120
	LLLineLength int

	// Minimum confidence interval to pass to golint - default 0.8
	GolintMinConfidence float64

	// Minimum occurrences to pass to goconst - default 3
	GoconstMinOccurrences int

	// Minimum token sequence as a clone for dupl - default 50
	DuplThreshold int
}

func NewDefaulConfig() *Cfg {
	return &Cfg{
		CheckLastCommitHash:   false,
		RedisConn:             "127.0.0.1:6379",
		RedisPass:             "",
		Debug:                 true,
		Concurrency:           16,
		Cyclo:                 10,
		LLLineLength:          120,
		GolintMinConfidence:   0.8,
		GoconstMinOccurrences: 3,
		DuplThreshold:         50,
	}
}

func Load(path string) (*Cfg, error) {
	conf := Cfg{}

	// load config from file
	if path == "" {
		return nil, fmt.Errorf("Empty path! Please provide path to load configuration!")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}

	conf.Print()

	return &conf, nil
}

func (c *Cfg) Print() {
	if c == nil {
		return
	}

	s := reflect.ValueOf(c).Elem()
	typeOfT := s.Type()

	log.Printf("Configuration:\n")
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		log.Printf("%s %s = %v\n", typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}