package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrNoWalletID = errors.New("no wallet with id")

func (s *Storage) Balance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	var balance float64
	if err := s.getBalance.QueryRowContext(ctx, walletID).Scan(&balance); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("balance: %w %v", ErrNoWalletID, walletID)
		}
		return 0, fmt.Errorf("balance: %w", err)
	}

	return balance, nil
}
