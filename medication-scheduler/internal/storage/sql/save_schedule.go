package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
)

func (s *Storage) SaveSchedule(ctx context.Context, schedule *storage.Schedule) (int64, error) {
	if err := s.insertSchedule.QueryRowContext(ctx,
		schedule.NameMedicine,
		schedule.NumPerDay,
		storage.TimeOnlyArray(schedule.Times),
		schedule.AllLife,
		time.Time(schedule.BeginDate),
		time.Time(schedule.EndDate),
		schedule.UserID,
	).Scan(&schedule.ID); err != nil {
		return 0, fmt.Errorf("saving schedule to db: %w", err)
	}

	return schedule.ID, nil
}
