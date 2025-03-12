package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (s *Storage) ListSchedulesIDs(ctx context.Context, userID uuid.UUID) ([]int64, error) {
	slListID := make([]int64, 0)

	rows, err := s.getListID.QueryContext(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return slListID, nil
		}
		return nil, fmt.Errorf("can't get list id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		slListID = append(slListID, id)
	}
	if err := rows.Err(); err != nil {
		return slListID, fmt.Errorf("rows error: %w", err)
	}

	return slListID, nil
}
