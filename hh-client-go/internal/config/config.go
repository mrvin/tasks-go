package config

import (
	"log/slog"
	"os"

	clienthttp "github.com/mrvin/tasks-go/hh-client-go/internal/client/http"
	"github.com/mrvin/tasks-go/hh-client-go/internal/logger"
)

//nolint:tagliatelle
type Config struct {
	AuthHH clienthttp.ConfAPIhh `yaml:"auth_hh"`
	Logger logger.Conf          `yaml:"authorization_code"`
}

// LoadFromEnv will load configuration solely from the environment.
func (c *Config) LoadFromEnv() {
	if clientID := os.Getenv("CLIENT_ID"); clientID != "" {
		c.AuthHH.ClientID = clientID
	} else {
		slog.Warn("Empty client id")
	}
	if clientSecret := os.Getenv("CLIENT_SECRET"); clientSecret != "" {
		c.AuthHH.ClientSecret = clientSecret
	} else {
		slog.Warn("Empty client secret")
	}
	if authorizationCode := os.Getenv("AUTHORIZATION_CODE"); authorizationCode != "" {
		c.AuthHH.AuthorizationCode = authorizationCode
	} else {
		slog.Warn("Empty authorization code")
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
