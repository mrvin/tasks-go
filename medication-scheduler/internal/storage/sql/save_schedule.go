package sqlstorage

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
)

func (s *Storage) SaveSchedule(ctx context.Context, schedule *storage.Schedule) (int64, error) {
	if err := s.insertSchedule.QueryRowContext(ctx,
		schedule.NameMedicine,
		schedule.NumPerDay,
		pq.Array(schedule.TimesInt64),
		schedule.AllLife,
		schedule.BeginDate.Time,
		schedule.EndDate.Time,
		schedule.UserID,
	).Scan(&schedule.ID); err != nil {
		return 0, fmt.Errorf("saving schedule to db: %w", err)
	}

	return schedule.ID, nil
}
