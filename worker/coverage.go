package worker

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/qfarm/qfarm"
)

type CoverageChecker struct {
	cfg      *Cfg
	notifier *Notifier
}

func NewCoverageChecker(cfg *Cfg, notifier *Notifier) *CoverageChecker {
	return &CoverageChecker{cfg: cfg, notifier: notifier}
}

func (c *CoverageChecker) Start(cfg qfarm.BuildCfg, ft *FilesMap) error {
	report, err := c.runCoverageAnalysis(cfg)
	if err != nil {
		c.notifier.SendEvent(cfg.Repo, fmt.Sprintf("Coverage error in repo %s", cfg.Repo), EventTypeCoverageErr)
	}

	c.notifier.SendEvent(cfg.Repo, fmt.Sprintf("Coverage for repo %s done", cfg.Repo), EventTypeCoverageDone)

	if report.Failed {
		return nil
	}
	if err := ft.ApplyCover(report); err != nil {
		return err
	}

	return nil
}

func (c *CoverageChecker) runCoverageAnalysis(cfg qfarm.BuildCfg) (*qfarm.CoverageReport, error) {
	packages := make([]qfarm.PackageReport, 0)

	// list all packages
	out, err := exec.Command("go", "list", path.Join(cfg.Repo, "...")).Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		packages = append(packages, qfarm.PackageReport{Name: scanner.Text()})
	}

	// run per package:
	// go tool cover
	// go test -cover
	var rootTotal, rootCovered int64
	for i, pac := range packages {
		c.debug("Starting coverage analysis of pkg: %s", pac.Name)
		start := time.Now()
		stamp := fmt.Sprint(start.Nanosecond())
		cmd := exec.Command("bash", "-c", "go test -v -covermode=set -coverprofile=/tmp/"+stamp+" "+pac.Name)
		var stdErr bytes.Buffer
		cmd.Stderr = &stdErr

		out, err := cmd.Output()
		if err != nil {
			warning("Some tests in package %s failed", pac.Name)
		}
		testOut := string(out)
		packages[i].Time = time.Now().Sub(start)

		if strings.Contains(testOut, "[no test files]") {
			c.debug("No tests for package(%s). Continuing", pac.Name)
			continue
		}

		// find coverage
		startI := strings.Index(testOut, "coverage:")
		endI := strings.Index(testOut, "% of statements")

		if startI == -1 || endI == -1 {
			c.debug("Can't parse test output of package(%s). Continuing", pac.Name)
			continue
		} else {
			startI = startI + 10
		}

		if startI >= endI || endI >= len(testOut) {
			c.debug("Can't parse test output of package(%s). Indexes are wrong. Continuing", pac.Name)
			continue
		}

		value := testOut[startI:endI]
		packages[i].Coverage, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}

		cp, err := os.Open("/tmp/" + stamp)
		if err != nil {
			return nil, fmt.Errorf("can't open coverage profile file: %v", err)
		}
		defer cp.Close()

		scanner := bufio.NewScanner(cp)
		scanner.Scan() // Skip file header.
		files := make(map[string]qfarm.CoverFileReport)
		for scanner.Scan() {
			line := strings.Split(strings.TrimPrefix(scanner.Text(), pac.Name+"/"), ":")
			var fileName, blocks string
			if len(line) > 1 {
				fileName = line[0]
				blocks = line[1]
			}
			block := strings.Split(blocks, " ")
			var cursors string
			var numStmt, count int
			if len(block) > 2 {
				cursors = block[0]
				numStmt, err = strconv.Atoi(block[1])
				if err != nil {
					return nil, fmt.Errorf("can't parse NumStmt: %v", err)
				}
				count, err = strconv.Atoi(block[2])
				if err != nil {
					return nil, fmt.Errorf("can't parse Count: %v", err)
				}
			}
			curs := strings.Split(cursors, ",")
			if len(curs) < 2 {
				return nil, fmt.Errorf("bad curs len: %+v", curs)
			}
			var start, end qfarm.Cursor
			st := strings.Split(curs[0], ".")
			if len(st) > 1 {
				start.Line, err = strconv.Atoi(st[0])
				if err != nil {
					return nil, fmt.Errorf("can't parse Line: %v", err)
				}
				start.Col, err = strconv.Atoi(st[1])
				if err != nil {
					return nil, fmt.Errorf("can't parse Col: %v", err)
				}
			}
			en := strings.Split(curs[1], ".")
			if len(en) > 1 {
				end.Line, err = strconv.Atoi(en[0])
				if err != nil {
					return nil, fmt.Errorf("can't parse Line: %v", err)
				}
				end.Col, err = strconv.Atoi(en[1])
				if err != nil {
					return nil, fmt.Errorf("can't parse Col: %v", err)
				}
			}
			coverBlock := qfarm.CoverBlock{Start: start, End: end, NumStmt: int64(numStmt), Count: int64(count)}

			file := files[fileName]
			file.Blocks = append(file.Blocks, coverBlock)
			files[fileName] = file
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("scanner error while reading coverage profile file: %v", err)
		}

		var pacTotal, pacCovered int64
		for f, report := range files {
			var total, covered int64
			for _, b := range report.Blocks {
				total += b.NumStmt
				pacTotal += b.NumStmt
				rootTotal += b.NumStmt
				if b.Count > 0 {
					covered += b.NumStmt
					pacCovered += b.NumStmt
					rootCovered += b.NumStmt
				}
			}
			if total == 0 {
				report.Coverage = 0
				files[f] = report
				continue
			}
			report.Coverage = float64(covered) / float64(total) * 100
			files[f] = report
		}

		if pacTotal > 0 {
			packages[i].Coverage = float64(pacCovered) / float64(pacTotal) * 100
		}

		// find no of failed and passed tests
		packages[i].PassedNo = strings.Count(testOut, "--- PASS")
		packages[i].FailedNo = strings.Count(testOut, "--- FAIL")
		packages[i].TestsNo = packages[i].PassedNo + packages[i].FailedNo
		packages[i].Files = files

		if packages[i].FailedNo > 0 {
			packages[i].Failed = true
		}

		c.debug("Coverage anlysis of package(%s): %f", pac.Name, packages[i].Coverage)
	}

	report := qfarm.CoverageReport{Repo: cfg.Repo, Packages: packages}

	for _, pkg := range report.Packages {
		report.TotalFailedNo += pkg.FailedNo
		report.TotalPassedNo += pkg.PassedNo
		report.TotalTestsNo += pkg.TestsNo
		report.TotalTime += pkg.Time
	}

	if rootTotal > 0 {
		report.TotalCoverage = float64(rootCovered) / float64(rootTotal) * 100
	}
	c.debug("Root total coverage: %f", report.TotalCoverage)

	if report.TotalFailedNo > 0 {
		report.Failed = true
	}

	return &report, nil
}

func (m *CoverageChecker) debug(format string, args ...interface{}) {
	if m.cfg.Debug {
		fmt.Fprintf(os.Stderr, "DEBUG: "+format+"\n", args...)
	}
}
