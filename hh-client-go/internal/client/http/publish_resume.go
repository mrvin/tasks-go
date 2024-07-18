package httpclient

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

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
		return err
	}

	bearer := "Bearer " + c.userAuth.accessToken
	req.Header.Add("Authorization", bearer)
	req.Header.Add("HH-User-Agent", c.hhUserAgent)

	slog.Info("Publish resume: ", slog.String("url", requestURL))

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		err := fmt.Errorf("http response status code: %d", resp.StatusCode)
		if 400 <= resp.StatusCode && resp.StatusCode < 500 {
			// Detailed error output.
		}

		return err
	}

	return nil
}
