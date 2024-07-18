package httpclient

import (
	"context"
	"crypto/tls"
	"fmt"

	"net/http"
)

const requestTimeout = 10

type ConfAPIhh struct {
	ClientID          string `yaml:"client_id"`
	ClientSecret      string `yaml:"client_secret"`
	AuthorizationCode string `yaml:"authorization_code"`
}

type Conf struct {
	ClientCrt string `yaml:"cert_file"`
	ClientKey string `yaml:"key_file"`
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

func New(ctx context.Context, conf *Conf, confHH *ConfAPIhh, appInfo *AppInfo) (*Client, error) {
	var client Client
	cert, err := tls.LoadX509KeyPair(conf.ClientCrt, conf.ClientKey)
	if err != nil {
		return nil, err
	}
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConf,
	}
	client.Client = http.Client{
		Transport: transport,
	}
	client.hhUserAgent = appInfo.Name + "/" + appInfo.Version + " (" + appInfo.Email + ")"

	client.userAuth, err = client.getToken(ctx, confHH.ClientID, confHH.ClientSecret, confHH.AuthorizationCode)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	return &client, nil
}
