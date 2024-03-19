package handler

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sp4rd4/wrkpln/planner"
)

type PlanningHandler struct {
	*gin.Engine
	plan planner.Work
}

func New(logger *slog.Logger, plan planner.Work) PlanningHandler {
	h := PlanningHandler{Engine: gin.New(), plan: plan}
	setRoutes(h, logger)
	return h
}

func (h PlanningHandler) CreateWorker(c *gin.Context) {
	worker := planner.Worker{}
	if errorReturned := parseJson(c, &worker); !errorReturned {
		return
	}

	worker, err := h.plan.CreateWorker(c.Request.Context(), worker)
	if err != nil {
		var planErr planner.Error
		if errors.As(err, &planErr) {
			hadnlePlanningError(c, planErr)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("create worker error", "error", err)
		return
	}

	c.JSON(http.StatusCreated, worker)
}

func (h PlanningHandler) Workers(c *gin.Context) {
	wf := workersFilter(c.Request.URL.Query())
	workers, err := h.plan.Workers(c.Request.Context(), wf)
	if err != nil {
		var planErr planner.Error
		if errors.As(err, &planErr) {
			hadnlePlanningError(c, planErr)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("list workers error", "error", err)
		return
	}

	c.JSON(http.StatusOK, workers)
}

func workersFilter(query url.Values) planner.WorkersFilter {
	wf := planner.WorkersFilter{}
	if name := query.Get("name"); name != "" {
		wf.Name = &name
	}
	return wf
}

func (h PlanningHandler) CreateShift(c *gin.Context) {
	shift := planner.Shift{}
	if errorReturned := parseJson(c, &shift); !errorReturned {
		return
	}

	shift, err := h.plan.CreateShift(c.Request.Context(), shift)
	if err != nil {
		var planErr planner.Error
		if errors.As(err, &planErr) {
			hadnlePlanningError(c, planErr)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("create shift error", "error", err)
		return
	}

	c.JSON(http.StatusCreated, shift)
}

func (h PlanningHandler) Shifts(c *gin.Context) {
	sf, err := shiftsFilter(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shifts, err := h.plan.Shifts(c.Request.Context(), sf)
	if err != nil {
		var planErr planner.Error
		if errors.As(err, &planErr) {
			hadnlePlanningError(c, planErr)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("list shifts error", "error", err)
		return
	}

	c.JSON(http.StatusOK, shifts)
}

func shiftsFilter(query url.Values) (planner.ShiftsFilter, error) {
	sf := planner.ShiftsFilter{}
	if workerIDStr := query.Get("worker_id"); workerIDStr != "" {
		workerID, err := uuid.Parse(workerIDStr)
		if err != nil {
			return planner.ShiftsFilter{}, fmt.Errorf("worker_id: %w", err)
		}
		sf.WorkerID = &workerID
	}
	if dateStr := query.Get("date"); dateStr != "" {
		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return planner.ShiftsFilter{}, fmt.Errorf("date: %w", err)
		}
		sf.Date = &date
	}
	return sf, nil
}

func hadnlePlanningError(c *gin.Context, err planner.Error) {
	switch err {
	case planner.ErrDayAlreadyBooked:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case planner.ErrNoRecord:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
}

func jsonValidationError(errs validator.ValidationErrors) (string, error) {
	buf := bytes.NewBufferString("invalid fields:")
	for i, err := range errs {
		if i != len(errs)-1 {
			_, err := buf.WriteString(" " + err.Field() + ",")
			if err != nil {
				return "", fmt.Errorf("buffer write: %w", err)
			}
		}
		_, err := buf.WriteString(" " + err.Field() + ".")
		if err != nil {
			return "", fmt.Errorf("buffer write: %w", err)
		}
	}
	return buf.String(), nil
}

func parseJson(c *gin.Context, obj any) bool {
	err := c.ShouldBindJSON(obj)
	ve := validator.ValidationErrors{}
	if errors.As(err, &ve) {
		msg, err := jsonValidationError(ve)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			slog.Error("validate json error", "error", err)
			return false
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
		return false
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
