package sqllite

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/sp4rd4/wrkpln/planner"
	driver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func New(dbFilepath, schema string) (DB, error) {
	db, err := gorm.Open(driver.Open(dbFilepath))
	if err != nil {
		return DB{}, fmt.Errorf("open db: %w", err)
	}
	schemaSQL, err := os.ReadFile(schema)
	if err != nil {
		return DB{}, fmt.Errorf("read db schema: %w", err)
	}
	res := db.Exec(string(schemaSQL))
	if res.Error != nil {
		return DB{}, fmt.Errorf("load db schema: %w", err)
	}
	return DB{DB: db}, nil
}

func (db DB) CreateWorker(ctx context.Context, worker planner.Worker) error {
	res := db.WithContext(ctx).Create(worker)
	if res.Error != nil {
		return fmt.Errorf("create worker: %w", res.Error)
	}
	return nil
}

func (db DB) Worker(ctx context.Context, id uuid.UUID) (planner.Worker, error) {
	worker := planner.Worker{}
	res := db.WithContext(ctx).Take(&worker, "id = ?", id)
	switch {
	case errors.Is(res.Error, gorm.ErrRecordNotFound):
		return planner.Worker{}, planner.ErrNoRecord
	case res.Error != nil:
		return planner.Worker{}, fmt.Errorf("get worker: %w", res.Error)
	default:
		return worker, nil
	}
}

func (db DB) Workers(ctx context.Context, filter planner.WorkersFilter) ([]planner.Worker, error) {
	workers := []planner.Worker{}
	query := db.WithContext(ctx)
	if filter.Name != nil {
		query = query.Where("name LIKE ?", "%"+*filter.Name+"%")
	}

	res := query.Find(&workers)
	if res.Error != nil {
		return nil, fmt.Errorf("list workers: %w", res.Error)
	}
	return workers, nil
}

func (db DB) CreateShift(ctx context.Context, shift planner.Shift) error {
	res := db.WithContext(ctx).Create(shift)
	if res.Error != nil {
		return fmt.Errorf("create shift: %w", res.Error)
	}
	return nil
}

func (db DB) Shifts(ctx context.Context, filter planner.ShiftsFilter) ([]planner.Shift, error) {
	shifts := []planner.Shift{}
	query := db.WithContext(ctx)
	if filter.Date != nil {
		query = query.Where("date = ?", *filter.Date)
	}
	if filter.WorkerID != nil {
		query = query.Where("worker_id = ?", *filter.WorkerID)
	}

	res := query.Find(&shifts)
	if res.Error != nil {
		return nil, fmt.Errorf("list shifts: %w", res.Error)
	}
	return shifts, nil
}

func (db DB) Transaction(ctx context.Context, action func(planner.Repository) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDB := DB{DB: tx}
		return action(txDB)
	})
}
