package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
)

func (s *Storage) GetAllTaking(ctx context.Context, userID uuid.UUID, now time.Time) ([]storage.AllTaking, error) {
	slAllTaking := make([]storage.AllTaking, 0)

	rows, err := s.getAllTaking.QueryContext(ctx, userID, now)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return slAllTaking, nil
		}
		return nil, fmt.Errorf("can't get taking: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var taking storage.AllTaking
		err = rows.Scan(&taking.NameMedicine, pq.Array(&taking.Times))
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		slAllTaking = append(slAllTaking, taking)
	}
	if err := rows.Err(); err != nil {
		return slAllTaking, fmt.Errorf("rows error: %w", err)
	}

	return slAllTaking, nil
}
