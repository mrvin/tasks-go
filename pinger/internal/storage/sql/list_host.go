package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"

	"github.com/mrvin/tasks-go/pinger/internal/storage"
)

func (s *Storage) ListHost(ctx context.Context) ([]storage.Host, error) {
	slHost := make([]storage.Host, 0)
	rows, err := s.selectIP.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return slHost, nil
		}
		return nil, fmt.Errorf("can't get ip: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ipStr string
		var host storage.Host
		err = rows.Scan(&host.Name, &ipStr)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		host.IP = net.ParseIP(ipStr)
		slHost = append(slHost, host)
	}
	if err := rows.Err(); err != nil {
		return slHost, fmt.Errorf("rows error: %w", err)
	}

	return slHost, nil
}
