package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flowforge/backend/internal/events"
	"flowforge/backend/internal/httpapi"
	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
	"flowforge/backend/internal/runmanager"
	"flowforge/backend/internal/scheduler"
)

const (
	serverAddress   = ":8080"
	shutdownTimeout = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("application stopped with an error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	appCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	pipelineRepository := pipeline.NewMemoryRepository()
	pipelineService := pipeline.NewService(pipelineRepository)
	runRepository := runmodel.NewMemoryRepository()
	eventBus := events.NewBus(events.DefaultSubscriberBuffer)
	eventObserver := events.NewObserver(eventBus)
	snapshotCtx := context.WithoutCancel(appCtx)
	concurrentScheduler := scheduler.NewConcurrentScheduler(
		map[pipeline.ExecutorType]scheduler.Executor{
			pipeline.ExecutorTypeDelay: scheduler.DelayExecutor{},
		},
		func(snapshot runmodel.Run) error {
			return errors.Join(
				runRepository.Save(snapshotCtx, snapshot),
				eventObserver.Observe(snapshot),
			)
		},
	)
	runManager := runmanager.New(appCtx, pipelineService, runRepository, concurrentScheduler)
	router := httpapi.NewRouter(httpapi.Dependencies{
		Pipelines:      pipelineService,
		Runs:           runManager,
		AllowedOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
	})

	server := &http.Server{
		Addr:              serverAddress,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	slog.Info("HTTP server started", "address", serverAddress)

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil

	case <-appCtx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shut down HTTP server: %w", err)
	}
	runManager.Wait()
	eventBus.Close()

	if err := <-serverErr; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP after shutdown: %w", err)
	}

	slog.Info("HTTP server stopped")
	return nil
}
