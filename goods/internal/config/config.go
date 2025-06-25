package config

import (
	"log/slog"
	"os"

	"github.com/mrvin/tasks-go/goods/internal/httpserver"
	"github.com/mrvin/tasks-go/goods/internal/logger"
	natsmq "github.com/mrvin/tasks-go/goods/internal/queue/nats"
	chstorage "github.com/mrvin/tasks-go/goods/internal/storage/clickhouse"
	sqlstorage "github.com/mrvin/tasks-go/goods/internal/storage/sql"
)

type Config struct {
	HTTP       httpserver.Conf
	Logger     logger.Conf
	Postgres   sqlstorage.Conf
	ClickHouse chstorage.Conf
	MQ         natsmq.Conf
}

// LoadFromEnv will load configuration solely from the environment.
func (c *Config) LoadFromEnv() {
	if addr := os.Getenv("API_HTTP_ADDR"); addr != "" {
		c.HTTP.Addr = addr
	} else {
		slog.Warn("Empty http addr")
	}

	if logFilePath := os.Getenv("LOGGER_FILEPATH"); logFilePath != "" {
		c.Logger.FilePath = logFilePath
	} else {
		slog.Warn("Empty log file path")
	}
	if logLevel := os.Getenv("LOGGER_LEVEL"); logLevel != "" {
		c.Logger.Level = logLevel
	} else {
		slog.Warn("Empty log level")
	}

	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		c.Postgres.Host = host
	} else {
		slog.Warn("Empty postgres host")
	}
	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		c.Postgres.Port = port
	} else {
		slog.Warn("Empty postgres port")
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		c.Postgres.User = user
	} else {
		slog.Warn("Empty postgres user")
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		c.Postgres.Password = password
	} else {
		slog.Warn("Empty postgres password")
	}
	if name := os.Getenv("POSTGRES_DB"); name != "" {
		c.Postgres.Name = name
	} else {
		slog.Warn("Empty postgres db name")
	}

	if host := os.Getenv("CLICKHOUSE_HOST"); host != "" {
		c.ClickHouse.Host = host
	} else {
		slog.Warn("Empty clickhouse host")
	}
	if port := os.Getenv("CLICKHOUSE_PORT"); port != "" {
		c.ClickHouse.Port = port
	} else {
		slog.Warn("Empty clickhouse port")
	}
	if user := os.Getenv("CLICKHOUSE_USER"); user != "" {
		c.ClickHouse.User = user
	} else {
		slog.Warn("Empty clickhouse user")
	}
	if password := os.Getenv("CLICKHOUSE_PASSWORD"); password != "" {
		c.ClickHouse.Password = password
	} else {
		slog.Warn("Empty clickhouse password")
	}
	if name := os.Getenv("CLICKHOUSE_DB"); name != "" {
		c.ClickHouse.Name = name
	} else {
		slog.Warn("Empty clickhouse db name")
	}

	if host := os.Getenv("NATS_HOST"); host != "" {
		c.MQ.Host = host
	} else {
		slog.Warn("Empty nats host")
	}
	if port := os.Getenv("NATS_PORT"); port != "" {
		c.MQ.Port = port
	} else {
		slog.Warn("Empty nats port")
	}
	if subject := os.Getenv("NATS_SUBJECT"); subject != "" {
		c.MQ.Subject = subject
	} else {
		slog.Warn("Empty nats subject")
	}
}
