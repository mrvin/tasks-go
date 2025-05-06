package sqlstorage

import (
	"context"
	"fmt"

	"github.com/mrvin/tasks-go/persons/internal/storage"
)

func (s *Storage) Delete(ctx context.Context, id int64) error {
	res, err := s.deletePerson.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("delete person: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete url: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w %v", storage.ErrNoPersonID, id)
	}
	return nil
}
