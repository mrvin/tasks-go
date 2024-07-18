package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type resumeID struct {
	ID string `json:"id"`
}

type respListResumeID struct {
	Found int        `json:"found"`
	Pages int        `json:"pages"`
	Items []resumeID `json:"items"`
}

func (c *Client) ListResumeID(ctx context.Context) ([]string, error) {
	const requestURL = "https://api.hh.ru/resumes/mine"

	ctx, cancel := context.WithTimeout(ctx, requestTimeout*time.Second)
	defer cancel()

	// Create a new request using http
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	bearer := "Bearer " + c.userAuth.accessToken
	req.Header.Add("Authorization", bearer)
	req.Header.Add("HH-User-Agent", c.hhUserAgent)

	slog.Info("Get list resume id: ", slog.String("url", requestURL))

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("http response status code: %d", resp.StatusCode)
		if 400 <= resp.StatusCode && resp.StatusCode < 500 {
			// Detailed error output.
		}
		return nil, err
	}

	var respData respListResumeID
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("unmarshal body response: %w", err)
	}

	if respData.Pages > 1 {
		slog.Warn("too many resumes")
	}

	listResumeID := make([]string, 0, respData.Found)

	for _, resume := range respData.Items {
		listResumeID = append(listResumeID, resume.ID)
	}

	return listResumeID, nil
}
