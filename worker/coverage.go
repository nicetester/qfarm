package worker

import (
	"bufio"
	"bytes"
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
	"github.com/qfarm/qfarm"
)

type CoverageChecker struct {
	cfg *Cfg
	notifier *Notifier
}

func NewCoverageChecker(cfg *Cfg, notifier *Notifier) *CoverageChecker {
	return &CoverageChecker{cfg: cfg, notifier: notifier}
}

func (c *CoverageChecker) Start(cfg qfarm.BuildCfg) error {
	if err := c.runCoverageAnalysis(cfg); err != nil {
		c.notifier.SendEvent(cfg.Repo, fmt.Sprintf("Coverage error in repo %s", cfg.Repo), EventTypeCoverageErr)
	}

	c.notifier.SendEvent(cfg.Repo, fmt.Sprintf("Coverage for repo %s done", cfg.Repo), EventTypeCoverageDone)
	return nil
}

func (c *CoverageChecker) runCoverageAnalysis(cfg qfarm.BuildCfg) error {
	packages := make([]qfarm.PackageReport, 0)

	// list all packages
	out, err := exec.Command("go", "list", path.Join(cfg.Repo, "...")).Output()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		packages = append(packages, qfarm.PackageReport{Name: scanner.Text()})
	}

	// run per package:
	// go tool cover
	// go test -cover
	for i, pac := range packages {
		c.debug("Starting coverage analysis of pkg: %s", pac.Name)
		start := time.Now()
		hash := fmt.Sprintf("%d", hash(pac.Name))
		cmd := exec.Command("bash", "-c", "go test -v -covermode=count -coverprofile=/tmp/"+hash+" "+pac.Name)
		var stdErr bytes.Buffer
		cmd.Stderr = &stdErr

		out, err := cmd.Output()
		if err != nil {
			warning("Some tests in package %s failed", pac.Name)
		}
		testOut := string(out)
		packages[i].Time = time.Now().Sub(start)

		if strings.Contains(testOut, "[no test files]") {
			c.debug("No tests for package(%s). Continueing", pac.Name)
			continue
		}

		// find coverage
		startI := strings.Index(testOut, "coverage:")
		endI := strings.Index(testOut, "% of statements")

		if startI == -1 || endI == -1 {
			c.debug("Can't parse test output of package(%s). Continueing", pac.Name)
			continue
		} else {
			startI = startI + 10
		}

		if startI >= endI || endI >= len(testOut) {
			c.debug("Can't parse test output of package(%s). Indexes are wriong. Continueing", pac.Name)
			continue
		}

		value := testOut[startI:endI]
		packages[i].Coverage, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		// find no of failed and passed tests
		packages[i].PassedNo = strings.Count(testOut, "--- PASS")
		packages[i].FailedNo = strings.Count(testOut, "--- FAIL")
		packages[i].TestsNo = packages[i].PassedNo + packages[i].FailedNo

		if packages[i].FailedNo > 0 {
			packages[i].Failed = true
		}

		cmd = exec.Command("go", "tool", "cover", "-html=/tmp/"+hash, "-o=/dev/stdout")
		cmd.Stderr = &stdErr
		out, err = cmd.Output()
		if err != nil {
			warning("Error in go tool cover command: %s", stdErr.String())
			return err
		}

		packages[i].Html = string(out)
		c.debug("Coverage anlysis of package(%s) Done", pac.Name)
	}

	report := qfarm.CoverageReport{Repo: cfg.Repo, Packages: packages}

	coverageAgg := 0.0
	for _, pkg := range report.Packages {
		coverageAgg += pkg.Coverage
		report.TotalFailedNo += pkg.FailedNo
		report.TotalPassedNo += pkg.PassedNo
		report.TotalTestsNo += pkg.TestsNo
		report.TotalTime += pkg.Time
	}

	report.TotalCoverage = coverageAgg / float64(len(report.Packages))

	if report.TotalFailedNo > 0 {
		report.Failed = true
	}

	// TODO: store report in redis here!

	return nil
}

func (m *CoverageChecker) debug(format string, args ...interface{}) {
	if m.cfg.Debug {
		fmt.Fprintf(os.Stderr, "DEBUG: "+format+"\n", args...)
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
