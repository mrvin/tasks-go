package sqlstorage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Storage) Deposit(ctx context.Context, walletID uuid.UUID, amount float64) error {
	if _, err := s.deposit.ExecContext(ctx, walletID, amount); err != nil {
		return fmt.Errorf("deposit: %w", err)
	}
	return nil
}
