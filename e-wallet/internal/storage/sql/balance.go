package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

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
