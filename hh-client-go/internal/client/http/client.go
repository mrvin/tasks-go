package httpclient

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

const requestTimeout = 10

//nolint:tagliatelle
type ConfAPIhh struct {
	ClientID          string `yaml:"client_id"`
	ClientSecret      string `yaml:"client_secret"`
	AuthorizationCode string `yaml:"authorization_code"`
}

type AppInfo struct {
	Name    string
	Version string
	Email   string
}

type Client struct {
	http.Client
	hhUserAgent   string
	userAuth      *auth
	mutexUserAuth sync.RWMutex
}

func New(ctx context.Context, confHH *ConfAPIhh, appInfo *AppInfo) (*Client, error) {
	var client Client

	client.hhUserAgent = appInfo.Name + "/" + appInfo.Version + " (" + appInfo.Email + ")"

	userAuth, err := client.getToken(ctx, confHH.ClientID, confHH.ClientSecret, confHH.AuthorizationCode)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	client.mutexUserAuth.Lock()
	client.userAuth = userAuth
	client.mutexUserAuth.Unlock()

	client.mutexUserAuth.RLock()
	fmt.Println(client.userAuth)
	time.AfterFunc(client.userAuth.expiresIn, client.refreshToken)
	slog.Info("Refresh token will start", slog.String("duration", client.userAuth.expiresIn.String()))
	client.mutexUserAuth.RUnlock()

	return &client, nil
}
