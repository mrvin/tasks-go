package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
)

func (s *Storage) GetBalance(ctx context.Context, userName string) (uint64, error) {
	var balance uint64

	if err := s.getBalance.QueryRowContext(ctx, userName).Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%w: %s", storage.ErrAccountNotFound, userName)
		}
		return 0, fmt.Errorf("can't scan user with name: %s: %w", userName, err)
	}

	return balance, nil
}
