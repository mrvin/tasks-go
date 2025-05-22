package config

import (
	"log/slog"
	"os"

	"github.com/mrvin/tasks-go/quotes/internal/httpserver"
	"github.com/mrvin/tasks-go/quotes/internal/logger"
)

type Config struct {
	HTTP   httpserver.Conf
	Logger logger.Conf
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
}
