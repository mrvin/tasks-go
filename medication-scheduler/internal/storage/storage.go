package storage

import (
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
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
	ID           int64      `json:"id"`
	NameMedicine string     `json:"name_medicine"`
	NumPerDay    int16      `json:"num_per_day"`
	Times        []TimeOnly `json:"times"`
	AllLife      bool       `json:"all_life"`
	BeginDate    DateOnly   `json:"begin_date,omitempty"`
	EndDate      DateOnly   `json:"end_date,omitempty"`
	UserID       uuid.UUID  `json:"user_id"`
	Status       string     `json:"status"`
}

//nolint:tagliatelle
type Taking struct {
	NameMedicine string   `json:"name_medicine"`
	Time         TimeOnly `json:"time"`
}

type AllTaking struct {
	NameMedicine string
	Times        []TimeOnly
}

type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	date, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err //nolint:wrapcheck
	}
	*d = DateOnly(date)

	return nil
}

func (d *DateOnly) MarshalJSON() ([]byte, error) {
	return []byte("\"" + (*time.Time)(d).Format(time.DateOnly) + "\""), nil
}

type TimeOnly time.Time

func (t *TimeOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	timeOnly, err := time.Parse(time.TimeOnly, s)
	if err != nil {
		return err //nolint:wrapcheck
	}
	*t = TimeOnly(timeOnly)

	return nil
}

func (t *TimeOnly) MarshalJSON() ([]byte, error) {
	return []byte("\"" + (*time.Time)(t).Format(time.TimeOnly) + "\""), nil
}

type TimeOnlyArray []TimeOnly //nolint:recvcheck

// Scan implements the sql.Scanner interface.
func (a *TimeOnlyArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanTimes(src)
	case string:
		return a.scanTimes([]byte(src))
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("cannot convert %T to TimeOnlyArray", src)
}

func (a *TimeOnlyArray) scanTimes(src []byte) error {
	elems := bytes.Split(src[1:len(src)-1], []byte(","))
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(TimeOnlyArray, len(elems))
		for i, v := range elems {
			t, err := time.Parse(time.TimeOnly, string(v))
			if err != nil {
				return fmt.Errorf("parsing array element index %d: cannot convert", i)
			}
			b[i] = TimeOnly(t)
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a TimeOnlyArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil //nolint:nilnil
	}

	if n := len(a); n > 0 {
		b := make([]byte, 1, 1+8*n)
		b[0] = '{'

		t := (time.Time)(a[0]).Format(time.TimeOnly)
		b = append(b, []byte(t)...)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			t = (time.Time)(a[i]).Format(time.TimeOnly)
			b = append(b, []byte(t)...)
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}
