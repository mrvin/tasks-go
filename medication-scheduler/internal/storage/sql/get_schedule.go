package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
)

func (s *Storage) GetSchedule(ctx context.Context, userID uuid.UUID, scheduleID int64) (*storage.Schedule, error) {
	var schedule storage.Schedule

	if err := s.getSchedule.QueryRowContext(ctx, userID, scheduleID).Scan(
		&schedule.ID,
		&schedule.NameMedicine,
		&schedule.NumPerDay,
		(*storage.TimeOnlyArray)(&schedule.Times),
		&schedule.AllLife,
		(*time.Time)(&schedule.BeginDate),
		(*time.Time)(&schedule.EndDate),
		&schedule.UserID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: %d", storage.ErrScheduleNotFound, scheduleID)
		}
		return nil, fmt.Errorf("get schedule with id: %d: %w", scheduleID, err)
	}

	return &schedule, nil
}
