package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmi3midd/hws"
	"github.com/dmi3midd/shkvcache"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := hws.LoadConfig()
	if err != nil {
		slog.Error(
			"failed to load config",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	cache, err := shkvcache.NewCache[string](ctx, &shkvcache.Options{
		ShardCount:      8,
		CleanerInterval: 60,
		RunCleaner:      false,
	})
	defer cache.Close()
	if err != nil {
		slog.Error(
			"failed to create cache",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	service := hws.NewURLService(cache)

	server := hws.NewServer(
		cfg,
		service,
	)
	slog.Info(
		"server is running",
		slog.String("address", cfg.Address),
	)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(
				"failed to run server",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	slog.Info("received shutdown signal, stopping application...")

	server.Close()
}
