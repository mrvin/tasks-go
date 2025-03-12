package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrScheduleNotFound = errors.New("schedule not found")

type SchedulerStorage interface {
	SaveSchedule(ctx context.Context, schedule *Schedule) (int64, error)
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
	BeginDate    Date      `json:"begin_date,omitempty"`
	EndDate      Date      `json:"end_date,omitempty"`
	UserID       uuid.UUID `json:"user_id"`
	Status       string    `json:"status"`
}

//nolint:tagliatelle
type Taking struct {
	NameMedicine string `json:"name_medicine"`
	Time         string `json:"time"`
}

type AllTaking struct {
	NameMedicine string
	Times        []int64
}
type Date struct {
	time.Time
}

func (t *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	date, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err //nolint:wrapcheck
	}
	t.Time = date

	return nil
}

func (t *Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.Time.Format(time.DateOnly) + "\""), nil
}
