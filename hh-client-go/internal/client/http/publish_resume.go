package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

//nolint:tagliatelle
type RequestError struct {
	StatusCode  int
	RequestID   string `json:"request_id"`
	Description string `json:"description"`
	Errors      []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	OauthError string `json:"oauth_error"`
}

func (e *RequestError) Error() string {
	errStr := fmt.Sprintf("http response status code: %d;", e.StatusCode)
	if e.RequestID != "" {
		errStr += fmt.Sprintf("request_id: %q;", e.RequestID)
	}
	if e.Description != "" {
		errStr += fmt.Sprintf("request_id: %q;", e.Description)
	}
	for _, err := range e.Errors {
		if err.Type != "" || err.Value != "" {
			errStr += fmt.Sprintf("Type: %q, Value: %q;", err.Type, err.Value)
		}
	}
	if e.OauthError != "" {
		errStr += fmt.Sprintf("oauth_error: %q;", e.OauthError)
	}

	return errStr
}

func requestBuildPublishResume(resumeID string) string {
	requestURL := &url.URL{
		Scheme: "https",
		Host:   "api.hh.ru",
		Path:   "resumes",
	}

	requestURL = requestURL.JoinPath(resumeID, "publish")

	return requestURL.String()
}

func (c *Client) PublishResume(ctx context.Context, resumeID string) error {
	requestURL := requestBuildPublishResume(resumeID)

	ctx, cancel := context.WithTimeout(ctx, requestTimeout*time.Second)
	defer cancel()

	// Create a new request using http
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	c.mutexUserAuth.RLock()
	bearer := "Bearer " + c.userAuth.accessToken
	c.mutexUserAuth.RUnlock()
	req.Header.Add("Authorization", bearer)
	req.Header.Add("HH-User-Agent", c.hhUserAgent) //nolint:canonicalheader

	slog.Info("Publish resume", slog.String("url", requestURL))

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var errReq RequestError
		errReq.StatusCode = resp.StatusCode
		if resp.StatusCode == http.StatusBadRequest ||
			resp.StatusCode == http.StatusForbidden ||
			resp.StatusCode == http.StatusNotFound ||
			resp.StatusCode == http.StatusTooManyRequests {
			if err := json.NewDecoder(resp.Body).Decode(&errReq); err != nil {
				slog.Error("Unmarshal body response: " + err.Error())
			}
		}
		return &errReq
	}

	return nil
}
