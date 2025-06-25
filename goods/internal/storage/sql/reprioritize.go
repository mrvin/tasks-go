package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/goods/internal/storage"
)

func (s *Storage) Reprioritize(ctx context.Context, id, projectID, newPriority int64) (*storage.Good, []storage.Priority, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var good storage.Good
	sqlSelectGood := `
		SELECT id, project_id, name, description, priority, removed, created_at
		FROM goods
		WHERE id = $1 AND project_id = $2
	`
	err = tx.QueryRowContext(ctx, sqlSelectGood, id, projectID).Scan(
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
			return nil, nil, fmt.Errorf("%w %v", storage.ErrNoGoodID, id)
		}
		return nil, nil, fmt.Errorf("failed get good: %w", err)
	}

	if good.Priority == newPriority {
		return &good, []storage.Priority{}, nil
	}

	var updateQuery string

	if newPriority < good.Priority {
		updateQuery = `
            UPDATE goods 
            SET priority = priority + 1 
            WHERE priority >= $1 AND priority < $2
            RETURNING id, priority
        `
	} else {
		updateQuery = `
            UPDATE goods 
            SET priority = priority - 1 
            WHERE priority > $2 AND priority <= $1
            RETURNING id, priority
        `
	}

	rows, err := tx.QueryContext(ctx, updateQuery, newPriority, good.Priority)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update priorities: %w", err)
	}
	defer rows.Close()

	var priorities []storage.Priority
	for rows.Next() {
		var p storage.Priority
		if err := rows.Scan(&p.ID, &p.Priority); err != nil {
			return nil, nil, fmt.Errorf("failed to scan priority: %w", err)
		}
		priorities = append(priorities, p)
	}
	sqlUpdatePriority := `
		UPDATE goods
		SET priority = $1 
		WHERE id = $2 AND project_id = $3
	`
	_, err = tx.ExecContext(ctx, sqlUpdatePriority, newPriority, id, projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update current good: %w", err)
	}

	priorities = append(priorities, storage.Priority{
		ID:       id,
		Priority: newPriority,
	})

	if err := tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return &good, priorities, nil
}
