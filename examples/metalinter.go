package main

import (
	"github.com/qfarm/qfarm/worker"
	"time"
	"fmt"
)

func main() {
	cfg := worker.NewDefaulConfig()
	cfg.Debug = false
	notifier := worker.NewNotifier(nil)
	metalinter := worker.NewMetalinter(cfg, notifier)
	//
	//	if err := metalinter.InstallAllLinters(); err != nil {
	//		panic(err)
	//	}

	repoCfg, err := worker.LoadRepoCfg("github.com/qfarm/bad-go-code", "/home/md/.gvm/pkgsets/go1.6/global/src/github.com/qfarm/bad-go-code")
	fmt.Printf("Build config:  %+v ", repoCfg)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	if err := metalinter.Start(*repoCfg); err != nil {
		panic(err)
	}

	fmt.Printf("Analysis done! Took: %v", time.Now().Sub(start))
}
