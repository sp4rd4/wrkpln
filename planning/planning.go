package planning

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrDayAlreadyBooked = Error("day_already_booked")
	ErrNoRecord         = Error("no_record")
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
	StartHour int       `json:"start_hour" binding:"required"`
	EndHour   int       `json:"end_hour" binding:"required"`
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

type Planner struct {
	repo Repository
}

func (p Planner) CreateWorker(ctx context.Context, worker Worker) (Worker, error) {
	worker.ID = uuid.New()
	if err := p.repo.CreateWorker(ctx, worker); err != nil {
		return Worker{}, fmt.Errorf("creating worker: %w", err)
	}
	return worker, nil
}

func (p Planner) Worker(ctx context.Context, id uuid.UUID) (Worker, error) {
	worker, err := p.repo.Worker(ctx, id)
	if err != nil {
		return Worker{}, fmt.Errorf("get worker: %w", err)
	}
	return worker, nil
}

func (p Planner) Workers(ctx context.Context, filter WorkersFilter) ([]Worker, error) {
	workers, err := p.repo.Workers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list workers: %w", err)
	}
	return workers, nil
}

func (p Planner) CreateShift(ctx context.Context, shift Shift) (Shift, error) {
	shift.ID = uuid.New()
	err := p.repo.Transaction(ctx, func(repo Repository) error {
		shifts, err := repo.Shifts(ctx, ShiftsFilter{WorkerID: &shift.WorkerID, Date: &shift.Date})
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

func (p Planner) Shifts(ctx context.Context, filter ShiftsFilter) ([]Shift, error) {
	shifts, err := p.repo.Shifts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list shift: %w", err)
	}
	return shifts, nil
}
