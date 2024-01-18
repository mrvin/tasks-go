package sqlstorage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Storage) Create(ctx context.Context, balance float64) (uuid.UUID, error) {
	var id uuid.UUID
	if err := s.insertWallet.QueryRowContext(ctx, balance).Scan(&id); err != nil {
		return uuid.Nil, fmt.Errorf("create: %w", err)
	}

	return id, nil
}
