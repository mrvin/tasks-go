package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/goods/internal/storage"
)

func (s *Storage) Delete(ctx context.Context, id, projectID int64) (*storage.Good, error) {
	sqlDeleteGood := `
		UPDATE goods
		SET removed = true
		WHERE id = $1 AND project_id = $2
		RETURNING id, project_id, name, description, priority, removed, created_at
	`
	var good storage.Good
	err := s.db.QueryRowContext(ctx, sqlDeleteGood, id, projectID).Scan(
		&good.ID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w %v", storage.ErrNoGoodID, id)
		}
		return nil, fmt.Errorf("can't get good: %w", err)
	}

	return &good, nil
}
