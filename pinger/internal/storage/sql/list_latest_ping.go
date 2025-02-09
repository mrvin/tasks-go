package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"

	"github.com/mrvin/tasks-go/pinger/internal/storage"
)

func (s *Storage) ListLatestPing(ctx context.Context) ([]storage.Ping, error) {
	slPing := make([]storage.Ping, 0)
	rows, err := s.selectLatestPing.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return slPing, nil
		}
		return nil, fmt.Errorf("can't get ping: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ping storage.Ping
		var ipStr string
		err = rows.Scan(&ipStr, &ping.Time, &ping.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		ping.IP = net.ParseIP(ipStr)
		slPing = append(slPing, ping)
	}
	if err := rows.Err(); err != nil {
		return slPing, fmt.Errorf("rows error: %w", err)
	}

	return slPing, nil
}
