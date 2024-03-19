package planner_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/sp4rd4/wrkpln/planner"
	repomock "github.com/sp4rd4/wrkpln/repository/mock"
)

func ptr[T any](v T) *T {
	return &v
}

var fixedID = uuid.New()

func genID() uuid.UUID {
	return fixedID
}

type transaction func(planner.Repository) error

func TestCreateWorker(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   planner.Worker
		want    planner.Worker
		repoErr error
		expErr  error
	}{
		{
			name:  "Success",
			input: planner.Worker{Name: "Buddy Guy"},
			want:  planner.Worker{ID: fixedID, Name: "Buddy Guy"},
		},
		{
			name:    "Repo error",
			input:   planner.Worker{Name: "Buddy Guy"},
			want:    planner.Worker{},
			repoErr: net.UnknownNetworkError("error"),
			expErr:  net.UnknownNetworkError("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := repomock.NewMockRepository(ctrl)

			plan := planner.New(repo, planner.UUIDGenerator(genID))
			expected := tt.input
			expected.ID = fixedID
			repo.EXPECT().CreateWorker(ctx, expected).Return(tt.repoErr)

			result, err := plan.CreateWorker(ctx, tt.input)
			if tt.expErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, result, tt.want)
			} else {
				assert.ErrorIs(t, err, tt.expErr)
			}

		})
	}
}

func TestWorkers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   planner.WorkersFilter
		want    []planner.Worker
		repoErr error
		expErr  error
	}{
		{
			name:  "Success",
			input: planner.WorkersFilter{Name: ptr("Buddy")},
			want:  []planner.Worker{{ID: fixedID, Name: "Buddy Guy"}, {ID: fixedID, Name: "Buddy Friend"}},
		},
		{
			name:    "Repo error",
			input:   planner.WorkersFilter{Name: ptr("Buddy")},
			want:    nil,
			repoErr: net.UnknownNetworkError("error"),
			expErr:  net.UnknownNetworkError("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := repomock.NewMockRepository(ctrl)

			plan := planner.New(repo, planner.UUIDGenerator(genID))
			repo.EXPECT().Workers(ctx, tt.input).Return(tt.want, tt.repoErr)

			result, err := plan.Workers(ctx, tt.input)
			if tt.expErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, result, tt.want)
			} else {
				assert.ErrorIs(t, err, tt.expErr)
			}

		})
	}
}

func TestCreateShift(t *testing.T) {
	t.Parallel()
	id1 := uuid.New()
	date := time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		input        planner.Shift
		worker       planner.Worker
		workerShifts []planner.Shift
		want         planner.Shift

		repoWorkerErr  error
		repoShiftsErr  error
		repoCreaterErr error
		expErr         error
	}{
		{
			name:         "Success",
			input:        planner.Shift{WorkerID: id1, Date: date, StartHour: 8, EndHour: 16},
			worker:       planner.Worker{ID: id1, Name: "Buddy Guy"},
			want:         planner.Shift{ID: fixedID, WorkerID: id1, Date: date, StartHour: 8, EndHour: 16},
			workerShifts: nil,
		},
		{
			name:         "Day booked",
			input:        planner.Shift{WorkerID: id1, Date: date, StartHour: 8, EndHour: 16},
			worker:       planner.Worker{ID: id1, Name: "Buddy Guy"},
			workerShifts: []planner.Shift{{WorkerID: id1, Date: date, StartHour: 0, EndHour: 8}},
			expErr:       planner.ErrDayAlreadyBooked,
		},
		{
			name:          "No such worker",
			input:         planner.Shift{WorkerID: id1, Date: date, StartHour: 8, EndHour: 16},
			worker:        planner.Worker{},
			repoWorkerErr: planner.ErrNoRecord,
			expErr:        planner.ErrNoRecord,
		},
		{
			name:          "Shifts repo err",
			input:         planner.Shift{WorkerID: id1, Date: date, StartHour: 8, EndHour: 16},
			worker:        planner.Worker{ID: id1, Name: "Buddy Guy"},
			repoShiftsErr: net.UnknownNetworkError("error"),
			expErr:        net.UnknownNetworkError("error"),
		},
		{
			name:           "Create shift repo err",
			input:          planner.Shift{WorkerID: id1, Date: date, StartHour: 8, EndHour: 16},
			worker:         planner.Worker{ID: id1, Name: "Buddy Guy"},
			repoCreaterErr: net.UnknownNetworkError("error"),
			expErr:         net.UnknownNetworkError("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := repomock.NewMockRepository(ctrl)

			plan := planner.New(repo, planner.UUIDGenerator(genID))
			expectations := []any{
				repo.EXPECT().Transaction(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, f transaction) error { return f(repo) }),
				repo.EXPECT().Worker(ctx, id1).Return(tt.worker, tt.repoWorkerErr),
			}
			if tt.repoWorkerErr == nil {
				expectations = append(
					expectations,
					repo.EXPECT().Shifts(ctx, planner.ShiftsFilter{WorkerID: &id1, Date: &date}).Return(tt.workerShifts, tt.repoShiftsErr),
				)
				if tt.repoShiftsErr == nil && len(tt.workerShifts) == 0 {
					expected := tt.input
					expected.ID = fixedID
					expectations = append(
						expectations,
						repo.EXPECT().CreateShift(ctx, expected).Return(tt.repoCreaterErr),
					)
				}
			}
			gomock.InOrder(expectations...)

			result, err := plan.CreateShift(ctx, tt.input)
			if tt.expErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, result, tt.want)
			} else {
				assert.ErrorIs(t, err, tt.expErr)
			}

		})
	}
}

func TestShifts(t *testing.T) {
	t.Parallel()
	date := time.Date(2025, 11, 3, 14, 22, 15, 0, time.UTC)
	dateTrunc := time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC)
	id1 := uuid.New()
	tests := []struct {
		name    string
		input   planner.ShiftsFilter
		want    []planner.Shift
		repoErr error
		expErr  error
	}{
		{
			name:  "Success",
			input: planner.ShiftsFilter{WorkerID: &fixedID, Date: &date},
			want:  []planner.Shift{{ID: id1, Date: dateTrunc, StartHour: 8, EndHour: 16}},
		},
		{
			name:    "Repo error",
			input:   planner.ShiftsFilter{WorkerID: &fixedID, Date: &date},
			want:    nil,
			repoErr: net.UnknownNetworkError("error"),
			expErr:  net.UnknownNetworkError("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := repomock.NewMockRepository(ctrl)

			plan := planner.New(repo, planner.UUIDGenerator(genID))
			tt.input.Date = &dateTrunc
			repo.EXPECT().Shifts(ctx, tt.input).Return(tt.want, tt.repoErr)

			result, err := plan.Shifts(ctx, tt.input)
			if tt.expErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, result, tt.want)
			} else {
				assert.ErrorIs(t, err, tt.expErr)
			}

		})
	}
}
