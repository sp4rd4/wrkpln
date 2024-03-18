package service

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sp4rd4/wrkpln/config"
	"github.com/sp4rd4/wrkpln/handler"
	"golang.org/x/sync/errgroup"
)

func Start(ctx context.Context, logger *slog.Logger, cfg config.Config) error {
	h := handler.New(logger)

	server := &http.Server{
		Addr:           ":" + strconv.Itoa(cfg.Port),
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
		Handler:        h,
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		return server.Shutdown(ctx)
	})

	slog.Info("http service: started")
	defer slog.Info("http service: stopped")

	return eg.Wait()
}
