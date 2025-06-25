package sqlstorage

import (
	"context"
	"fmt"

	"github.com/mrvin/tasks-go/goods/internal/storage"
)

func (s *Storage) Create(ctx context.Context, projectID int64, name, description string) (*storage.Good, error) {
	var good storage.Good
	sqlInsertGood := `
		INSERT INTO goods (
			project_id,
			name,
			description,
			priority
		)
		VALUES ($1, $2, $3, (SELECT COALESCE(MAX(priority), 0) FROM goods) + 1)
		RETURNING id, project_id, name, description, priority, removed, created_at
	`
	err := s.db.QueryRowContext(ctx, sqlInsertGood, projectID, name, description).Scan(
		&good.ID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("can't create good: %w", err)
	}

	return &good, nil
}
