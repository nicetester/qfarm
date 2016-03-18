package main

import (
	"encoding/json"
	"github.com/qfarm/qfarm/worker"
	"github.com/qfarm/qfarm"
	"fmt"
)

func main() {
	dir := "/home/md/.gvm/pkgsets/go1.6/global/src/github.com/qfarm/bad-go-code"

	fm, err := worker.BuildTree(dir)
	if err != nil {
		panic(err)
	}

	i1 := &qfarm.Issue{Path: "/home/md/.gvm/pkgsets/go1.6/global/src/github.com/qfarm/bad-go-code/logger/logger.go", Severity: qfarm.Error}
	i2 := &qfarm.Issue{Path: "/home/md/.gvm/pkgsets/go1.6/global/src/github.com/qfarm/bad-go-code/main.go", Severity: qfarm.Warning}
	i3 := &qfarm.Issue{Path: "/home/md/.gvm/pkgsets/go1.6/global/src/github.com/qfarm/bad-go-code/logger/logger.go", Severity: qfarm.Error}
	fm.ApplyIssue(i1)
	fm.ApplyIssue(i2)
	fm.ApplyIssue(i3)

	data, err := json.Marshal(fm.FilesMap)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON %+v", string(data))
}
