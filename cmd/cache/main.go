package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"vc/internal/cache/apiv1"
	"vc/internal/cache/httpserver"
	"vc/internal/cache/simplequeue"
	"vc/pkg/configuration"
	"vc/pkg/kvclient"
	"vc/pkg/logger"
	"vc/pkg/trace"
)

type service interface {
	Close(ctx context.Context) error
}

func main() {
	var wg sync.WaitGroup
	ctx := context.Background()

	services := make(map[string]service)

	cfg, err := configuration.Parse(ctx, logger.NewSimple("Configuration"))
	if err != nil {
		panic(err)
	}

	log, err := logger.New("vc_cache", cfg.Common.Log.FolderPath, cfg.Common.Production)
	if err != nil {
		panic(err)
	}
	tracer, err := trace.New(ctx, cfg, log, "vc", "cache")
	if err != nil {
		panic(err)
	}

	kvClient, err := kvclient.New(ctx, cfg, tracer, log.New("kvClient"))
	services["kvService"] = kvClient
	if err != nil {
		panic(err)
	}

	queueService, err := simplequeue.New(ctx, kvClient, tracer, cfg, log.New("queue"))
	services["queueService"] = queueService
	if err != nil {
		panic(err)
	}

	apiv1Client, err := apiv1.New(ctx, kvClient, tracer, cfg, log.New("apiv1"))
	if err != nil {
		panic(err)
	}
	httpService, err := httpserver.New(ctx, cfg, apiv1Client, tracer, log.New("httpserver"))
	services["httpService"] = httpService
	if err != nil {
		panic(err)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan // Blocks here until interrupted

	mainLog := log.New("main")
	mainLog.Info("HALTING SIGNAL!")

	for serviceName, service := range services {
		if err := service.Close(ctx); err != nil {
			mainLog.Trace("serviceName", serviceName, "error", err)
		}
	}

	if err := tracer.Shutdown(ctx); err != nil {
		mainLog.Error(err, "Tracer shutdown")
	}

	wg.Wait() // Block here until are workers are done

	mainLog.Info("Stopped")
}
