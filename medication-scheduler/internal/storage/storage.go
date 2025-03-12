package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrScheduleNotFound = errors.New("schedule not found")

type SchedulerStorage interface {
	CreateSchedule(ctx context.Context, schedule *Schedule) (int64, error)
	GetSchedule(ctx context.Context, userID uuid.UUID, scheduleID int64) (*Schedule, error)

	ListSchedulesIDs(ctx context.Context, userID uuid.UUID) ([]int64, error)
	GetAllTaking(ctx context.Context, userID uuid.UUID, now time.Time) ([]AllTaking, error)
}

//nolint:tagliatelle
type Schedule struct {
	ID           int64     `json:"id"`
	NameMedicine string    `json:"name_medicine"`
	NumPerDay    int16     `json:"num_per_day"`
	Times        []string  `json:"times"`
	TimesInt64   []int64   `json:"-"`
	AllLife      bool      `json:"all_life"`
	BeginDate    time.Time `json:"begin_date,omitempty"`
	EndDate      time.Time `json:"end_date,omitempty"`
	UserID       uuid.UUID `json:"user_id"`
	Status       string    `json:"status"`
}

type AllTaking struct {
	NameMedicine string
	Times        []int64
}

//nolint:tagliatelle
type Taking struct {
	NameMedicine string `json:"name_medicine"`
	Time         string `json:"time"`
}
