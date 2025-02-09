package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/mrvin/tasks-go/pinger/internal/logger"
	"github.com/mrvin/tasks-go/pinger/internal/pinger/app"
	"github.com/mrvin/tasks-go/pinger/internal/pinger/httpclient"
)

type Config struct {
	App    app.Conf
	HTTP   httpclient.Conf
	Logger logger.Conf
}

const defaultPingPeriod = 30 * time.Second
const defaultPingTimeout = 3 * time.Second
const defaultPingCountPackets = 1

// LoadFromEnv will load configuration solely from the environment.
func (c *Config) LoadFromEnv() {
	if pingPeriodStr := os.Getenv("PING_PERIOD"); pingPeriodStr != "" {
		pingPeriod, err := time.ParseDuration(pingPeriodStr)
		if err != nil {
			slog.Warn("Invalid ping period " + pingPeriodStr)
		}
		c.App.PingPeriod = pingPeriod
	} else {
		c.App.PingPeriod = defaultPingPeriod
		slog.Warn("Ping period set default " + time.Duration(defaultPingPeriod).String())
	}

	if pingTimeoutStr := os.Getenv("PING_PERIOD"); pingTimeoutStr != "" {
		pingTimeout, err := time.ParseDuration(pingTimeoutStr)
		if err != nil {
			slog.Warn("Invalid ping timeout " + pingTimeoutStr)
		}
		c.App.PingTimeout = pingTimeout
	} else {
		c.App.PingTimeout = defaultPingTimeout
		slog.Warn("Ping timeout set default " + time.Duration(defaultPingTimeout).String())
	}

	if pingCountPacketsStr := os.Getenv("PING_COUNT_PACKETS"); pingCountPacketsStr != "" {
		pingCountPackets, err := strconv.Atoi(pingCountPacketsStr)
		if err != nil {
			slog.Warn("Invalid ping count packets " + pingCountPacketsStr)
		}
		c.App.PingCountPackets = pingCountPackets
	} else {
		c.App.PingCountPackets = defaultPingCountPackets
		slog.Warn("Ping count packets set default " + fmt.Sprintf("%d", defaultPingCountPackets))
	}

	if addr := os.Getenv("BACKEND_HTTP_ADDR"); addr != "" {
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
