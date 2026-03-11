package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

type Worker struct {
	cron *cron.Cron
	job  func(ctx context.Context) error
}

func (w *Worker) Start(ctx context.Context, schedule string) error {
	_, err := w.cron.AddFunc(schedule, func() {
		slog.Info("cron triggered: starting job")

		jobCtx, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()

		if err := w.job(jobCtx); err != nil {
			slog.Error("do job failed", "error_msg", err.Error())
			return
		}
		slog.Info("job done")
	})

	if err != nil {
		return err
	}

	w.cron.Start()
	slog.Info("cron scheduler started", "schedule", schedule)
	return nil
}

func (w *Worker) Stop() {
	ctx := w.cron.Stop()
	<-ctx.Done()
	slog.Info("cron scheduler stopped")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := &Worker{
		cron: cron.New(),
		job: func(jobCtx context.Context) error {
			slog.Info("working hard...")

			// The job MUST listen to jobCtx.Done() to actually stop when the timeout occurs.
			// Here we simulate work that takes 5 seconds (which is longer than our 2s timeout).
			select {
			case <-time.After(5 * time.Second):
				slog.Info("work finished naturally")
				return nil

			case <-jobCtx.Done():
				// If the timeout is reached, jobCtx.Done() receives a signal.
				// We return the context error (usually context.DeadlineExceeded).
				return jobCtx.Err()
			}
		},
	}

	if err := w.Start(ctx, "@every 5s"); err != nil {
		slog.Error("failed to start worker", "error", err.Error())
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	slog.Info("shutting down application...")

	cancel()
	w.Stop()
	slog.Info("application exited")
}
