package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/nikunicke/reaktorw/cmd/reaktorw/service"
	"github.com/nikunicke/reaktorw/cmd/reaktorw/service/updater"
	"github.com/nikunicke/reaktorw/warehouse/store/memory"
	"github.com/sirupsen/logrus"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// serviceRunner
	var serviceGroup service.Group
	var updaterConf updater.Config
	rootLogger := logrus.New()
	logger := rootLogger.WithFields(logrus.Fields{
		"app": "reaktor-warehouse",
	})
	logger.WithField("start_time", time.Now().Truncate(time.Second).String()).
		Info("starting reaktor-warehouse")
	warehouse := memory.NewInMemoryWarehouse()
	updaterConf.WarehouseAPI = warehouse
	updaterConf.UpdateInterval = 3 * time.Minute
	updaterConf.Logger = logger.WithField("service", "warehouse-updater")
	if s, err := updater.NewService(updaterConf); err == nil {
		serviceGroup = append(serviceGroup, s)
	} else {
		log.Fatal(err)
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGHUP)
		select {
		case s := <-sigCh:
			logger.WithField("signal", s.String()).Info("shutting down due to signal")
			cancelFn()
		case <-ctx.Done():
		}
	}()
	if err := serviceGroup.Run(ctx); err != nil {
		fmt.Println(err)
	}
	logger.WithField("end_time", time.Now().Truncate(time.Second).String()).Info("stopping reaktor-warehouse")
}
