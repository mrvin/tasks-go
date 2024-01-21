package sqlstorage

import (
	"context"
	"fmt"
)

func (s *Storage) Delete(ctx context.Context, id int64) error {
	if _, err := s.deletePerson.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("delete person: %w", err)
	}
	return nil
}
