package sqlstorage

import (
	"context"
	"fmt"

	"github.com/mrvin/tasks-go/pinger/internal/storage"
)

func (s *Storage) CreatePing(ctx context.Context, ping *storage.Ping) error {
	if _, err := s.insertPing.ExecContext(ctx, ping.IP.String(), ping.Time, ping.CreatedAt); err != nil {
		return fmt.Errorf("create ping: %w", err)
	}

	return nil
}
