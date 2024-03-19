package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func setRoutes(handler PlanningHandler, logger *slog.Logger) {
	handler.Use(sloggin.New(logger), gin.Recovery())
	handler.POST("/worker", handler.CreateWorker, ContentTypeCheck)
	handler.GET("/worker/:id", handler.Worker)
	handler.GET("/workers", handler.Workers)

	handler.POST("/shift", handler.CreateShift, ContentTypeCheck)
	handler.GET("/shifts", handler.Shifts)

	handler.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 page not found"})
	})
	handler.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "405 method not allowed"})
	})
}

func ContentTypeCheck(c *gin.Context) {
	if c.ContentType() != "application/json" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "415 unsupported media type"})
		return
	}
	c.Next()
}
