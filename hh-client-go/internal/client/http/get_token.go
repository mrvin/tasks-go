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

type respGetToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type auth struct {
	accessToken  string
	refreshToken string
	expiresAt    time.Time
}

func (c *Client) getToken(ctx context.Context, clientID, clientSecret, code string) (*auth, error) {
	const requestURL = "https://api.hh.ru/token"

	ctx, cancel := context.WithTimeout(ctx, requestTimeout*time.Second)
	defer cancel()

	data := bytes.NewBufferString("grant_type=authorization_code&client_id=" + clientID + "&client_secret=" + clientSecret + "&code=" + code)
	// Create a new request using http
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("HH-User-Agent", c.hhUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	slog.Info("Get token: ", slog.String("url", requestURL))

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("http response status code: %d", resp.StatusCode)
		if 400 <= resp.StatusCode && resp.StatusCode < 500 {
			if resp.StatusCode == 400 {
				dataError := struct {
					Summary          string `json:"summary"`
					Error            string `json:"error"`
					ErrorDescription string `json:"error_description"`
				}{}
				if err := json.NewDecoder(resp.Body).Decode(&dataError); err != nil {
					slog.Error("Unmarshal body response: " + err.Error())
				}
				err = fmt.Errorf("%w; summary: %s; error: %s; error_description: %s",
					err,
					dataError.Summary,
					dataError.Error,
					dataError.ErrorDescription,
				)
			}
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
		expiresAt:    time.Now().Add(time.Second * time.Duration(respData.ExpiresIn)),
	}, nil
}
