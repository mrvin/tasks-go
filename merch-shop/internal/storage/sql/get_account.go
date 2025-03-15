package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
)

func (s *Storage) GetAccount(ctx context.Context, userName string) (*storage.Account, error) {
	var account storage.Account

	if err := s.getAccount.QueryRowContext(ctx, userName).Scan(&account.Name, &account.HashPassword, &account.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", storage.ErrAccountNotFound, userName)
		}
		return nil, fmt.Errorf("can't scan user with name: %s: %w", userName, err)
	}

	return &account, nil
}
