package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/mrvin/tasks-go/goods/internal/storage"
)

func (s *Storage) Update(ctx context.Context, id, projectID int64, name, description string) (*storage.Good, error) {
	// Create sql query
	sqlUpdate := squirrel.Update("goods")
	sqlUpdate = sqlUpdate.Where(squirrel.Eq{"id": id})
	sqlUpdate = sqlUpdate.Where(squirrel.Eq{"project_id": projectID})
	sqlUpdate = sqlUpdate.Set("name", name)
	if description != "" {
		sqlUpdate = sqlUpdate.Set("description", description)
	}
	sqlUpdate = sqlUpdate.Suffix("RETURNING id, project_id, name, description, priority, removed, created_at")

	sqlUpdateGood, args, err := sqlUpdate.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to build UPDATE query: %w", err)
	}

	var good storage.Good
	err = s.db.QueryRowContext(ctx, sqlUpdateGood, args...).Scan(
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
