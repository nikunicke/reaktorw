package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/nikunicke/reaktorw/cmd/reaktorw/service"
	"github.com/nikunicke/reaktorw/cmd/reaktorw/service/frontend"
	"github.com/nikunicke/reaktorw/cmd/reaktorw/service/updater"
	"github.com/nikunicke/reaktorw/warehouse/store/memory"
	"github.com/sirupsen/logrus"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	// cpu profiling
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	mainlogger := logrus.New()
	logger := mainlogger.WithField("app", "reaktor-warehouse")
	if err := runApp(logger); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"exit_time": time.Now().String(),
		}).Error("Exiting due to error")
	}

	logger.WithField("exit_time", time.Now().Truncate(time.Second).String()).Info("stopping reaktor-warehouse")
}

func runApp(logger *logrus.Entry) error {
	serviceGroup, err := setupServices(logger)
	if err != nil {
		return err
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
	logger.WithField("start_time", time.Now().String()).Info("starting app")
	return serviceGroup.Run(ctx)
}

func setupServices(logger *logrus.Entry) (service.Group, error) {
	var (
		updaterConf  updater.Config
		frontendConf frontend.Config

		serviceGroup service.Group
	)

	// warehouse
	warehouse := memory.NewInMemoryWarehouse()

	// updater
	updaterConf.WarehouseAPI = warehouse
	updaterConf.UpdateInterval = 5 * time.Minute
	updaterConf.Logger = logger.WithField("service", "warehouse-updater")
	if service, err := updater.NewService(updaterConf); err == nil {
		serviceGroup = append(serviceGroup, service)
	} else {
		return nil, err
	}
	// frontend
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	frontendConf.WarehouseAPI = warehouse
	frontendConf.ListenAddr = ":" + port
	frontendConf.Logger = logger.WithField("service", "frontend")
	if service, err := frontend.NewService(frontendConf); err == nil {
		serviceGroup = append(serviceGroup, service)
	} else {
		return nil, err
	}
	return serviceGroup, nil
}
