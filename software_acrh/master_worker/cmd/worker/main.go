package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"play-ground/software_acrh/master_worker/internal/worker/app"
	"play-ground/software_acrh/master_worker/internal/worker/platform/config"
)

func main() {
	cfg := config.MustLoad()

	w, err := app.Initialize(cfg)
	if err != nil {
		log.Fatalf("init worker failed: %v", err)
	}
	defer w.Close()

	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := w.Run(appCtx); err != nil {
		log.Fatalf("agent start error: %v", err)
	}

	slog.Info("worker is running!")
	<-appCtx.Done()
	slog.Info("app shutdown gracefully!")
}
