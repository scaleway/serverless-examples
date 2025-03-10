package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

var (
	signalChan = make(chan os.Signal, 1)
	wg         sync.WaitGroup

	disableGracefulShutdown = os.Getenv("DISABLE_GRACEFUL_SHUTDOWN") != ""

	longRunningJobDuration        = os.Getenv("LONG_RUNNING_JOB_DURATION")
	parsedLongRunningJobDuration  time.Duration
	defaultLongRunningJobDuration = 14 * time.Minute
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))

	if longRunningJobDuration != "" {
		var err error

		parsedLongRunningJobDuration, err = time.ParseDuration(longRunningJobDuration)
		if err != nil {
			slog.Error("failed to parse LONG_RUNNING_JOB_DURATION", "error", err)
			os.Exit(1)
		}

		slog.Info("LONG_RUNNING_JOB_DURATION set", "duration", parsedLongRunningJobDuration)
	} else {
		parsedLongRunningJobDuration = defaultLongRunningJobDuration
		slog.Info("LONG_RUNNING_JOB_DURATION not set, using default", "duration", parsedLongRunningJobDuration)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handler),
	}

	// Handle SIGTERM.
	if !disableGracefulShutdown {
		signal.Notify(signalChan, syscall.SIGTERM)
	}

	// Start HTTP server.
	if disableGracefulShutdown {
		slog.Info("server started")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	} else {
		go func() {
			slog.Info("server started")

			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("server failed to start", "error", err)
				os.Exit(1)
			}
		}()

		sig := <-signalChan
		slog.Info("received signal", "signal", sig)

		// This should not take a lot of time because the handler is non-blocking.
		if err := srv.Shutdown(context.Background()); err != nil {
			slog.Error("server failed to shutdown", "error", err)
		}

		// This is where we wait for the long running tasks to finish.
		wg.Wait()
	}
}

func handler(w http.ResponseWriter, _ *http.Request) {
	wg.Add(1)

	taskID := uuid.NewString()
	go longBackgroundTask(taskID, parsedLongRunningJobDuration)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Started long background task: " + taskID))
}

func longBackgroundTask(taskID string, duration time.Duration) {
	defer wg.Done()

	slog.Info("long background task started", "taskID", taskID, "duration", duration)

	sleepTime := 1 * time.Second
	numIterations := int(duration / sleepTime)

	for i := range numIterations {
		slog.Info("long background task running...",
			"iteration", i,
			"total", numIterations,
			"taskID", taskID)

		time.Sleep(sleepTime)
	}

	slog.Info("long background task finished", "taskID", taskID)
}
