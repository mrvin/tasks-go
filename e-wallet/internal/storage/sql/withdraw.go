package sqlstorage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Storage) Withdraw(ctx context.Context, walletID uuid.UUID, amount float64) error {
	if _, err := s.withdraw.ExecContext(ctx, walletID, amount); err != nil {
		return fmt.Errorf("withdraw: %w", err)
	}
	return nil
}
