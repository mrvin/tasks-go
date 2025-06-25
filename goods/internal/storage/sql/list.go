package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/goods/internal/storage"
)

func (s *Storage) List(ctx context.Context, limit, offset uint64) ([]storage.Good, error) {
	sqlSelectList := `
		SELECT id, project_id, name, description, priority, removed, created_at
		FROM goods
		LIMIT $1
		OFFSET $2
	`
	goods := make([]storage.Good, 0)
	rows, err := s.db.QueryContext(ctx, sqlSelectList, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goods, nil
		}
		return nil, fmt.Errorf("can't get good: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var good storage.Good
		err = rows.Scan(
			&good.ID,
			&good.ProjectID,
			&good.Name,
			&good.Description,
			&good.Priority,
			&good.Removed,
			&good.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		goods = append(goods, good)
	}
	if err := rows.Err(); err != nil {
		return goods, fmt.Errorf("rows error: %w", err)
	}

	return goods, nil
}

func (s *Storage) Meta(ctx context.Context) (int64, int64, error) {
	var total, removed int64
	sqlSelectMeta := `
		SELECT 
			COUNT(*) AS total,
			COALESCE(SUM(CASE WHEN removed THEN 1 ELSE 0 END), 0)  AS removed
        FROM goods
	`
	err := s.db.QueryRowContext(ctx, sqlSelectMeta).Scan(&total, &removed)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get goods meta: %w", err)
	}

	return total, removed, nil
}
