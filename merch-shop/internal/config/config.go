package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/mrvin/tasks-go/merch-shop/internal/app"
	"github.com/mrvin/tasks-go/merch-shop/internal/auth"
	"github.com/mrvin/tasks-go/merch-shop/internal/httpserver"
	"github.com/mrvin/tasks-go/merch-shop/internal/logger"
	sqlstorage "github.com/mrvin/tasks-go/merch-shop/internal/storage/sql"
)

type Config struct {
	App    app.Conf
	Auth   auth.Conf
	HTTP   httpserver.Conf
	DB     sqlstorage.Conf
	Logger logger.Conf
}

const (
	defaultStartingBalance     = 1000
	defaultTokenValidityPeriod = 24 * time.Hour
)

// LoadFromEnv will load configuration solely from the environment.
func (c *Config) LoadFromEnv() {
	if startingBalanceStr := os.Getenv("STARTING_BALANCE"); startingBalanceStr != "" {
		startingBalance, err := strconv.ParseUint(startingBalanceStr, 10, 64)
		if err != nil {
			slog.Warn("Invalid starting balance " + startingBalanceStr)
		}
		c.App.StartingBalance = startingBalance
	} else {
		c.App.StartingBalance = defaultStartingBalance
		slog.Warn("Starting balance set default " + fmt.Sprintf("%d", defaultStartingBalance))
	}
	if secretKey := os.Getenv("SECRET_KEY"); secretKey != "" {
		c.Auth.SecretKey = secretKey
	} else {
		slog.Warn("Empty secret key")
	}
	if tokenValidityPeriodStr := os.Getenv("TOKEN_VALIDITY_PERIOD"); tokenValidityPeriodStr != "" {
		tokenValidityPeriod, err := time.ParseDuration(tokenValidityPeriodStr)
		if err != nil {
			slog.Warn("Invalid token validity period " + tokenValidityPeriodStr)
		}
		c.Auth.TokenValidityPeriod = tokenValidityPeriod
	} else {
		c.Auth.TokenValidityPeriod = defaultTokenValidityPeriod
		slog.Warn("Token validity period set default " + time.Duration(defaultTokenValidityPeriod).String())
	}

	if addr := os.Getenv("API_HTTP_ADDR"); addr != "" {
		c.HTTP.Addr = addr
	} else {
		slog.Warn("Empty http addr")
	}

	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		c.DB.Host = host
	} else {
		slog.Warn("Empty postgres host")
	}
	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		c.DB.Port = port
	} else {
		slog.Warn("Empty postgres port")
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		c.DB.User = user
	} else {
		slog.Warn("Empty postgres user")
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		c.DB.Password = password
	} else {
		slog.Warn("Empty postgres password")
	}
	if name := os.Getenv("POSTGRES_DB"); name != "" {
		c.DB.Name = name
	} else {
		slog.Warn("Empty postgres db name")
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
