package main
import (
	"fmt"
	"github.com/qfarm/qfarm/worker"
)

func main() {
	cfg := worker.NewDefaulConfig()
	notifier := worker.NewNotifier(nil)

	repo := "github.com/hashicorp/consul/api"

	checker := worker.NewCoverageChecker(cfg, notifier)

	repoCfg, err := worker.LoadRepoCfg(repo, "/home/md/.gvm/pkgsets/go1.6/global/src/" + repo)
	if err != nil {
		panic(err)
	}

	r, err := checker.RunCoverageAnalysis(*repoCfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", r)
}
