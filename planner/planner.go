package planner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrDayAlreadyBooked = Error("day already booked")
	ErrNoRecord         = Error("no record")
)

type Worker struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name" binding:"required"`
}

type WorkersFilter struct {
	Name *string `json:"name"`
}

type Shift struct {
	ID        uuid.UUID `json:"id"`
	WorkerID  uuid.UUID `json:"worker_id" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	StartHour int       `json:"start_hour" binding:"gte=0,lte=23"`
	EndHour   int       `json:"end_hour" binding:"gte=1,lte=24,gtfield=StartHour"`
}

type ShiftsFilter struct {
	WorkerID *uuid.UUID `json:"worker_id"`
	Date     *time.Time `json:"date"`
}

type Repository interface {
	CreateWorker(ctx context.Context, worker Worker) error
	Worker(ctx context.Context, id uuid.UUID) (Worker, error)
	Workers(ctx context.Context, filter WorkersFilter) ([]Worker, error)

	CreateShift(ctx context.Context, shift Shift) error
	Shifts(ctx context.Context, filter ShiftsFilter) ([]Shift, error)

	Transaction(ctx context.Context, action func(Repository) error) error
}

type Work struct {
	repo Repository
}

func New(repo Repository) Work {
	return Work{repo: repo}
}

func (w Work) CreateWorker(ctx context.Context, worker Worker) (Worker, error) {
	worker.ID = uuid.New()
	if err := w.repo.CreateWorker(ctx, worker); err != nil {
		return Worker{}, fmt.Errorf("creating worker: %w", err)
	}
	return worker, nil
}

func (w Work) Worker(ctx context.Context, id uuid.UUID) (Worker, error) {
	worker, err := w.repo.Worker(ctx, id)
	if err != nil {
		return Worker{}, fmt.Errorf("get worker: %w", err)
	}
	return worker, nil
}

func (w Work) Workers(ctx context.Context, filter WorkersFilter) ([]Worker, error) {
	workers, err := w.repo.Workers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list workers: %w", err)
	}
	return workers, nil
}

func (w Work) CreateShift(ctx context.Context, shift Shift) (Shift, error) {
	shift.ID = uuid.New()
	shift.Date = time.Date(
		shift.Date.Year(), shift.Date.Month(), shift.Date.Day(),
		0, 0, 0, 0, time.UTC,
	)
	err := w.repo.Transaction(ctx, func(repo Repository) error {
		_, err := repo.Worker(ctx, shift.WorkerID)
		if errors.Is(err, ErrNoRecord) {
			return fmt.Errorf("worker: %w", ErrNoRecord)
		}
		if err != nil {
			return fmt.Errorf("get worker: %w", err)
		}

		shifts, err := repo.Shifts(
			ctx, ShiftsFilter{WorkerID: &shift.WorkerID, Date: &shift.Date},
		)
		if err != nil {
			return fmt.Errorf("list shifts: %w", err)
		}
		if len(shifts) > 0 {
			return ErrDayAlreadyBooked
		}

		if err := repo.CreateShift(ctx, shift); err != nil {
			return fmt.Errorf("creating shift: %w", err)
		}
		return nil
	})
	if err != nil {
		return Shift{}, fmt.Errorf("create shift transaction: %w", err)
	}

	return shift, nil
}

func (w Work) Shifts(ctx context.Context, filter ShiftsFilter) ([]Shift, error) {
	if filter.Date != nil {
		date := time.Date(
			filter.Date.Year(), filter.Date.Month(), filter.Date.Day(),
			0, 0, 0, 0, time.UTC,
		)
		filter.Date = &date
	}
	shifts, err := w.repo.Shifts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list shift: %w", err)
	}
	return shifts, nil
}
