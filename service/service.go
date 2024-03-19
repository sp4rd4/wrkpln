package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sp4rd4/wrkpln/config"
	handler "github.com/sp4rd4/wrkpln/handler/http"
	"github.com/sp4rd4/wrkpln/planner"
	"github.com/sp4rd4/wrkpln/repository/sqllite"
	"golang.org/x/sync/errgroup"
)

func Start(ctx context.Context, logger *slog.Logger, cfg config.Config) error {
	repo, err := sqllite.New(cfg.DBPath, cfg.DBSchemaPath)
	if err != nil {
		return fmt.Errorf("repository init: %w", err)
	}
	planner := planner.New(repo)
	h := handler.New(logger, planner)

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
			return fmt.Errorf("server listen: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}
		return nil
	})

	slog.Info("http service: started")
	defer slog.Info("http service: stopped")

	return eg.Wait()
}
