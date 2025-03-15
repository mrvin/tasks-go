package sqlstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
)

func (s *Storage) CreateAccount(ctx context.Context, userName, hashPassword string, startingBalance uint64) error {
	if _, err := s.insertAccount.ExecContext(ctx, userName, hashPassword, startingBalance); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code.Name() == "unique_violation" {
				return storage.ErrAccountExists
			}
		}

		return fmt.Errorf("create account: %w", err)
	}

	return nil
}
