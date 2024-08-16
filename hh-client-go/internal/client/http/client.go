package httpclient

import (
	"context"
	"fmt"
	"net/http"
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
	userAuth    *auth
	hhUserAgent string
}

func New(ctx context.Context, confHH *ConfAPIhh, appInfo *AppInfo) (*Client, error) {
	var err error
	var client Client

	client.hhUserAgent = appInfo.Name + "/" + appInfo.Version + " (" + appInfo.Email + ")"

	client.userAuth, err = client.getToken(ctx, confHH.ClientID, confHH.ClientSecret, confHH.AuthorizationCode)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	fmt.Println(client.userAuth)

	return &client, nil
}
