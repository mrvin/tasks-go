package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/notes/internal/storage"
)

func (s *Storage) CreateNote(ctx context.Context, userName string, note *storage.Note) (int64, error) {
	if err := s.insertNote.QueryRowContext(ctx,
		note.Title,
		note.Description,
		userName,
	).Scan(&note.ID); err != nil {
		return 0, fmt.Errorf("create note: %w", err)
	}

	return note.ID, nil
}

func (s *Storage) ListNotes(ctx context.Context, userName string) ([]storage.Note, error) {
	notes := make([]storage.Note, 0)

	rows, err := s.listNotes.QueryContext(ctx, userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notes, nil
		}
		return nil, fmt.Errorf("can't get list notes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var note storage.Note
		err = rows.Scan(
			&note.ID,
			&note.Title,
			&note.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return notes, fmt.Errorf("rows error: %w", err)
	}

	return notes, nil
}
