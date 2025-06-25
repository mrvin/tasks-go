package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/mrvin/tasks-go/goods/internal/storage"
)

const (
	maxOpenConns    = 10
	maxIdleConns    = 5
	connMaxLifetime = 1 // Hour
)

type Conf struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Storage struct {
	conn driver.Conn
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", conf.Host, conf.Port)},
		Auth: clickhouse.Auth{
			Database: conf.Name,
			Username: conf.User,
			Password: conf.Password,
		},
		Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &Storage{conn}, nil
}

func (s *Storage) BatchInsert(ctx context.Context, batchEvents []storage.Event) error {
	batch, err := s.conn.PrepareBatch(ctx, "INSERT INTO goods_events")
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}
	defer batch.Close()

	for _, event := range batchEvents {
		if err := batch.Append(
			event.ID,
			event.ProjectID,
			event.Name,
			event.Description,
			event.Priority,
			event.Removed,
			event.Time,
		); err != nil {
			return fmt.Errorf("append batch: %w", err)
		}
	}

	return batch.Send() //nolint:wrapcheck
}

func (s *Storage) Close() error {
	return s.conn.Close() //nolint:wrapcheck
}
