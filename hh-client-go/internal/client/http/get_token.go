package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

//nolint:tagliatelle
type respGetToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

//nolint:tagliatelle
type errorGetToken struct {
	Summary          string `json:"summary"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type auth struct {
	accessToken  string
	refreshToken string
	expiresIn    time.Duration
}

func (c *Client) getToken(ctx context.Context, clientID, clientSecret, code string) (*auth, error) {
	const requestURL = "https://api.hh.ru/token"

	ctx, cancel := context.WithTimeout(ctx, requestTimeout*time.Second)
	defer cancel()

	var data *bytes.Buffer
	if clientID == "" || clientSecret == "" || code == "" {
		c.mutexUserAuth.RLock()
		data = bytes.NewBufferString("grant_type=refresh_token&refresh_token=" + c.userAuth.refreshToken)
		c.mutexUserAuth.RUnlock()
	} else {
		data = bytes.NewBufferString("grant_type=authorization_code&client_id=" + clientID + "&client_secret=" + clientSecret + "&code=" + code)
	}

	// Create a new request using http
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, data)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Add("HH-User-Agent", c.hhUserAgent) //nolint:canonicalheader
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	slog.Info("Get token", slog.String("url", requestURL))

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("http response status code: %d", resp.StatusCode)
		if resp.StatusCode == http.StatusBadRequest {
			var errorData errorGetToken
			if err := json.NewDecoder(resp.Body).Decode(&errorData); err != nil {
				slog.Error("Unmarshal body response: " + err.Error())
			}
			err = fmt.Errorf("%w; summary: %s; error: %s; error_description: %s",
				err,
				errorData.Summary,
				errorData.Error,
				errorData.ErrorDescription,
			)
		}
		return nil, err
	}

	var respData respGetToken
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("unmarshal body response: %w", err)
	}

	return &auth{
		accessToken:  respData.AccessToken,
		refreshToken: respData.RefreshToken,
		expiresIn:    time.Second * time.Duration(respData.ExpiresIn),
	}, nil
}
