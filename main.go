package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/sp4rd4/wrkpln/config"
	"github.com/sp4rd4/wrkpln/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		slog.Error("Failed to read config", "error", err)
		os.Exit(1)
	}
	logLevel := new(slog.Level)
	if err := logLevel.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		slog.Error("Incorrect log level", "error", err)
	}
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	logger := slog.New(h)
	slog.SetDefault(logger)

	if err := service.Start(ctx, logger, cfg); err != nil {
		slog.Error("server error", "error", err)
	}
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		slog.Info("Received interrupt signal")
		stop()
	}()
}
