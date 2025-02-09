package sqlstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/mrvin/tasks-go/pinger/internal/storage"
)

func (s *Storage) CreateHost(ctx context.Context, host *storage.Host) error {
	if _, err := s.insertHost.ExecContext(ctx, host.Name, host.IP.String()); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code.Name() == "unique_violation" {
				return storage.ErrHostExists
			}
		}

		return fmt.Errorf("create host: %w", err)
	}

	return nil
}
