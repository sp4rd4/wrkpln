package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type PlanningHandler struct {
	*gin.Engine
}

type Option func(*PlanningHandler)

func New(logger *slog.Logger, options ...Option) PlanningHandler {
	h := PlanningHandler{Engine: gin.New()}
	for _, opt := range options {
		opt(&h)
	}
	h.Use(sloggin.New(logger), gin.Recovery())
	return h
}
